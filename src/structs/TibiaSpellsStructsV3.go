package structs

// SpellsOverview stores a list of spells
type SpellsOverview struct {
	SpellsVocationFilter string          `json:"spells_filter"`
	Spells               []SpellOverview `json:"spell_list"`
}

// SpellOverview stores information of a specific spell
type SpellOverview struct {
	Name         string `json:"name"`
	Spell        string `json:"spell_id"`
	Formula      string `json:"formula"`
	Level        int    `json:"level"`
	Mana         int    `json:"mana"`
	Price        int    `json:"price"`
	GroupAttack  bool   `json:"group_attack"`
	GroupHealing bool   `json:"group_healing"`
	GroupSupport bool   `json:"group_support"`
	TypeInstant  bool   `json:"type_instant"`
	TypeRune     bool   `json:"type_rune"`
	PremiumOnly  bool   `json:"premium_only"`
}

// Spell sotres a specific spell data
type Spell struct {
	Name                string           `json:"name"`
	Spell               string           `json:"spell_id"`
	ImageURL            string           `json:"image_url"`
	Description         string           `json:"description"`
	HasSpellInformation bool             `json:"has_spell_information"`
	SpellInformation    SpellInformation `json:"spell_information"`
	HasRuneInformation  bool             `json:"has_rune_information"`
	RuneInformation     RuneInformation  `json:"rune_information"`
}

// SpellInformation stores a specific spell information
type SpellInformation struct {
	Formula       string   `json:"formula"`
	Vocation      []string `json:"vocation"`
	GroupAttack   bool     `json:"group_attack"`
	GroupHealing  bool     `json:"group_healing"`
	GroupSupport  bool     `json:"group_support"`
	TypeInstant   bool     `json:"type_instant"`
	TypeRune      bool     `json:"type_rune"`
	DamageType    string   `json:"damage_type"`
	CooldownAlone int      `json:"cooldown_alone"`
	CooldownGroup int      `json:"cooldown_group"`
	SoulPoints    int      `json:"soul_points"`
	Amount        int      `json:"amount"`
	Level         int      `json:"level"`
	Mana          int      `json:"mana"`
	Price         int      `json:"price"`
	City          []string `json:"city"`
	Premium       bool     `json:"premium_only"`
}

// RuneInformation stores a specific rune information
type RuneInformation struct {
	Vocation     []string `json:"vocation"`
	GroupAttack  bool     `json:"group_attack"`
	GroupHealing bool     `json:"group_healing"`
	GroupSupport bool     `json:"group_support"`
	DamageType   string   `json:"damage_type"`
	Level        int      `json:"level"`
	MagicLevel   int      `json:"magic_level"`
}
