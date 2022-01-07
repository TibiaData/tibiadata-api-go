package structs

// Creature stores a creature information
type Creature struct {
	Name             string   `json:"name"`
	Race             string   `json:"race"`
	ImageURL         string   `json:"image_url"`
	Description      string   `json:"description"`
	Behaviour        string   `json:"behaviour"`
	Hitpoints        int      `json:"hitpoints"`
	ImmuneTo         []string `json:"immune"`
	StrongAgainst    []string `json:"strong"`
	WeaknessAgainst  []string `json:"weakness"`
	BeParalysed      bool     `json:"be_paralysed"`
	BeSummoned       bool     `json:"be_summoned"`
	SummonMana       int      `json:"summoned_mana"`
	BeConvinced      bool     `json:"be_convinced"`
	ConvincedMana    int      `json:"convinced_mana"`
	SeeInvisible     bool     `json:"see_invisible"`
	ExperiencePoints int      `json:"experience_points"`
	IsLootable       bool     `json:"is_lootable"`
	LootList         []string `json:"loot_list"`
	Featured         bool     `json:"featured"`
}

// Creatures stores a list of creatures and boosted creatures
type Creatures struct {
	Boosted   Creature   `json:"boosted"`
	Creatures []Creature `json:"creature_list"`
}
