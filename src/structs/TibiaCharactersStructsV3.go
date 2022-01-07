package structs

// Characters stores a character and account profile information
type Characters struct {
	Character          Character          `json:"character"`
	AccountBadges      []AccountBadges    `json:"account_badges"`
	Achievements       []Achievements     `json:"achievements"`
	Deaths             Deaths             `json:"deaths"`
	AccountInformation AccountInformation `json:"account_information"`
	OtherCharacters    []OtherCharacters  `json:"other_characters"`
}

// Character stores a tibia character information
type Character struct {
	Name              string            `json:"name"`
	FormerNames       []string          `json:"former_names"`
	Traded            bool              `json:"traded"`
	DeletionDate      string            `json:"deletion_date"`
	Sex               string            `json:"sex"`
	Title             string            `json:"title"`
	UnlockedTitles    int               `json:"unlocked_titles"`
	Vocation          string            `json:"vocation"`
	Level             int               `json:"level"`
	AchievementPoints int               `json:"achievement_points"`
	World             string            `json:"world"`
	FormerWorlds      []string          `json:"former_worlds"`
	Residence         string            `json:"residence"`
	MarriedTo         string            `json:"married_to"`
	Houses            []CharacterHouses `json:"houses"`
	Guild             CharacterGuild    `json:"guild"`
	LastLogin         string            `json:"last_login"`
	AccountStatus     string            `json:"account_status"`
	Comment           string            `json:"comment"`
}

// Houses stores houses owned by a specific character
type CharacterHouses struct {
	Name    string `json:"name"`
	Town    string `json:"town"`
	Paid    string `json:"paid"`
	HouseID int    `json:"houseid"`
}

// Guild stores a character's guild information
type CharacterGuild struct {
	GuildName string `json:"name"`
	Rank      string `json:"rank"`
}

// AccountBadges stores accounts badges information
type AccountBadges struct {
	Name        string `json:"name"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
}

// Achievements stores accounts achievements information
type Achievements struct {
	Name   string `json:"name"`
	Grade  int    `json:"grade"`
	Secret bool   `json:"secret"`
}

// Deaths stores all deaths from a character
type Deaths struct {
	DeathEntries    []DeathEntries `json:"death_entries"`
	TruncatedDeaths bool           `json:"truncated"` // deathlist can be truncated.. but we don't have logic for that atm
}

// DeathEntries stores a character's specific death information
type DeathEntries struct {
	Time    string    `json:"time"`
	Level   int       `json:"level"`
	Killers []Killers `json:"killers"`
	Assists []Killers `json:"assists"`
	Reason  string    `json:"reason"`
}

// Killers stores a character's killers from a specific death
type Killers struct {
	Name   string `json:"name"`
	Player bool   `json:"player"`
	Traded bool   `json:"traded"`
	Summon string `json:"summon"`
}

// AccountInformation stores some account specific information
type AccountInformation struct {
	Position     string `json:"position"`
	Created      string `json:"created"`
	LoyaltyTitle string `json:"loyalty_title"`
}

// OtherCharacters stores others characters from an account
type OtherCharacters struct {
	Name    string `json:"name"`
	World   string `json:"world"`
	Status  string `json:"status"`  // online/offline
	Deleted bool   `json:"deleted"` // don't know how to do that yet..
	Main    bool   `json:"main"`
	Traded  bool   `json:"traded"`
}
