package structs

// Highscores stores a highscore list
type Highscores struct {
	World         string      `json:"world"`
	Category      string      `json:"category"`
	Vocation      string      `json:"vocation"`
	HighscoreAge  int         `json:"highscore_age"`
	HighscoreList []Highscore `json:"highscore_list"`
}

// Highscore stores a specific highscore entry
type Highscore struct {
	Rank     int    `json:"rank"`            // Rank column
	Name     string `json:"name"`            // Name column
	Vocation string `json:"vocation"`        // Vocation column
	World    string `json:"world"`           // World column
	Level    int    `json:"level"`           // Level column
	Value    int    `json:"value"`           // Points/SkillLevel column
	Title    string `json:"title,omitempty"` // Title column (when category: loyalty)
}
