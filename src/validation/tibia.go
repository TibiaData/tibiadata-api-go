package validation

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// characterNameRegex is used to check if the character name provided is valid
	// nowadays Tibia only accepts a-zA-Z, but we have to consider old names
	characterNameRegex = regexp.MustCompile(`[^\s'\p{L}\-\.\+]`)

	// creatureAndSpellNameRegex is used to check if the creature or spell name provided is valid
	creatureAndSpellNameRegex = regexp.MustCompile(`[^\s'a-zA-Z]`)

	// guildNameRegex is used to check if the guild name provided is valid
	guildNameRegex = regexp.MustCompile(`[^\sa-zA-Z]`)

	// validVocations stores all valid tibia vocations
	validVocations = []string{"none", "knight", "knights", "paladin", "paladins", "sorcerer", "sorcerers", "druid", "druids", "all"}
)

// IsNewsIDValid reports wheter the provided int represents a valid news ID
// Check if error == nil to see whether the ID is valid or not
func IsNewsIDValid(ID int) error {
	// If the ID is 0 or a negative number it is invalid
	if ID <= 0 {
		return ErrorInvalidNewsID
	}

	return nil
}

// IsVocationValid reports wheter the provided string represents a valid vocation
// Check if error == nil to see whether the vocation is valid or not
func IsVocationValid(vocation string) error {
	for _, voc := range validVocations {
		if strings.EqualFold(vocation, voc) {
			return nil
		}
	}

	return ErrorVocationDoesNotExist
}

// IsCharacterNameValid reports wheter the provided string represents a valid character name
// Check if error == nil to see whether the name is valid or not
func IsCharacterNameValid(name string) error {
	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return ErrorCharacterNameEmpty
	case lenName < MinRunesAllowedInACharacterName: // Name is too small
		return ErrorCharacterNameTooSmall
	case lenName > MaxRunesAllowedInACharacterName: // Name is too big
		return ErrorCharacterNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return ErrorCharacterNameIsOnlyWhiteSpace
	}

	// Check if any word in the name has a length > 14
	strs := strings.Fields(name)
	for _, str := range strs {
		if utf8.RuneCountInString(str) > MaxRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooBig
		}

		if utf8.RuneCountInString(str) < MinRunesAllowedInACharacterNameWord {
			return ErrorCharacterWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := characterNameRegex.MatchString(name)
	if matched {
		return ErrorCharacterNameInvalid
	}

	return nil
}

