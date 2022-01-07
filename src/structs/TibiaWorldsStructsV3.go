package structs

// Worlds stores a list of worlds with some information
type Worlds struct {
	PlayersOnline    int     `json:"players_online"` // Calculated value
	RecordPlayers    int     `json:"record_players"` // Overall Maximum:
	RecordDate       string  `json:"record_date"`    // Overall Maximum:
	RegularWorlds    []World `json:"regular_worlds"`
	TournamentWorlds []World `json:"tournament_worlds"`
}

// World stores a specific world information
// Notice that on the overview endpoint (.../worlds) it will be
// partially filled
type World struct {
	Name                string         `json:"name"`
	Status              string         `json:"status"`                       // Status:
	PlayersOnline       int            `json:"players_online"`               // Players Online:
	RecordPlayers       int            `json:"record_players,omitempty"`     // Online Record:
	RecordDate          string         `json:"record_date,omitempty"`        // Online Record:
	CreationDate        string         `json:"creation_date,omitempty"`      // Creation Date: -> convert to YYYY-MM
	Location            string         `json:"location"`                     // Location:
	PvpType             string         `json:"pvp_type"`                     // PvP Type:
	PremiumOnly         bool           `json:"premium_only"`                 // Premium Type: premium = true / else: false
	TransferType        string         `json:"transfer_type"`                // Transfer Type: regular (if not present) / locked / blocked
	WorldsQuestTitles   []string       `json:"world_quest_titles,omitempty"` // World Quest Titles:
	BattleyeProtected   bool           `json:"battleye_protected"`           // BattlEye Status: true if protected / false if "Not protected by BattlEye."
	BattleyeDate        string         `json:"battleye_date"`                // BattlEye Status: null if since release / else show date?
	GameWorldType       string         `json:"game_world_type"`              // Game World Type: regular / experimental / tournament (if Tournament World Type exists)
	TournamentWorldType string         `json:"tournament_world_type"`        // Tournament World Type: "" (default?) / regular / restricted
	OnlinePlayers       []OnlinePlayer `json:"online_players,omitempty"`
}

// OnlinePLayers stores a world's online player information
type OnlinePlayer struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
}
