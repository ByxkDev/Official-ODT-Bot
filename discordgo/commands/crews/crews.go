package crews

import (
	"fmt"
	"strings"
	"time"

	"discordgo/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	DefaultCrewTag = "ROS"
	DefaultCrewID  = 1337
	MaxCrewTag     = 4
	MaxMotto       = 30
)

func CrewCode(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	userID := m.Author.ID
	var crewID, crewTag, crewOwner string

	err := utils.DB().QueryRow(`SELECT crew_id, crew_tag FROM members WHERE discordid = ?`, userID).Scan(&crewID, &crewTag)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "You are not in a crew.", 0xff0000)
		return
	}

	err = utils.DB().QueryRow(`SELECT crew_owner FROM crews WHERE crew_id = ?`, crewID).Scan(&crewOwner)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Crew not found.", 0xff0000)
		return
	}

	if crewOwner != userID {
		utils.Send(s, m.ChannelID, "Error", "Only the crew owner can generate invite codes.", 0xff0000)
		return
	}

	code := utils.GenerateCode()
	_, err = utils.DB().Exec(`UPDATE crews SET crew_invite = ? WHERE crew_id = ?`, code, crewID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Failed to generate invite code.", 0xff0000)
		return
	}

	utils.Send(s, m.ChannelID, "Crew Invite Code", fmt.Sprintf("Your crew invite code is: **%s**", code), 0x00ff00)
}

func CrewStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	var linkdiscord, crewPublic, memberCount int
	var crewTag, crewOwner, crewColor, crewInvite string

	err := utils.DB().QueryRow(`SELECT linkdiscord, crew_tag FROM members WHERE discordid = ?`, userID).Scan(&linkdiscord, &crewTag)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "No information found for your account. Please link your account first.", 0xff0000)
		return
	}

	if linkdiscord != 1 {
		utils.Send(s, m.ChannelID, "Error", "You are not linked to a crew system account. Please link your account first.", 0xff0000)
		return
	}

	if crewTag == "" {
		crewTag = DefaultCrewTag
	}

	if crewTag == DefaultCrewTag {
		utils.Send(s, m.ChannelID, "Special Crew Status", "You are in ROS.\nUse:\n?crewcreate\n?crewjoin\n?crewcode <code>", 0xff0000)
		return
	}

	err = utils.DB().QueryRow(`SELECT crew_owner, crew_color, crew_public, crew_invite FROM crews WHERE crew_tag = ?`, crewTag).Scan(&crewOwner, &crewColor, &crewPublic, &crewInvite)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Crew not found. You have been reset to ROS.", 0xff0000)
		_, _ = utils.DB().Exec(`UPDATE members SET crew_tag = 'ROS' WHERE discordid = ?`, userID)
		return
	}

	_ = utils.DB().QueryRow(`SELECT COUNT(*) FROM members WHERE crew_tag = ?`, crewTag).Scan(&memberCount)

	isOwner := crewOwner == userID
	ownerName := "Unknown"
	if isOwner {
		ownerName = "You (Crew Owner)"
	} else if u, err := s.User(crewOwner); err == nil {
		ownerName = u.Username
	}

	inviteText := ""
	if isOwner && crewPublic == 0 {
		inviteText = fmt.Sprintf("\nInvite Code: %s", crewInvite)
	}

	utils.Send(s, m.ChannelID, "Your Crew", fmt.Sprintf("Tag: %s\nOwner: %s\nMembers: %d%s", crewTag, ownerName, memberCount, inviteText), utils.ParseColor(crewColor))
}

func CrewInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	var gamertag, crewTag string

	err := utils.DB().QueryRow(`SELECT gamertag, crew_tag FROM members WHERE discordid = ?`, userID).Scan(&gamertag, &crewTag)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "No data found.", 0xff0000)
		return
	}

	if crewTag == "" || crewTag == DefaultCrewTag {
		crewTag = "No Crew"
	}

	utils.Send(s, m.ChannelID, "User Info", fmt.Sprintf("Gamertag: %s\nCrew: %s", gamertag, crewTag), 0x00ffcc)
}