// IsGuildNameValid reports wheter the provided string represents a valid guild name
// Check if error == nil to see whether the name is valid or not
func IsGuildNameValid(name string) error {
	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return ErrorGuildNameEmpty
	case lenName < MinRunesAllowedInAGuildName: // Name is too small
		return ErrorGuildNameTooSmall
	case lenName > MaxRunesAllowedInAGuildName: // Name is too big
		return ErrorGuildNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return ErrorGuildNameIsOnlyWhiteSpace
	}

	// Check if any word in the name has a length > 14
	strs := strings.Fields(name)
	for _, str := range strs {
		if utf8.RuneCountInString(str) > MaxRunesAllowedInAGuildNameWord {
			return ErrorGuildWordTooBig
		}

		if utf8.RuneCountInString(str) < MinRunesAllowedInAGuildNameWord {
			return ErrorGuildWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := guildNameRegex.MatchString(name)
	if matched {
		return ErrorGuildNameInvalid
	}

	return nil
}

// IsCreatureNameValid reports wheter the provided string represents a valid creature name
// Check if error == nil to see whether the creature is valid or not
// It will also return the creature endpoint
func IsCreatureNameValid(name string) (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return "", ErrorCreatureNameEmpty
	case lenName < smallestCreatureNameRuneCount: // Name is too small
		return "", ErrorCreatureNameTooSmall
	case lenName > biggestCreatureNameRuneCount: // Name is too big
		return "", ErrorCreatureNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return "", ErrorCreatureNameIsOnlyWhiteSpace
	}

	// Check words length
	strs := strings.Fields(name)
	for _, str := range strs {
		utfCount := utf8.RuneCountInString(str)

		if utfCount > biggestCreatureWordRuneCount {
			return "", ErrorCreatureWordTooBig
		}

		if utfCount < smallestCreatureWordRuneCount {
			return "", ErrorCreatureWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := creatureAndSpellNameRegex.MatchString(name)
	if matched {
		return "", ErrorCreatureNameInvalid
	}

	var (
		found    bool
		endpoint string
	)

	// Check if creature exists
	for _, creature := range val.Creatures {
		if strings.EqualFold(name, creature.Endpoint) || strings.EqualFold(name, creature.Name) || strings.EqualFold(name, creature.PluralName) {
			found = true
			endpoint = creature.Endpoint
			break
		}
	}

	if !found {
		return "", ErrorCreatureNotFound
	}

	return endpoint, nil
}

// IsSpellNameOrFormulaValid reports wheter the provided string represents a valid spell name or formula
// Check if error == nil to see whether the creature is valid or not
// It will also return the spell endpoint
func IsSpellNameOrFormulaValid(name string) (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	// Getting the length of the name
	lenName := utf8.RuneCountInString(name)

	switch {
	case lenName == 0: // Name is an empty string
		return "", ErrorSpellNameEmpty
	case lenName < smallestSpellNameOrFormulaRuneCount: // Name is too small
		return "", ErrorSpellNameTooSmall
	case lenName > biggestSpellNameOrFormulaRuneCount: // Name is too big
		return "", ErrorSpellNameTooBig
	}

	// Check if name consists of whitespaces only
	if strings.TrimSpace(name) == "" {
		return "", ErrorSpellNameIsOnlyWhiteSpace
	}

	// Check words length
	strs := strings.Fields(name)
	for _, str := range strs {
		utfCount := utf8.RuneCountInString(str)

		if utfCount > biggestSpellWordRuneCount {
			return "", ErrorSpellWordTooBig
		}

		if utfCount < smallestSpellWordRuneCount {
			return "", ErrorSpellWordTooSmall
		}
	}

	// Check if name matches the regex
	matched := creatureAndSpellNameRegex.MatchString(name)
	if matched {
		return "", ErrorSpellNameInvalid
	}

	var (
		found    bool
		endpoint string
	)

	// Check if spell exists
	for _, spell := range val.Spells {
		if strings.EqualFold(name, spell.Endpoint) || strings.EqualFold(name, spell.Name) || strings.EqualFold(name, spell.Formula) {
			found = true
			endpoint = spell.Endpoint
			break
		}
	}

	if !found {
		return "", ErrorSpellNotFound
	}

	return endpoint, nil
}

// GetWorlds returns a list of all existing worlds
func GetWorlds() ([]string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Worlds, nil
}

// WorldExists reports whether the specified world exists
// This function is case insensitive
func WorldExists(world string) (bool, error) {
	// Check if the validator has been initiated
	if !initiated {
		return false, ErrorValidatorNotInitiated
	}

	// Try to find the world
	for _, w := range val.Worlds {
		if strings.EqualFold(w, world) {
			return true, nil
		}
	}

	return false, nil
}

// GetTowns returns a list of all existing towns
func GetTowns() ([]string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Towns, nil
}

// TowndExists reports whether the specified town exists
// This function is case insensitive
func TownExists(town string) (bool, error) {
	// Check if the validator has been initiated
	if !initiated {
		return false, ErrorValidatorNotInitiated
	}

	// Try to find the town
	for _, t := range val.Towns {
		if strings.EqualFold(t, town) {
			return true, nil
		}
	}

	return false, nil
}

// GetHouses returns a slice of all houses
func GetHouses() ([]House, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Houses, nil
}

// GetHouseRaw returns a house by it's ID, independently
// of what town the house is from
// This function will return a nil house AND a nil error
// if the specified ID doesn't exist
func GetHouseRaw(houseID int) (*House, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	// Try to find the house
	for _, h := range val.Houses {
		if h.ID == houseID {
			return &h, nil
		}
	}

	return nil, nil
}

// HouseExistsRaw reports whether a house exits, independently
// of what town the house is from
func HouseExistsRaw(houseID int) (bool, error) {
	house, err := GetHouseRaw(houseID)
	return house != nil, err
}

// GetHouseInTown returns a house by it's ID and town
// This function will return a nil house AND a nil error
// if the specified ID doesn't exist in the specified town
// or if the specified town doesn't exist
func GetHouseInTown(houseID int, town string) (*House, error) {
	// We don't need to check if the validator has been initiated
	// because TownExists will already check that for us
	townExists, err := TownExists(town)
	if err != nil {
		return nil, err
	}

	// Town doesn't exist
	if !townExists {
		return nil, nil
	}

	// Try to find the house
	for _, h := range val.Houses {
		if h.ID == houseID && strings.EqualFold(h.Town, town) {
			return &h, nil
		}
	}

	return nil, nil
}

// HouseExistsInTown reports whether a house exits in the specified town
// This function will return false AND a nil error if the specified ID
// doesn't exist in the specified town or if the specified town doesn't exist
func HouseExistsInTown(houseID int, town string) (bool, error) {
	house, err := GetHouseInTown(houseID, town)
	return house != nil, err
}

// GetCreatures returns a list of all existing creatures
func GetCreatures() ([]Creature, error) {
	// Check if the validator has been initiated
	if !initiated {
		return nil, ErrorValidatorNotInitiated
	}

	return val.Creatures, nil
}
