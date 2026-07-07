package information

import (
	"fmt"
	"strings"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

func PlayerCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	platform := "all"
	if len(args) > 1 {
		platform = strings.ToLower(args[1])
	}

	var filter string

	switch platform {
	case "all":
		filter = "last_online > NOW() - INTERVAL 2 MINUTE"

	case "ps3":
		filter = "last_online > NOW() - INTERVAL 2 MINUTE AND platform_name='ps3'"

	case "xbox360", "xbox":
		filter = "last_online > NOW() - INTERVAL 2 MINUTE AND platform_name='xbox360'"

	default:
		platform = "rpcs3"
		filter = "last_online > NOW() - INTERVAL 2 MINUTE AND platform_name NOT IN ('ps3','xbox360')"
	}

	type Player struct {
		Gamertag string
		Platform string
	}

	var players []Player
	var totalPlayers int

	rows, err := utils.DB().Query(fmt.Sprintf(`SELECT gamertag, platform_name FROM members WHERE %s ORDER BY gamertag ASC`, filter))
	if err != nil {
		utils.SendError(s, m, "Error", "Failed to fetch active players.")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Gamertag, &p.Platform); err != nil {
			continue
		}
		players = append(players, p)
	}

	if err := utils.DB().QueryRow(`SELECT COUNT(*) FROM members`).Scan(&totalPlayers); 
	err != nil {
		totalPlayers = 0
	}

	var list strings.Builder

	if len(players) == 0 {
		list.WriteString("No active players found.")
	} else {
		for _, p := range players {
			displayPlatform := utils.FallbackPlatform(p.Platform)

			if platform == "rpcs3" {
				displayPlatform = "RPCS3"
			}

			list.WriteString(fmt.Sprintf("> %s (%s)\n", p.Gamertag, displayPlatform,))
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Current Players Online - %s", strings.ToUpper(platform)),
		Description: fmt.Sprintf("**%d player(s) online**\n\n%s\n\n**Total Registered Players:** %d", len(players), list.String(), totalPlayers,),
		Color: 0xff0000,
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)
}