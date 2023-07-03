package validation

const (
	MaxRunesAllowedInACharacterName     = 29 // Larggest character name length possible
	MinRunesAllowedInACharacterName     = 2  // Smallest character name length possible
	MaxRunesAllowedInACharacterNameWord = 16 // Larggest character name word length possible (new names only 14, but older ones are longer)
	MinRunesAllowedInACharacterNameWord = 2  // Smallest character name word length possible

	MaxRunesAllowedInAGuildName     = 29 // Larggest guild name length possible
	MinRunesAllowedInAGuildName     = 3  // Smallest guild name length possible
	MaxRunesAllowedInAGuildNameWord = 14 // Larggest guild name word length possible
	MinRunesAllowedInAGuildNameWord = 2  // Smallest character name word length possible

	// AmountOfBoostableBosses is the amount of boostable bosses.
	// Last updated: Jun 30 2023
	AmountOfBoostableBosses = 91
)