func CrewCreate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	userID := m.Author.ID

	if len(args) < 5 {
		utils.Send(s, m.ChannelID, "Usage", "?crewcreate <name> <tag> <color> <motto...> <yes/no>", 0xffaa00)
		return
	}

	name := args[0]
	tag := strings.ToUpper(args[1])
	color := args[2]
	public := strings.ToLower(args[len(args)-1])
	motto := strings.Join(args[3:len(args)-1], " ")
	crewID := fmt.Sprintf("%d", time.Now().UnixNano())

	if len(tag) > MaxCrewTag {
		utils.Send(s, m.ChannelID, "Error", "Tag max 4 characters.", 0xff0000)
		return
	}

	if len(motto) > MaxMotto {
		utils.Send(s, m.ChannelID, "Error", "Motto max 30 characters.", 0xff0000)
		return
	}

	if utils.IsInCrew(userID) {
		utils.Send(s, m.ChannelID, "Error", "Leave your current crew first.", 0xff0000)
		return
	}

	invite := ""

	if public == "no" {
		invite = utils.GenerateCode()
	}

	_, err := utils.DB().Exec(`INSERT INTO crews (crew_id, crew_name, crew_tag, crew_owner, crew_color, crew_motto, crew_public, crew_invite) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, crewID, name, tag, userID, color, motto, public == "yes", invite)
	if err != nil {
		fmt.Println("[CrewCreate ERROR]", err)
		utils.Send(s, m.ChannelID, "Error", "Failed to create crew.", 0xff0000)
		return
	}

	_, _ = utils.DB().Exec(`UPDATE members SET crew_tag = ?, crew_id = ? WHERE discordid = ?`, tag, crewID, userID)
	utils.Send(s, m.ChannelID, "Crew Created", fmt.Sprintf("Created **%s (%s)**", name, tag), 0x00ff00)
}

func CrewJoin(s *discordgo.Session, m *discordgo.MessageCreate, code string) {
	userID := m.Author.ID
	var crewID, tag, name string

	if utils.IsInCrew(userID) {
		utils.Send(s, m.ChannelID, "Error", "Already in a crew.", 0xff0000)
		return
	}

	err := utils.DB().QueryRow(`SELECT crew_id, crew_tag, crew_name FROM crews WHERE crew_invite = ?`, code).Scan(&crewID, &tag, &name)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Invalid invite.", 0xff0000)
		return
	}

	_, err = utils.DB().Exec(`UPDATE members SET crew_tag = ?, crew_id = ? WHERE discordid = ?`, tag, crewID, userID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Database error.", 0xff0000)
		return
	}

	utils.Send(s, m.ChannelID, "Joined Crew", fmt.Sprintf("You joined **%s (%s)**", name, tag), 0x00ff00)
}

func CrewLeave(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	var tag, crewID string

	err := utils.DB().QueryRow(`SELECT crew_tag, crew_id FROM members WHERE discordid = ?`, userID).Scan(&tag, &crewID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Profile not found.", 0xff0000)
		return
	}

	if tag == DefaultCrewTag || tag == "" {
		utils.Send(s, m.ChannelID, "Error", "You are not in a real crew.", 0xff0000)
		return
	}

	_, err = utils.DB().Exec(`UPDATE members SET crew_tag = ?, crew_id = ? WHERE discordid = ?`, DefaultCrewTag, DefaultCrewID, userID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Failed to leave crew.", 0xff0000)
		return
	}

	utils.Send(s, m.ChannelID, "Left Crew", "You successfully left your crew and returned to ROS.", 0x00ff00)
}

func CrewDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	var crewID, crewTag, crewOwner string

	err := utils.DB().QueryRow(`SELECT crew_id, crew_tag, crew_owner FROM crews WHERE crew_owner = ?`, userID).Scan(&crewID, &crewTag, &crewOwner)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "You don't own a crew.", 0xff0000)
		return
	}

	_, err = utils.DB().Exec(`UPDATE members SET crew_tag = ?, crew_id = ? WHERE crew_id = ?`, DefaultCrewTag, DefaultCrewID, crewID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Failed to reset crew members.", 0xff0000)
		return
	}

	_, err = utils.DB().Exec(`DELETE FROM crews WHERE crew_id = ?`, crewID)
	if err != nil {
		utils.Send(s, m.ChannelID, "Error", "Failed to delete crew.", 0xff0000)
		return
	}

	_, _ = utils.DB().Exec(`UPDATE members SET crew_tag = ?, crew_id = ? WHERE discordid = ?`, DefaultCrewTag, DefaultCrewID, userID)
	utils.Send(s, m.ChannelID, "Crew Deleted", fmt.Sprintf("Crew **%s (%s)** has been deleted and members reset to ROS.", crewTag, crewID), 0xff0000)
}