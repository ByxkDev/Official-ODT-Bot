package information

import (
	"fmt"
	"strings"
    "discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

func Link(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordID := m.Author.ID
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		utils.SendError(s, m,
			"Missing Name",
			"Usage: ?link <gamertag>",
		)
		return
	}

	gamertag := args[1]
	var existing string

	err := utils.DB().QueryRow(`SELECT discordid FROM members WHERE discordid = ?`, discordID).Scan(&existing)
	if err == nil && existing != "" {
		utils.SendError(s, m, "Account Already Linked", "Your Discord account is already linked to a gamertag.",)
		return
	}

	var found string

	err = utils.DB().QueryRow(`SELECT gamertag FROM members WHERE gamertag = ?`, gamertag).Scan(&found)
	if err != nil || found == "" {
		utils.SendError(s, m, "Gamertag Not Found", "Launch GTA and try again with a valid gamertag.",)
		return
	}

	code := utils.GenerateCode()
	_, err = utils.DB().Exec(`UPDATE members SET discordcode = ?, discordid = ? WHERE gamertag = ?`, code, discordID, gamertag)
	if err != nil {
		utils.SendError(s, m, "Error", "Failed to link account. Please try again later.",)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Gamertag Linked Successfully!",
		Description: fmt.Sprintf("Your verification code: **%s**\nUse it with `?verify` to complete linking.", code,),
		Color: 0x00ff00,

		Footer: &discordgo.MessageEmbedFooter{
			Text: "Link Command",
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: embed,
	})
}