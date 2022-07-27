package main

import (
	"errors"
	"strings"
)

type HighscoreCategory int

const (
	achievements HighscoreCategory = iota + 1
	axefighting
	charmpoints
	clubfighting
	distancefighting
	experience
	fishing
	fistfighting
	goshnarstaint
	loyaltypoints
	magiclevel
	shielding
	swordfighting
	dromescore
	bosspoints
)

func (hc HighscoreCategory) String() (string, error) {
	seasons := [...]string{"achievements", "axefighting", "charmpoints", "clubfighting", "distancefighting", "experience", "fishing", "fistfighting", "goshnarstaint", "loyaltypoints", "magiclevel", "shielding", "swordfighting", "dromescore", "bosspoints"}
	if hc < achievements || hc > bosspoints {
		return "", errors.New("invalid HighscoreCategory value")
	}
	return seasons[hc-1], nil
}

func HighscoreCategoryFromString(input string) HighscoreCategory {
	// Sanatize of category value
	input = strings.ToLower(input)
	switch input {
	case "achievements", "achievement":
		return achievements
	case "axe", "axefighting":
		return axefighting
	case "charm", "charms", "charmpoints":
		return charmpoints
	case "club", "clubfighting":
		return clubfighting
	case "distance", "distancefighting":
		return distancefighting
	case "fishing":
		return fishing
	case "fist", "fistfighting":
		return fistfighting
	case "goshnar", "goshnars", "goshnarstaint":
		return goshnarstaint
	case "loyalty", "loyaltypoints":
		return loyaltypoints
	case "magic", "mlvl", "magiclevel":
		return magiclevel
	case "shielding", "shield":
		return shielding
	case "sword", "swordfighting":
		return swordfighting
	case "drome", "dromescore":
		return dromescore
	case "boss", "bosses", "bosspoints":
		return bosspoints
	default:
		return experience
	}
}
