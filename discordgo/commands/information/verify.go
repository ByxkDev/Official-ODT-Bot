package information

import (
	"fmt"
	"strings"
    "discordgo/utils"
	"github.com/bwmarrin/discordgo"
)


func Verify(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		utils.SendError(s, m, "Missing Code", "Usage: ?verify <code>",)
		return
	}

	code := args[1]
	discordID := m.Author.ID

	var gamertag string

	err := utils.DB().QueryRow(`SELECT gamertag FROM members WHERE discordcode = ?`, code).Scan(&gamertag)
	if err != nil || gamertag == "" {
		utils.SendError(s, m, "Invalid Verification Code", "The code you entered is not valid. Please try again.",)
		return
	}

	_, err = utils.DB().Exec(`UPDATE members SET linkdiscord = 1, discordid = ? WHERE discordcode = ?`, discordID, code)
	if err != nil {
		utils.SendError(s, m, "Error", "Failed to verify account. Please try again later.",)
		return
	}

	_ = utils.SendWebhook(s, fmt.Sprintf("New Verification\nGamertag: %s\nDiscord: <@%s>", gamertag, discordID,))

	embed := &discordgo.MessageEmbed{
		Title:       "Verification Successful",
		Description: fmt.Sprintf("Gamertag **%s** has been successfully linked.", gamertag),
		Color:       0x00ff00,

		Footer: &discordgo.MessageEmbedFooter{
			Text: "Verify Command",
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: embed,
	})
}