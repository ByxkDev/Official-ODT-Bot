package utils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
    "os"
	"encoding/base64"
	"io"
	"crypto/rand"
	"math/big"
	"net/http"
    "bytes"
	"encoding/json"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

var db *sql.DB
func SetDB(database *sql.DB) {
	db = database
}

func DB() *sql.DB {
	return db
}

const (
	DefaultCrewTag = "ROS"
	DefaultCrewID  = 1337
)

func QueryOne(query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

func Exec(query string, args ...any) error {
	_, err := db.Exec(query, args...)
	return err
}

var AllowedIconTypes = []string{
	".png",
	".jpg",
	".jpeg",
	".webp",
}

func IsAllowedIcon(ext string) bool {
	ext = strings.ToLower(ext)
	for _, allowed := range AllowedIconTypes {
		if ext == allowed {
			return true
		}
	}

	return false
}

func SetRoleIcon(s *discordgo.Session, guildID, roleID, imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"icon": "data:image/png;base64," + base64.StdEncoding.EncodeToString(data),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PATCH", "https://discord.com/api/v10/guilds/"+guildID+"/roles/"+roleID, bytes.NewBuffer(body),)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("discord api error: %s", string(b))
	}

	return nil
}

func IsInCrew(userID string) bool {
	var tag string
	err := db.QueryRow(`SELECT crew_tag FROM members WHERE discordid = ?`, userID).Scan(&tag)
	return err == nil && tag != "" && tag != DefaultCrewTag
}

func IsCrewOwner(discordID string) bool {
	var id string
	err := db.QueryRow(`SELECT crew_id FROM crews WHERE crew_owner = ?`, discordID,).Scan(&id)
	return err == nil
}

func GenerateCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	code := make([]byte, length)

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))),)
		if err != nil {
			return ""
		}
		code[i] = chars[index.Int64()]
	}

	return string(code)
}

func FormatTime(ts string) string {
	if ts == "" {
		return "Not Available"
	}

	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}

	return t.Format("02 Jan 2006 15:04")
}

func NormalizePlatform(p string) string {
	switch strings.ToLower(p) {
	case "xbox360":
		return "X360"
	case "ps3":
		return "PS3"
	default:
		return ""
	}
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func toLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func validateCrewTag(tag string) bool {
	tag = strings.TrimSpace(tag)
	return len(tag) > 0 && len(tag) <= 4
}

func validateCrewName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) > 0 && len(name) <= 20
}

func validateMotto(m string) bool {
	return len(m) <= 30
}

func normalizeColor(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "#") && len(input) == 7 {
		return input
	}

	return "#000001"
}

func onlineFilter() string {
	return "last_online > NOW() - INTERVAL 2 MINUTE"
}

func errWrap(prefix string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", prefix, err)
}

func SendError(s *discordgo.Session, i interface{}, title string, description string) {
	content := fmt.Sprintf("**%s**\n%s", title, description)
	switch v := i.(type) {

	case *discordgo.MessageCreate:
		_, _ = s.ChannelMessageSend(v.ChannelID, content)

	case *discordgo.InteractionCreate:
		_ = s.InteractionRespond(v.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   64, 
			},
		})
	}
}

func SendEmbed(s *discordgo.Session, i interface{}, title string, description string) error {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x00ff00,
	}

	switch v := i.(type) {

	case *discordgo.MessageCreate:
		_, err := s.ChannelMessageSendEmbed(v.ChannelID, embed)
		return err

	case *discordgo.InteractionCreate:
		return s.InteractionRespond(v.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  64,
			},
		})
	}

	return nil
}

func SendWebhook(s *discordgo.Session, content string) error {
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return nil
	}

	params := &discordgo.WebhookParams{
		Content: content,
	}

	_, err := s.WebhookExecute(webhookURL, "", false, params)
	return err
}

const ItemsPerPage = 5

func RenderCrewPage(data []CrewFull, page int, total int) *discordgo.MessageEmbed {
	start := page * ItemsPerPage
	end := start + ItemsPerPage

	if end > len(data) {
		end = len(data)
	}

	desc := ""
	for _, c := range data[start:end] {
		desc += fmt.Sprintf("**[%s] %s**\nOwner: %s\nMembers: %d\nVisibility: %s\n\n", c.Crew.Tag, c.Crew.Name, c.OwnerName, c.MemberCount, c.Visibility,)
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Crews (Page %d/%d)", page+1, total),
		Description: desc,
		Color:       0x00ff00,
	}
}

func FormatExactTime(ts string) string {
	if ts == "" {
		return "Not Available"
	}

	t, err := time.Parse(time.RFC3339, ts)

	if err != nil {
		return ts
	}

	return t.Format("02 Jan 2006 15:04")
}

func OwnedCrewText(has bool, name, tag string) string {
	if !has {
		return ""
	}
	return fmt.Sprintf("**Owned Crew:** %s (%s)", name, tag)
}

func MemberCrewText(tag, id string) string {
	if tag == "" {
		return ""
	}
	return fmt.Sprintf("**Member Crew:** %s (%s)", tag, id)
}

func Fallback(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func NullStr(v string, fallback ...string) string {
	if v == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return "N/A"
	}
	return v
}

func Truncate(s string) string {
	if len(s) > 25 {
		return s[:25] + "..."
	}
	return s
}

func FallbackPlatform(p string) string {
	switch strings.ToLower(p) {
	case "xbox360":
		return "Xbox 360"
	case "ps3":
		return "PlayStation 3"
	default:
		return strings.ToUpper(p)
	}
}

func SafeCrew(tag string) string {
	if tag == "" {
		return "None"
	}
	return tag
}

func ParseColor(input string) int {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "red":
		return 0xff0000
	case "green":
		return 0x00ff00
	case "blue":
		return 0x0099ff
	default:
		return 0x5865f2
	}
}

func Send(s *discordgo.Session, channelID, title, description string, color int) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}

	_, _ = s.ChannelMessageSendEmbed(channelID, embed)
}

func GenerateCrewID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
