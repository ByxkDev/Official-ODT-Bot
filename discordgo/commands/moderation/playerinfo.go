package moderation

import (
	"fmt"
	"strings"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

func PlayerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) < 2 {
		utils.SendError(s, m, "Missing Name", "Usage: ?playerinfo <gamertag>",)
		return
	}

	gamertag := args[1]

	var (
		xuid           string
		playerGamertag string
		crewID         string
		crewTag        string
		sessionTicket  string
		sessionKey     string
		linkdiscord    int
		discordID      string
		lastOnline     string
		banned         int
		discordCode    string
		platformName   string
	)

	err := utils.DB().QueryRow(`SELECT xuid, gamertag, crew_id, crew_tag, session_ticket, session_key, linkdiscord, discordid, last_online, banned, discordcode, platform_name FROM members WHERE gamertag = ?`, gamertag).Scan(&xuid, &playerGamertag, &crewID, &crewTag, &sessionTicket, &sessionKey, &linkdiscord, &discordID, &lastOnline, &banned, &discordCode, &platformName,)
	if err != nil {
		utils.SendError(s, m, "No Player Found", fmt.Sprintf("No player found with gamertag: %s", gamertag,),)
		return
	}

	isBanned := banned == 1

	var (
		crewName  string
		crewColor string
		crewOwner string
	)

	err = utils.DB().QueryRow(`SELECT crew_name, crew_color, crew_owner FROM crews WHERE crew_owner = ?`, discordID).Scan(&crewName, &crewColor, &crewOwner,)
	hasOwnedCrew := err == nil

	desc := fmt.Sprintf(
		"**Discord ID**: ||%s||\n"+
			"**Discord Code**: ||%s||\n"+
			"**Platform**: %s\n"+
			"**Gamertag**: ||%s||\n"+
			"**XUID**: ||%s||\n"+
			"**Clan Tag**: %s\n"+
			"**Crew ID**: %s\n"+
			"**Session Key**: ||%s||\n"+
			"**Session Ticket**: ||%s||\n"+
			"**Link Status**: %s\n"+
			"**Last Online**: %s\n"+
			"**Ban Status**: %s\n\n"+
			"%s\n%s",

		utils.NullStr(discordID),
		utils.NullStr(discordCode),
		utils.NullStr(platformName, "Unknown"),
		playerGamertag,
		utils.NullStr(xuid, "Not Available"),
		utils.NullStr(crewTag, "Not Set"),
		utils.NullStr(crewID, "Not Set"),
		utils.Truncate(sessionKey),
		utils.Truncate(sessionTicket),

		map[bool]string{
			true:  "Linked",
			false: "Not Linked",
		}[linkdiscord == 1],

		utils.FormatTime(lastOnline),

		map[bool]string{
			true:  "Banned",
			false: "Not Banned",
		}[isBanned], utils.OwnedCrewText(hasOwnedCrew, crewName, crewTag,), utils.MemberCrewText(crewTag, crewID,),)

	color := 0xff0000

	if hasOwnedCrew {
		color = utils.ParseColor(crewColor)
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Player Information: %s", playerGamertag,),
		Description: desc,
		Color: color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Player Info Command",
		},
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed,)
	if err != nil {
		fmt.Println(err)
	}
}