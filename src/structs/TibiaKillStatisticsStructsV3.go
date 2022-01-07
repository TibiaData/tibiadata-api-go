package structs

// KillStatistics stores a world's kill statistics
type KillStatistics struct {
	World   string  `json:"world"`
	Entries []Entry `json:"entries"`
	Total   Total   `json:"total"`
}

// Entry is a kill statistic entry
type Entry struct {
	Race                    string `json:"race"`
	LastDayKilledPlayers    int    `json:"last_day_players_killed"`
	LastDayKilledByPlayers  int    `json:"last_day_killed"`
	LastWeekKilledPlayers   int    `json:"last_week_players_killed"`
	LastWeekKilledByPlayers int    `json:"last_week_killed"`
}

// Total is a kill statistic total amount
type Total struct {
	LastDayKilledPlayers    int `json:"last_day_players_killed"`
	LastDayKilledByPlayers  int `json:"last_day_killed"`
	LastWeekKilledPlayers   int `json:"last_week_players_killed"`
	LastWeekKilledByPlayers int `json:"last_week_killed"`
}
