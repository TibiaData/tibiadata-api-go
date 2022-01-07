package structs

// Guilds stores all guilds from a specific world
type Guilds struct {
	World     string  `json:"world"`
	Active    []Guild `json:"active"`
	Formation []Guild `json:"formation"`
}

// Guild stores a guild information
// Notice that on the overview endpoint (.../guilds/world/:world) it will
// only be partially filled (Name, LogoURL and Description)
type Guild struct {
	Name               string         `json:"name"`
	World              string         `json:"world,omitempty"`
	LogoURL            string         `json:"logo_url"`
	Description        string         `json:"description"`
	Guildhalls         []Guildhall    `json:"guildhalls,omitempty"`
	Active             bool           `json:"active,omitempty"`
	Founded            string         `json:"founded,omitempty"`
	Applications       bool           `json:"open_applications,omitempty"`
	Homepage           string         `json:"homepage,omitempty"`
	InWar              bool           `json:"in_war,omitempty"`
	DisbandedDate      string         `json:"disband_date,omitempty"`
	DisbandedCondition string         `json:"disband_condition,omitempty"`
	PlayersOnline      int            `json:"players_online,omitempty"`
	PlayersOffline     int            `json:"players_offline,omitempty"`
	MembersTotal       int            `json:"members_total,omitempty"`
	MembersInvited     int            `json:"members_invited,omitempty"`
	Members            []GuildMembers `json:"members,omitempty"`
	Invited            []GuildInvited `json:"invites,omitempty"`
}

// Guildhall stores a guildhall information
type Guildhall struct {
	Name  string `json:"name"`
	World string `json:"world"` // Maybe duplicate info? Guild can only be on one world..
	/*
		Town      string `json:"town"`       // We can collect that from cached info?
		Status    string `json:"status"`     // rented (but maybe also auctioned)
		Owner     string `json:"owner"`      // We can collect that from cached info?
		HouseID   int    `json:"houseid"`    // We can collect that from cached info?
	*/
	PaidUntil string `json:"paid_until"` // Paid until date
}

// GuildMembers stores the members of a guild
type GuildMembers struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Rank     string `json:"rank"`
	Vocation string `json:"vocation"`
	Level    int    `json:"level"`
	Joined   string `json:"joined"`
	Status   string `json:"status"`
}

// GuildInvited stores characters that have been
// invited to a guild
type GuildInvited struct {
	Name string `json:"name"`
	Date string `json:"date"`
}
