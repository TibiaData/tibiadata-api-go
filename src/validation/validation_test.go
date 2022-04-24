package validation

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRaceCondition(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			_, err := GetHouses()
			if err != nil {
				panic(err)
			}
		}(i)
	}

	wg.Wait()
}

func TestNameValidator(t *testing.T) {
	names := []string{
		"Torbjörn",
		"Himmelhüpferin",
		"Himmili",
		"Elder'Kronus",
		"Dizzy IV",
		"Iceman IV",
		"Iv",
		"Kolskägg",
		"Näurin", "Hälge",
		"Li-Ho Crue",
		"Obi-Bok",
	}

	for _, n := range names {
		err := IsCharacterNameValid(n)
		if err != nil {
			t.Fatalf("Name %s should be valid, but it is not, err: %s", n, err)
		}
	}

	invalidNames := []string{
		"One1",
		"Two 2",
		"A",
		"a",
		"T2wo",
	}

	for _, n := range invalidNames {
		err := IsCharacterNameValid(n)
		if err == nil {
			t.Fatalf("Invalid name passed as valid, name: %s", n)
		}
	}

	jsonErrorNames := map[Error]string{
		ErrorCharacterNameEmpty:            "",
		ErrorCharacterNameTooSmall:         "o",
		ErrorCharacterNameTooBig:           "abcabcabcabcabcabcabcabcabcabc",
		ErrorCharacterWordTooBig:           "abcabcabcabcabc hello",
		ErrorCharacterNameIsOnlyWhiteSpace: "     ",
		ErrorCharacterNameInvalid:          "12",
	}

	for expectedErr, name := range jsonErrorNames {
		actualErr := IsCharacterNameValid(name)
		if expectedErr != actualErr {
			t.Fatalf("Wanted json error: %s, got %s instead", expectedErr, actualErr)
		}
	}
}

func TestCreatureValidator(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	names := []string{
		"Demon",
		"demon",
		"Demons",
		"demons",
		"Druid's Apparitions",
		"Druid's Apparition",
		"apparitionofadruid",
		"Elves",
		"elves",
		"Elf",
		"elf",
		"priesteSsofThEWildsUn",
		"Priestesses Of The Wild Sun",
		"Priestess Of The Wild Sun",
	}

	for _, n := range names {
		_, err := IsCreatureNameValid(n)
		if err != nil {
			t.Fatalf("Creature Name %s should be valid, but it is not, err: %s", n, err)
		}
	}

	invalidNames := []string{
		"One1",
		"Two 2",
		"A",
		"a",
		"T2wo",
		"ab",
		"abcdefghijklmnopqrstuvxzwyabcdefg",
		"",
		"abcabcabcabcabcabcabcabca abc",
		"a abca",
	}

	for _, n := range invalidNames {
		_, err := IsCreatureNameValid(n)
		if err == nil {
			t.Fatalf("Invalid creature name passed as valid, name: %s", n)
		}
	}

	jsonErrorNames := map[Error]string{
		ErrorCreatureNameEmpty:            "",
		ErrorCreatureNameTooSmall:         "ab",
		ErrorCreatureNameTooBig:           "abcdefghijklmnopqrstuvxzwyabcdefg",
		ErrorCreatureWordTooBig:           "abcabcabcabcabcabcabcabca a",
		ErrorCreatureWordTooSmall:         "abca a abc",
		ErrorCreatureNameIsOnlyWhiteSpace: "     ",
		ErrorCreatureNameInvalid:          "12345",
		ErrorCreatureNotFound:             "gaz'haragoth",
	}

	for expectedErr, name := range jsonErrorNames {
		_, actualErr := IsCreatureNameValid(name)
		if expectedErr != actualErr {
			t.Fatalf("Wanted json error: %s, got %s instead", expectedErr, actualErr)
		}
	}

	endpoints := map[string]string{
		"Demon":                       "demon",
		"demon":                       "demon",
		"Demons":                      "demon",
		"demons":                      "demon",
		"Druid's Apparitions":         "apparitionofadruid",
		"Druid's Apparition":          "apparitionofadruid",
		"apparitionofadruid":          "apparitionofadruid",
		"Elves":                       "elf",
		"elves":                       "elf",
		"Elf":                         "elf",
		"elf":                         "elf",
		"priesteSsofThEWildsUn":       "priestessofthewildsun",
		"Priestesses Of The Wild Sun": "priestessofthewildsun",
		"Priestess Of The Wild Sun":   "priestessofthewildsun",
	}

	for name, expectedEndpoint := range endpoints {
		actualEndpoint, err := IsCreatureNameValid(name)
		if err != nil {
			t.Fatal(err)
		}

		if expectedEndpoint != actualEndpoint {
			t.Fatalf("Wanted %s creature endpoint but got %s creature endpoint", expectedEndpoint, actualEndpoint)
		}
	}
}

