package discord

import (
	"fmt"
	"os"
	"strings"
	"time"

	"discordgo/commands/crews"
	"discordgo/commands/early"
	"discordgo/commands/information"
	"discordgo/commands/moderation"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

var (
	CartelBoss     string
	CartelCoBoss   string
	CartelAdmin    string
	CartelStaff    string
	EarlySupporter string
	EventWinner    string
)

var startTime = time.Now()

func loadRoles() {
	CartelBoss = os.Getenv("CARTEL_BOSSID")
	CartelCoBoss = os.Getenv("CARTEL_COBOSSID")
	CartelAdmin = os.Getenv("CARTEL_ADMINID")
	CartelStaff = os.Getenv("CARTEL_STAFFID")
	EarlySupporter = os.Getenv("EARLY_SUPPORTERID")
	EventWinner = os.Getenv("EVENT_WINNERID")
}

func hasStaffRole(member *discordgo.Member) bool {

	if member == nil {
		return false
	}

	for _, role := range member.Roles {
		switch role {
		case CartelBoss,
			CartelCoBoss,
			CartelAdmin,
			CartelStaff:
			return true
		}
	}

	return false
}


func requireStaffRole(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if hasStaffRole(m.Member) {
		return true
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" You don't have permission to use this command.",)
	return false
}

func allCommandsHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendEmbed(s, m, "All Bot Commands", "**Public Commands**\n"+ "`?help` - Show public help\n"+ "`?claim` - Claim early supporter role\n"+ "`?link` -  important links\n"+ "`?verify` - Verify yourself\n\n"+ "**Crew Commands**\n"+ "`?crewhelp` - Show crew help\n"+ "`?crewstatus` - View crew status\n"+ "`?crewinfo` - View crew information\n"+ "`?crewcreate <name>` - Create a crew\n"+ "`?crewcode` - Get crew code\n"+ "`?crewjoin <code>` - Join a crew\n"+ "`?crewleave` - Leave crew\n"+ "`?crewdelete` - Delete crew\n\n"+ "**Staff Commands**\n"+ "`?staffhelp` - Show staff help\n"+ "`?crewemblem` - Manage crew emblems\n"+ "`?crewlist` - List crews\n"+ "`?info` - Server information\n"+ "`?playercount` - Show player count\n"+ "`?playerinfo <player>` - Lookup player",)
}

func hasCrewRole(member *discordgo.Member) bool {
	if member == nil {
		return false
	}

	for _, role := range member.Roles {
		switch role {
		case CartelBoss,
			CartelCoBoss,
			CartelAdmin,
			CartelStaff,
			EarlySupporter,
			EventWinner:
			return true
		}
	}

	return false
}


func requireCrewRole(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if hasCrewRole(m.Member) {
		return true
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"You don't have permission to use this command.",)
	return false
}

func publicHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendEmbed(s, m, "Public Commands", "`?claim` - Claim early supporter role\n"+ "`?link` - Link gamertag to Discord\n"+ "`?verify` - Verify yourself\n"+ "`?help` - Show public commands",)
}

func crewHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendEmbed(s, m, "Crew Commands", "`?crewstatus` - View your crew status\n"+ "`?crewinfo` - View crew information\n"+ "`?crewcreate <name>` - Create a crew\n"+ "`?crewcode` - Get crew code\n"+ "`?crewjoin <code>` - Join a crew\n"+ "`?crewleave` - Leave your crew\n"+ "`?crewdelete` - Delete your crew\n"+ "`?crewhelp` - Show crew commands",)
}

func staffHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendEmbed(s, m, "Staff Commands", "`?crewemblem` - Manage crew emblems\n"+ "`?crewlist` - List crews\n"+ "`?info` - Server information\n"+ "`?playercount` - Show player count\n"+ "`?playerinfo` - Player lookup\n"+ "`?staffhelp` - Show staff commands",)
}

func New(token string) (*Bot, error) {
	loadRoles()

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent
	bot := &Bot{
		Session: session,
	}

	session.AddHandler(bot.MessageCreate)
	return bot, nil
}


func (b *Bot) Start() error {
	return b.Session.Open()
}

func sendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, title string, description string,) {
	name := m.Author.Username

	if m.Member != nil && m.Member.Nick != "" {
		name = m.Member.Nick
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Description: m.Author.Mention() + "\n\n" + description,
		Color: 0x3498DB,
		Author: &discordgo.MessageEmbedAuthor{
			Name: name,
			IconURL: m.Author.AvatarURL(""),
		},

		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed,)
	if err != nil {
		fmt.Println(err)
	}

}

func (b *Bot) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate,) {
	if m.Author.Bot {
		return
	}

	content := strings.TrimSpace(m.Content)
	if !strings.HasPrefix(content, "?") {
		return
	}

	_ = s.ChannelMessageDelete(m.ChannelID, m.ID,)

	switch {

	//PUBLIC

	case content == "?help":
		publicHelp(s, m)

	case content == "?claim":
		early.Claim(s, m)

	case strings.HasPrefix(content, "?link"):
		information.Link(s, m)

	case strings.HasPrefix(content, "?verify"):
		information.Verify(s, m)

	//CREW

	case content == "?crewhelp":
		if !requireCrewRole(s, m) {
			return
		}
		crewHelp(s, m)

	case strings.HasPrefix(content, "?crewstatus"):
		if !requireCrewRole(s, m) {
			return
		}
		crews.CrewStatus(s, m)

	case content == "?crewinfo":
		if !requireCrewRole(s, m) {
			return
		}
		crews.CrewInfo(s, m)

	case strings.HasPrefix(content, "?crewcreate"):
		if !requireCrewRole(s, m) {
			return
		}
		parts := strings.Fields(content)
		if len(parts) < 2 {
			sendEmbed(s, m, "Error", "Usage: `?crewcreate <name>`")
			return
		}
		crews.CrewCreate(s, m, parts[1:])

	case strings.HasPrefix(content, "?crewcode"):
		if !requireCrewRole(s, m) {
			return
		}
		parts := strings.Fields(content)
		crews.CrewCode(s, m, parts)

	case strings.HasPrefix(content, "?crewjoin"):
		if !requireCrewRole(s, m) {
			return
		}
		code := strings.TrimSpace(strings.TrimPrefix(content, "?crewjoin"),)
		crews.CrewJoin(s, m, code)

	case content == "?crewleave":
		if !requireCrewRole(s, m) {
			return
		}
		crews.CrewLeave(s, m)

	case content == "?crewdelete":
		if !requireCrewRole(s, m) {
			return
		}
		crews.CrewDelete(s, m)

	//STAFF

    case content == "?allcommands":
		if !requireStaffRole(s, m) {
			return
		}
	    allCommandsHelp(s, m)

	case content == "?staffhelp":
		if !requireStaffRole(s, m) {
			return
		}
		staffHelp(s, m)

	case strings.HasPrefix(content, "?crewemblem"):
		if !requireStaffRole(s, m) {
			return
		}
		crews.Emblem(s, m)

	case content == "?crewlist":
		if !requireStaffRole(s, m) {
			return
		}
		crews.List(s, m)

	case content == "?info":
		if !requireStaffRole(s, m) {
			return
		}
		information.Information(s, m)

	case content == "?playercount":
		if !requireStaffRole(s, m) {
			return
		}
		information.PlayerCount(s, m)

	case strings.HasPrefix(content, "?playerinfo"):
		if !requireStaffRole(s, m) {
			return
		}
		moderation.PlayerInfo(s, m)
	}

}