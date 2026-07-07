package utils

type Member struct {
	DiscordID   string `db:"discordid"`
	Gamertag    string `db:"gamertag"`
	CrewTag     string `db:"crew_tag"`
	CrewID      string `db:"crew_id"`
	LinkDiscord int    `db:"linkdiscord"`
	LastOnline  string `db:"last_online"`
	Platform    string `db:"platform_name"`
	Banned      int    `db:"banned"`
}

type Crew struct {
	CrewID     string `db:"crew_id"`
	Name       string `db:"crew_name"`
	Tag        string `db:"crew_tag"`
	Owner      string `db:"crew_owner"`
	Color      string `db:"crew_color"`
	Motto      string `db:"crew_motto"`
	Public     int    `db:"crew_public"`
	InviteCode string `db:"crew_invite"`
}

type CrewWithStats struct {
	Crew
	MemberCount int `db:"member_count"`
}

type OnlinePlayer struct {
	Gamertag string `db:"gamertag"`
	Platform string `db:"platform_name"`
}

type CrewJoinResult struct {
	CrewID string
	Tag    string
	Name   string
	Color  string
}

type CrewInfo struct {
	Name       string
	Tag        string
	Owner      string
	Color      string
	Motto      string
	Public     bool
	InviteCode string
}

type PaginatedResult[T any] struct {
	Data       []T
	Total      int
	TotalPages int
	Page       int
	PageSize   int
}

type InteractionContext struct {
	UserID    string
	GuildID   string
	ChannelID string
	MessageID string
}

type CrewCreateInput struct {
	Name   string
	Tag    string
	Color  string
	Motto  string
	Public bool
}

type CrewUpdateInput struct {
	Name   *string
	Tag    *string
	Color  *string
	Motto  *string
	Public *bool
}

type PlatformCount struct {
	PlatformName string
	PlayerCount  int
}

type ActionResult struct {
	Success bool
	Message string
	Data    any
}

type InviteLookup struct {
	CrewID string
	Name   string
	Tag    string
	Color  string
}

type CrewFull struct {
	Crew        Crew
	OwnerName   string
	MemberCount int
	Visibility  string
}