func TestSpellValidator(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	names := []string{
		"apprenticestrike",
		"Apprentice's Strike",
		"exori min flam",
		"berserk",
		"exori",
		"Conjure Wand of Darkness",
		"conjurewandofdarkness",
		"exevo gran mort",
	}

	for _, n := range names {
		_, err := IsSpellNameOrFormulaValid(n)
		if err != nil {
			t.Fatalf("Spell Name %s should be valid, but it is not, err: %s", n, err)
		}
	}

	invalidNames := []string{
		"One1",
		"Two 2",
		"A",
		"a",
		"T2wo",
		"ab",
		"abcdefghijklmnopqrstuvxzwyabcdefg",
		"",
		"abcabcabcabcabcabcabcabca abc",
		"a abca",
	}

	for _, n := range invalidNames {
		_, err := IsSpellNameOrFormulaValid(n)
		if err == nil {
			t.Fatalf("Invalid spell name passed as valid, name: %s", n)
		}
	}

	jsonErrorNames := map[Error]string{
		ErrorSpellNameEmpty:            "",
		ErrorSpellNameTooSmall:         "ab",
		ErrorSpellNameTooBig:           "abcdefghijklmnopqrstuvxzwyabcdefg",
		ErrorSpellWordTooBig:           "abcabcabcabcabcabcabcab ",
		ErrorSpellWordTooSmall:         "abca a abc",
		ErrorSpellNameIsOnlyWhiteSpace: "     ",
		ErrorSpellNameInvalid:          "12345",
		ErrorSpellNotFound:             "exori'mort",
	}

	for expectedErr, name := range jsonErrorNames {
		_, actualErr := IsSpellNameOrFormulaValid(name)
		if expectedErr != actualErr {
			t.Fatalf("Wanted json error: %s, got %s instead", expectedErr, actualErr)
		}
	}

	endpoints := map[string]string{
		"apprenticestrike":         "apprenticestrike",
		"Apprentice's Strike":      "apprenticestrike",
		"exori min flam":           "apprenticestrike",
		"berserk":                  "berserk",
		"exori":                    "berserk",
		"Conjure Wand of Darkness": "conjurewandofdarkness",
		"conjurewandofdarkness":    "conjurewandofdarkness",
		"exevo gran mort":          "conjurewandofdarkness",
	}

	for name, expectedEndpoint := range endpoints {
		actualEndpoint, err := IsSpellNameOrFormulaValid(name)
		if err != nil {
			t.Fatal(err)
		}

		if expectedEndpoint != actualEndpoint {
			t.Fatalf("Wanted %s spell endpoint but got %s spell endpoint", expectedEndpoint, actualEndpoint)
		}
	}
}

