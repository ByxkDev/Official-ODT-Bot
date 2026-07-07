package crews

import (
	"fmt"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

func List(s *discordgo.Session, m *discordgo.MessageCreate) {

	rows, err := utils.DB().Query(`SELECT crew_tag, crew_name, crew_owner, crew_motto, crew_invite FROM crews`)
	if err != nil {
		utils.SendError(s, m, "Error", "Failed to load crews.")
		return
	}

	defer rows.Close()

	var crews []utils.Crew
	for rows.Next() {
		var c utils.Crew

		err := rows.Scan(&c.Tag, &c.Name, &c.Owner, &c.Motto, &c.InviteCode,)

		if err != nil {
			continue
		}

		crews = append(crews, c)
	}

	if len(crews) == 0 {
		utils.SendError(s, m, "No Crews Found", "There are currently no crews available.")
		return
	}

	var full []utils.CrewFull
	for _, c := range crews {

		var count int

		err := utils.DB().QueryRow("SELECT COUNT(*) FROM members WHERE crew_tag = ?", c.Tag,).Scan(&count)
		if err != nil {
			count = 0
		}

		ownerName := "Unknown"
		if user, err := s.User(c.Owner); err == nil {
			ownerName = user.Username
		}

		visibility := "Public"
		if c.InviteCode == "" {
			visibility = "Private"
		}

		full = append(full, utils.CrewFull{
			Crew:        c,
			OwnerName:   ownerName,
			MemberCount: count,
			Visibility:  visibility,
		})
	}

	embed := utils.RenderCrewPage(full, 0, 1)
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: embed,
	})

	if err != nil {
		fmt.Println("[LISTCREWS] send error:", err)
		return
	}

	_ = msg
}