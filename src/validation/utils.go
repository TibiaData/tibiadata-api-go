package validation

import "unicode"

// GetSha256Sum returns the sha256sum of the data.min.json file being used
func GetSha256Sum() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return sha256sum, nil
}

// GetSha512Sum returns the sha512sum of the data.min.json file being used
func GetSha512Sum() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return sha512sum, nil
}

// DoesStringContainDigits returns whether there is a digit rune in the string
func DoesStringContainDigits(str string) bool {
	for _, s := range str {
		if unicode.IsDigit(s) {
			return true
		}
	}

	return false
}

// DoesStringContainsNumbers returns whether there is a number rune in the string
func DoesStringContainsNumbers(str string) bool {
	for _, s := range str {
		if unicode.IsNumber(s) {
			return true
		}
	}

	return false
}

// GetSmallestCreatureName returns the name of the creature with the smallest name
func GetSmallestCreatureName() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return smallestCreatureName, nil
}

// GetBiggestCreatureName returns the name of the creature with the biggest name
func GetBiggestCreatureName() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return biggestCreatureName, nil
}

// GetBiggestCreatureWord returns the biggest word in a creature name
func GetBiggestCreatureWord() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return biggestCreatureWord, nil
}

// GetSmallestCreatureWord returns the smallest word in a creature name
func GetSmallestCreatureWord() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return smallestCreatureWord, nil
}

// GetSmallestCreatureNameRuneCount returns the length of the smallest creature name
func GetSmallestCreatureNameRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return smallestCreatureNameRuneCount, nil
}

// GetBiggestCreatureNameRuneCount returns the length of the biggest creature name
func GetBiggestCreatureNameRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return biggestCreatureNameRuneCount, nil
}

// GetSmallestCreatureWordRuneCount returns the length of the smallest creature word
func GetSmallestCreatureWordRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return smallestCreatureWordRuneCount, nil
}

// GetBiggestCreatureWordRuneCount returns the length of the biggest creature word
func GetBiggestCreatureWordRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return biggestCreatureWordRuneCount, nil
}

// GetSmallestSpellNameOrFormula returns the name of the spell with the smallest name or formula
func GetSmallestSpellNameOrFormula() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return smallestSpellNameOrFormula, nil
}

// GetBiggestSpellNameOrFormula returns the name of the spell with the biggest name or formula
func GetBiggestSpellNameOrFormula() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return biggestSpellNameOrFormula, nil
}

// GetBiggestSpellWord returns the biggest word in a spell name or formula
func GetBiggestSpellWord() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return biggestSpellWord, nil
}

// GetSmallestSpellWord returns the smallest word in a spell name or formula
func GetSmallestSpellWord() (string, error) {
	// Check if the validator has been initiated
	if !initiated {
		return "", ErrorValidatorNotInitiated
	}

	return smallestSpellWord, nil
}

// GetSmallestSpellNameOrFormulaRuneCount returns the length of the smallest spell name
func GetSmallestSpellNameOrFormulaRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return smallestSpellNameOrFormulaRuneCount, nil
}

// GetBiggestSpellNameOrFormulaRuneCount returns the length of the biggest spell name
func GetBiggestSpellNameOrFormulaRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return biggestSpellNameOrFormulaRuneCount, nil
}

// GetSmallestSpellWordRuneCount returns the length of the smallest spell word
func GetSmallestSpellWordRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return smallestSpellWordRuneCount, nil
}

// GetBiggestSpellWordRuneCount returns the length of the biggest spell word
func GetBiggestSpellWordRuneCount() (int, error) {
	// Check if the validator has been initiated
	if !initiated {
		return -1, ErrorValidatorNotInitiated
	}

	return biggestSpellWordRuneCount, nil
}