func TestGuildValidator(t *testing.T) {
	names := []string{
		"Pax",
		"Own Way",
		"Hill",
		"Re Evolution",
		"On Top",
	}

	for _, n := range names {
		err := IsGuildNameValid(n)
		if err != nil {
			t.Fatalf("Name %s should be valid, but it is not, err: %s", n, err)
		}
	}

	invalidNames := []string{
		"One1",
		"Two 2",
		"A",
		"a",
		"T2wo",
		"ab",
		"abc'def",
	}

	for _, n := range invalidNames {
		err := IsGuildNameValid(n)
		if err == nil {
			t.Fatalf("Invalid name passed as valid, name: %s", n)
		}
	}

	jsonErrorNames := map[Error]string{
		ErrorGuildNameEmpty:            "",
		ErrorGuildNameTooSmall:         "ob",
		ErrorGuildNameTooBig:           "abcabcabcabcabcabcabcabcabcabc",
		ErrorGuildWordTooBig:           "abcabcabcabcabc hello",
		ErrorGuildNameIsOnlyWhiteSpace: "     ",
		ErrorGuildNameInvalid:          "123",
	}

	for expectedErr, name := range jsonErrorNames {
		actualErr := IsGuildNameValid(name)
		if expectedErr != actualErr {
			t.Fatalf("Wanted json error: %s, got %s instead", expectedErr, actualErr)
		}
	}
}

func TestValidationFuncs(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := GetWorlds()
	if err != nil {
		t.Fatalf("GetWorlds error: %s", err)
	}

	exists, err := WorldExists("antica")
	if err != nil {
		t.Fatalf("WorldExists error with antica: %s", err)
	}

	if !exists {
		t.Fatal("WorldExists is reporting antica does not exist")
	}

	exists, _ = WorldExists("noworld")
	if exists {
		t.Fatal("WorldExists is reporting noworld exists")
	}

	_, err = GetTowns()
	if err != nil {
		t.Fatalf("GetTowns error: %s", err)
	}

	exists, err = TownExists("carlin")
	if err != nil {
		t.Fatalf("TownExists error with carlin: %s", err)
	}

	if !exists {
		t.Fatal("TownExists is reporting carlin does not exist")
	}

	exists, _ = TownExists("notown")
	if exists {
		t.Fatal("TownExists is reporting notown exists")
	}

	_, err = GetHouses()
	if err != nil {
		t.Fatalf("GetHouses error: %s", err)
	}

	house, err := GetHouseRaw(59052)
	if err != nil || house == nil {
		t.Fatalf("GetHouseRaw error with house 59052: %s", err)
	}

	exists, err = HouseExistsRaw(59051)
	if err != nil {
		t.Fatalf("HouseExistsRaw error with house 59051: %s", err)
	}

	if !exists {
		t.Fatal("HouseExistsRaw is reporting house 59051 does not exist")
	}

	exists, _ = HouseExistsRaw(1010)
	if exists {
		t.Fatal("HouseExistsRaw is reporting house 1010 exists")
	}

	house, err = GetHouseInTown(59054, "Ankrahmun")
	if err != nil || house == nil {
		t.Fatalf("GetHouseInTown error with house 59054 in Ankrahmun: %s", err)
	}

	exists, err = HouseExistsInTown(59054, "Ankrahmun")
	if err != nil {
		t.Fatalf("HouseExistsInTown error with house 59054 in Ankrahmun: %s", err)
	}

	if !exists {
		t.Fatal("HouseExistsInTown is reporting house 59054 does not exist in Ankrahmun")
	}

	exists, _ = HouseExistsInTown(59054, "Carlin")
	if exists {
		t.Fatal("HouseExistsInTown is reporting house 59054 exists in Carlin")
	}

	_, err = GetCreatures()
	if err != nil {
		t.Fatalf("GetCreatures error: %s", err)
	}

	_, err = GetSha256Sum()
	if err != nil {
		t.Fatalf("GetSha256Sum error: %s", err)
	}

	_, err = GetSha512Sum()
	if err != nil {
		t.Fatalf("GetSha512Sum error: %s", err)
	}
}

