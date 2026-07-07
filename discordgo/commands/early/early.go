package early

import (
	"database/sql"
	"fmt"
	"os"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

const MaxClaims = 100

func GetEarlySupporterRoleID() string {
	return os.Getenv("EARLY_SUPPORTERID")
}

func GetEarlySupporterCount() (int, error) {
	db := utils.DB()
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM members WHERE earlysupporter = 1`).Scan(&count)
	return count, err
}

func IsEarlySupporter(discordID string) (bool, error) {
	db := utils.DB()
	var value int
	err := db.QueryRow(`SELECT earlysupporter FROM members WHERE discordid = ? LIMIT 1`, discordID).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return value == 1, err
}

func HasLinkedDiscord(discordID string) bool {
	db := utils.DB()
	var id int
	err := db.QueryRow(`SELECT id FROM members WHERE discordid = ? LIMIT 1`, discordID).Scan(&id)
	return err == nil
}

func SetEarlySupporter(discordID string) error {
	db := utils.DB()
	_, err := db.Exec(`UPDATE members SET earlysupporter = 1 WHERE discordid = ?`, discordID)
	return err
}

func Claim(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" {
		return
	}

	roleID := GetEarlySupporterRoleID()

	if roleID == "" {
		fmt.Println("ERROR: EARLY_SUPPORTERID is not set")
		s.ChannelMessageSend(m.ChannelID, "Bot misconfigured (missing role ID).")
		return
	}

	if !HasLinkedDiscord(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Your Discord account isn't linked to a game account.\nUse `?link` first.",)
		return
	}

	claimed, err := IsEarlySupporter(m.Author.ID)
	if err != nil {
		fmt.Println("DB error (IsEarlySupporter):", err)
		s.ChannelMessageSend(m.ChannelID, "Database error.")
		return
	}

	if claimed {
		s.ChannelMessageSend(m.ChannelID, "You've already claimed Early Supporter.",)
		return
	}

	count, err := GetEarlySupporterCount()
	if err != nil {
		fmt.Println("DB error (GetEarlySupporterCount):", err)
		s.ChannelMessageSend(m.ChannelID, "Database error.")
		return
	}

	if count >= MaxClaims {
		s.ChannelMessageSend(m.ChannelID, "The first 100 Early Supporter spots have already been claimed.",)
		return
	}

	err = SetEarlySupporter(m.Author.ID)
	if err != nil {
		fmt.Println("DB error (SetEarlySupporter):", err)
		s.ChannelMessageSend(m.ChannelID, "Database error.")
		return
	}

	err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, roleID,)

	if err != nil {
		fmt.Println("ROLE ADD FAILED:", err)
		fmt.Println("GuildID:", m.GuildID)
		fmt.Println("UserID:", m.Author.ID)
		fmt.Println("RoleID:", roleID)
		s.ChannelMessageSend(m.ChannelID, "I saved your claim but couldn't give the Discord role.",)
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Congratulations! You claimed the **Early Supporter** role!\n\nRemaining spots: **%d/%d**", MaxClaims-(count+1), MaxClaims,),)
}