package information

import (
	"fmt"
	"strings"
	"time"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

const DefaultCrewTag = "ROS"

func Information(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordID := m.Author.ID

	var (
		gamertag    string
		linkdiscord int
		crewTag     string
		lastOnline  string
	)

	err := utils.DB().QueryRow(`SELECT gamertag, linkdiscord, crew_tag, last_online FROM members WHERE discordid = ?`, discordID).Scan(&gamertag, &linkdiscord, &crewTag, &lastOnline)
	if err != nil {
		utils.Send(s, m.ChannelID, "No Information Found", "No information found for your Discord ID. Please link your gamertag using `?link`.", 0xff0000,)
		return
	}

	if gamertag == "" {
		gamertag = "Not Set"
	}

	if crewTag == "" {
		crewTag = DefaultCrewTag
	}

	formattedLastOnline := "Not Available"
	if lastOnline != "" {
		formattedLastOnline = utils.FormatExactTime(lastOnline)
	}

	var (
		crewName  string
		crewColor string
		crewMotto string
		crewOwner string
	)

	if crewTag != DefaultCrewTag {
		_ = utils.DB().QueryRow(`SELECT crew_name, crew_color, crew_motto, crew_owner FROM crews WHERE crew_tag = ?`, crewTag).Scan(&crewName, &crewColor, &crewMotto, &crewOwner)
	}

	if crewColor == "" {
		crewColor = "Red"
	}

	isOwner := crewOwner == discordID
	currentCrew := []string{}

	if crewTag != DefaultCrewTag && !isOwner {

		ownerName := "Unknown"
		if crewOwner != "" {
			if user, err := s.User(crewOwner); err == nil {
				ownerName = user.Username
			}
		}

		currentCrew = []string{
			"",
			"**Current Crew**",
			fmt.Sprintf("**Name:** %s", utils.Fallback(crewName, "Not Set")),
			fmt.Sprintf("**Tag:** %s", crewTag),
			fmt.Sprintf("**Color:** %s", crewColor),
			fmt.Sprintf("**Motto:** %s", utils.Fallback(crewMotto, "Not Set")),
			fmt.Sprintf("**Owner:** %s", ownerName),
		}
	}

	ownedCrew := []string{}

	if isOwner {

		var (
			oName  string
			oTag   string
			oColor string
			oMotto string
		)

		err := utils.DB().QueryRow(`SELECT crew_name, crew_tag, crew_color, crew_motto FROM crews WHERE crew_owner = ?`, discordID).Scan(&oName, &oTag, &oColor, &oMotto)
		if err == nil {
			ownedCrew = []string{
				"",
				"**Owned Crew**",
				fmt.Sprintf("**Name:** %s", oName),
				fmt.Sprintf("**Tag:** %s", oTag),
				fmt.Sprintf("**Color:** %s", oColor),
				fmt.Sprintf("**Motto:** %s", oMotto),
			}
		}
	}

	description := strings.Join(
		append([]string{
			fmt.Sprintf("**Discord ID:** ||%s||", discordID),
			fmt.Sprintf("**Gamertag:** ||%s||", gamertag),
			fmt.Sprintf("**Crew Tag:** %s", crewTag),
			fmt.Sprintf("**Link Status:** %s", map[bool]string{
				true:  "Linked",
				false: "Not Linked",
			}[linkdiscord == 1]),
			fmt.Sprintf("**Last Online:** %s", formattedLastOnline),
		}, append(currentCrew, ownedCrew...)...), "\n",)

	embed := &discordgo.MessageEmbed{
		Title:       "User Information",
		Description: description,
		Color:       utils.ParseColor(crewColor),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Info Command",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)
}