func TestErrors(t *testing.T) {
	type inside struct {
		Code int
	}

	generalError := Error{errors.New("General Error")}
	fakeError := Error{errors.New("general error")}

	errs := map[Error]inside{
		generalError: {
			Code: 0,
		},
		fakeError: {
			Code: 0,
		},
		ErrorAlreadyRunning: {
			Code: 10,
		},
		ErrorValidatorNotInitiated: {
			Code: 11,
		},
		ErrorStringCanNotBeConvertedToInt: {
			Code: 9001,
		},
		ErrorCharacterNameEmpty: {
			Code: 10001,
		},
		ErrorCharacterNameTooSmall: {
			Code: 10002,
		},
		ErrorCharacterNameInvalid: {
			Code: 10003,
		},
		ErrorCharacterNameIsOnlyWhiteSpace: {
			Code: 10004,
		},
		ErrorCharacterNameTooBig: {
			Code: 10005,
		},
		ErrorCharacterWordTooBig: {
			Code: 10006,
		},
		ErrorCharacterWordTooSmall: {
			Code: 10007,
		},
		ErrorInvalidNewsID: {
			Code: 11001,
		},
		ErrorWorldDoesNotExist: {
			Code: 11002,
		},
		ErrorVocationDoesNotExist: {
			Code: 11003,
		},
		ErrorHighscoreCategoryDoesNotExist: {
			Code: 11004,
		},
		ErrorHouseDoesNotExist: {
			Code: 11005,
		},
		ErrorTownDoesNotExist: {
			Code: 11006,
		},
		ErrorCreatureNameEmpty: {
			Code: 12001,
		},
		ErrorCreatureNameTooSmall: {
			Code: 12002,
		},
		ErrorCreatureNameInvalid: {
			Code: 12003,
		},
		ErrorCreatureNameIsOnlyWhiteSpace: {
			Code: 12004,
		},
		ErrorCreatureNameTooBig: {
			Code: 12005,
		},
		ErrorCreatureWordTooBig: {
			Code: 12006,
		},
		ErrorCreatureWordTooSmall: {
			Code: 12007,
		},
		ErrorSpellNameEmpty: {
			Code: 13001,
		},
		ErrorSpellNameTooSmall: {
			Code: 13002,
		},
		ErrorSpellNameInvalid: {
			Code: 13003,
		},
		ErrorSpellNameIsOnlyWhiteSpace: {
			Code: 13004,
		},
		ErrorSpellNameTooBig: {
			Code: 13005,
		},
		ErrorSpellWordTooBig: {
			Code: 13006,
		},
		ErrorSpellWordTooSmall: {
			Code: 13007,
		},
		ErrorGuildNameEmpty: {
			Code: 14001,
		},
		ErrorGuildNameTooSmall: {
			Code: 14002,
		},
		ErrorGuildNameInvalid: {
			Code: 14003,
		},
		ErrorGuildNameIsOnlyWhiteSpace: {
			Code: 14004,
		},
		ErrorGuildNameTooBig: {
			Code: 14005,
		},
		ErrorGuildWordTooBig: {
			Code: 14006,
		},
		ErrorGuildWordTooSmall: {
			Code: 14007,
		},
		ErrorCharacterNotFound: {
			Code: 20001,
		},
		ErrorCreatureNotFound: {
			Code: 20002,
		},
		ErrorSpellNotFound: {
			Code: 20003,
		},
		ErrorGuildNotFound: {
			Code: 20004,
		},
	}

	for err, values := range errs {
		if err.Code() != values.Code {
			t.Fatalf("Err %s Code should return %d, but it returned %d", err, values.Code, err.Code())
		}
	}
}

func TestUtils(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	wordStrings := []string{
		"abc",
		"世界",
		"wêreld",
		"κόσμος",
		"ਸੰਸਾਰ",
		"ዓለም",
		"العالمية",
	}

	digitStrings := []string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"a1",
		"1a",
		"abc123abc",
	}

	numberStrings := []string{
		"Ⅰ",
		"Ⅷ",
		"½",
	}

	for _, w := range wordStrings {
		if DoesStringContainDigits(w) {
			t.Fatalf("DoesStringContainDigits is saying %s contains a digit, but it doesn't", w)
		}

		if DoesStringContainsNumbers(w) {
			t.Fatalf("DoesStringContainsNumbers is saying %s contains a number, but it doesn't", w)
		}
	}

	for _, d := range digitStrings {
		if !DoesStringContainDigits(d) {
			t.Fatalf("DoesStringContainDigits is saying %s doesn't contain a digit, but it does", d)
		}

		if !DoesStringContainsNumbers(d) {
			t.Fatalf("DoesStringContainsNumbers is saying %s doesn't contain a number, but it does", d)
		}
	}

	for _, n := range numberStrings {
		if DoesStringContainDigits(n) {
			t.Fatalf("DoesStringContainDigits is saying %s contains a digit, but it doesn't", n)
		}

		if !DoesStringContainsNumbers(n) {
			t.Fatalf("DoesStringContainsNumbers is saying %s doesn't contain a number, but it does", n)
		}
	}

	name, err := GetSmallestCreatureName()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("smallestCreatureName not set")
	}

	name, err = GetBiggestCreatureName()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("biggestCreatureName not set")
	}

	name, err = GetSmallestCreatureWord()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("smallestCreatureWord not set")
	}

	name, err = GetBiggestCreatureWord()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("biggestCreatureWord not set")
	}

	count, err := GetSmallestCreatureNameRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("smallestCreatureNameRuneCount not set")
	}

	count, err = GetBiggestCreatureNameRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("biggestCreatureNameRuneCount not set")
	}

	count, err = GetSmallestCreatureWordRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("smallestCreatureWordRuneCount not set")
	}

	count, err = GetBiggestCreatureWordRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("biggestCreatureWordRuneCount not set")
	}

	name, err = GetSmallestSpellNameOrFormula()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("smallestSpellNameOrFormula not set")
	}

	name, err = GetBiggestSpellNameOrFormula()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("biggestSpellNameOrFormula not set")
	}

	name, err = GetSmallestSpellWord()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("smallestSpellWord not set")
	}

	name, err = GetBiggestSpellWord()
	if err != nil {
		t.Fatal(err)
	}

	if name == "" {
		t.Fatal("biggestSpellWord not set")
	}

	count, err = GetSmallestSpellNameOrFormulaRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("smallestSpellNameOrFormulaRuneCount not set")
	}

	count, err = GetBiggestSpellNameOrFormulaRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("biggestSpellNameOrFormulaRuneCount not set")
	}

	count, err = GetSmallestSpellWordRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("smallestSpellWorldRuneCount not set")
	}

	count, err = GetBiggestSpellWordRuneCount()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal("biggestSpellWorldRuneCount not set")
	}
}

func TestFake(t *testing.T) {
	if !initiated {
		err := Initiate("TibiaData-API-Testing")
		if err != nil {
			t.Fatal(err)
		}
	}

	assert := assert.New(t)

	assert.Equal(29, MaxRunesAllowedInACharacterName)
	assert.Equal(2, MinRunesAllowedInACharacterName)
	assert.Equal(14, MaxRunesAllowedInACharacterNameWord)
	assert.Equal(2, MinRunesAllowedInACharacterNameWord)

	assert.Equal(29, MaxRunesAllowedInAGuildName)
	assert.Equal(3, MinRunesAllowedInAGuildName)
	assert.Equal(14, MaxRunesAllowedInAGuildNameWord)
	assert.Equal(2, MinRunesAllowedInAGuildNameWord)

	setVars()
	setCreaturesVars()
	setSpellsVars()
}

func TestNewsIDValidator(t *testing.T) {
	err := IsNewsIDValid(1)
	if err != nil {
		t.Fatalf("News ID 1 is being reported as invalid but should be valid, err: %s", err)
	}

	err = IsNewsIDValid(0)
	if err == nil {
		t.Fatal("News ID 0 is being reported as valid but should be invalid")
	}

	err = IsNewsIDValid(-1)
	if err == nil {
		t.Fatal("News ID 0 is being reported as valid but should be invalid")
	}
}

func TestVocationValidator(t *testing.T) {
	err := IsVocationValid("tibia")
	if err == nil {
		t.Fatal("Vocation tibia is being reported as valid but should be invalid")
	}

	err = IsVocationValid("sorcerer")
	if err != nil {
		t.Fatalf("Vocation sorcerer is being reported as invalid but should be valid, err: %s", err)
	}

	err = IsVocationValid("sorcerers")
	if err != nil {
		t.Fatalf("Vocation sorcerers is being reported as invalid but should be valid, err: %s", err)
	}
}
