package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of SpellInformation
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

// Child of RuneInformation
type RuneInformation struct {
	Vocation     []string `json:"vocation"`
	GroupAttack  bool     `json:"group_attack"`
	GroupHealing bool     `json:"group_healing"`
	GroupSupport bool     `json:"group_support"`
	DamageType   string   `json:"damage_type"`
	Level        int      `json:"level"`
	MagicLevel   int      `json:"magic_level"`
}

// Child of Spells
type SpellData struct {
	Name                string           `json:"name"`
	Spell               string           `json:"spell_id"`
	ImageURL            string           `json:"image_url"`
	Description         string           `json:"description"`
	HasSpellInformation bool             `json:"has_spell_information"`
	SpellInformation    SpellInformation `json:"spell_information"`
	HasRuneInformation  bool             `json:"has_rune_information"`
	RuneInformation     RuneInformation  `json:"rune_information"`
}

// Child of JSONData
type SpellsContainer struct {
	Spell SpellData `json:"spell"`
}

//
// The base includes two levels: Spell and Information
type SpellInformationResponse struct {
	Spells      SpellsContainer `json:"spells"`
	Information Information     `json:"information"`
}

var (
	SpellDataRowRegex      = regexp.MustCompile(`<td.*>(.*):<\/td><td.*>(.*)<\/td>`)
	SpellNameAndImageRegex = regexp.MustCompile(`<td><img src="(.*)" width=.*<h2>(.*)<\/h2>.*`)
	SpellCooldownRegex     = regexp.MustCompile(`([0-9]+)s \(.*:.([0-9]+)s\)`)
	SpellDescriptionRegex  = regexp.MustCompile(`(.*)\.(Spell|Rune) InformationName:.*`)
)

// TibiaSpellsSpellV3 func
func TibiaSpellsSpellV3Impl(spell string, BoxContentHTML string) (*SpellInformationResponse, error) {
	//TODO: There is currently a bug with description, it always comes back empty

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaSpellsSpellV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var (
		// creating empty vars for later use
		SpellsInfoVocation, SpellsInfoCity, RuneInfoVocation []string
		// var SpellsInfoName, RuneInfoName string
		SpellInformationSection, SpellName, SpellImageURL, SpellDescription, SpellsInfoFormula, SpellsInfoDamageType, RuneInfoDamageType                                                                                                                  string
		SpellsInfoCooldownAlone, SpellsInfoCooldownGroup, SpellsInfoSoulPoints, SpellsInfoAmount, SpellsInfoLevel, SpellsInfoMana, SpellsInfoPrice, RuneInfoLevel, RuneInfoMagicLevel                                                                     int
		SpellsInfoGroupAttack, SpellsInfoGroupHealing, SpellsInfoGroupSupport, SpellsInfoTypeInstant, SpellsInfoTypeRune, RuneInfoGroupAttack, RuneInfoGroupHealing, RuneInfoGroupSupport, SpellsInfoPremium, SpellsHasSpellSection, SpellsHasRuneSection bool

		insideError error
	)

	ReaderHTML.Find(".BoxContent").EachWithBreak(func(index int, s *goquery.Selection) bool {
		NameAndImageSection, err := s.Find("table tr").First().Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaSpellsSpellV3Impl failed at NameAndImageSection, err := s.Find, err: %s", err)
			return false
		}

		// Get the name and image
		subma2 := SpellNameAndImageRegex.FindAllStringSubmatch(NameAndImageSection, -1)
		if len(subma2) > 0 {
			SpellName = TibiaDataSanitize0026String(subma2[0][2])
			SpellImageURL = subma2[0][1]
		}

		s.Find(".TableContainer").Each(func(index int, s *goquery.Selection) {
			SectionName := s.Find(".CaptionInnerContainer div.Text").Text()

			// Determine if this is the spell or rune section
			if SectionName == "Spell Information" {
				SpellInformationSection = "spell"
				SpellsHasSpellSection = true
			} else if SectionName == "Rune Information" {
				SpellInformationSection = "rune"
				SpellsHasRuneSection = true
			}

			// Running query over each div
			s.Find("table.Table2 tbody tr").EachWithBreak(func(index int, s *goquery.Selection) bool {
				// Storing HTML into SpellDivHTML
				SpellDivHTML, err := s.Html()
				if err != nil {
					insideError = fmt.Errorf("[error] TibiaSpellsSpellV3Impl failed at SpellDivHTML, err := s.Html(), err: %s", err)
					return false
				}

				subma1 := SpellDataRowRegex.FindAllStringSubmatch(SpellDivHTML, -1)

				// check if regex return length is over 0 and the match of name is over 1
				if len(subma1) > 0 {
					// Creating easy to use vars (and unescape hmtl right string)
					LeftColumn := subma1[0][1]
					RightColumn := TibiaDataSanitizeEscapedString(subma1[0][2])

					// Formula
					if LeftColumn == "Formula" {
						SpellsInfoFormula = TibiaDataSanitizeDoubleQuoteString(RightColumn)
					}

					// Vocation
					if LeftColumn == "Vocation" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoVocation = strings.Split(RightColumn, ", ")
						case "rune":
							RuneInfoVocation = strings.Split(RightColumn, ", ")
						}
					}

					// Group information
					if LeftColumn == "Group" {
						switch SpellInformationSection {
						case "spell":
							switch RightColumn {
							case "Attack":
								SpellsInfoGroupAttack = true
							case "Healing":
								SpellsInfoGroupHealing = true
							case "Support":
								SpellsInfoGroupSupport = true
							}
						case "rune":
							switch RightColumn {
							case "Attack":
								RuneInfoGroupAttack = true
							case "Healing":
								RuneInfoGroupHealing = true
							case "Support":
								RuneInfoGroupSupport = true
							}
						}
					}

					// Spell type
					if LeftColumn == "Type" {
						switch RightColumn {
						case "Instant":
							SpellsInfoTypeInstant = true
						case "Rune":
							SpellsInfoTypeRune = true
						}
					}

					// Damage
					if LeftColumn == "Damage Type" || LeftColumn == "Magic Type" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoDamageType = strings.ToLower(RightColumn)
						case "rune":
							RuneInfoDamageType = strings.ToLower(RightColumn)
						}
					}

					// Cooldown
					if LeftColumn == "Cooldown" {
						subma3 := SpellCooldownRegex.FindAllStringSubmatch(SpellDivHTML, -1)
						if len(subma3) > 0 {
							SpellsInfoCooldownAlone = TibiaDataStringToIntegerV3(subma3[0][1])
							SpellsInfoCooldownGroup = TibiaDataStringToIntegerV3(subma3[0][2])
						}
					}

					// Soul Points
					if LeftColumn == "Soul Points" {
						SpellsInfoSoulPoints = TibiaDataStringToIntegerV3(RightColumn)
					}

					// Amount
					if LeftColumn == "Amount" {
						SpellsInfoAmount = TibiaDataStringToIntegerV3(RightColumn)
					}

					// Experience Level
					if LeftColumn == "Exp Lvl" {
						switch SpellInformationSection {
						case "spell":
							SpellsInfoLevel = TibiaDataStringToIntegerV3(RightColumn)
						case "rune":
							RuneInfoLevel = TibiaDataStringToIntegerV3(RightColumn)
						}
					}

					// Mana
					if LeftColumn == "Mana" {
						SpellsInfoMana = TibiaDataStringToIntegerV3(RightColumn)
					}

					// Price
					if LeftColumn == "Price" {
						if RightColumn == "free" {
							SpellsInfoPrice = 0
						} else {
							SpellsInfoPrice = TibiaDataStringToIntegerV3(RightColumn)
						}
					}

					// City
					if LeftColumn == "City" {
						SpellsInfoCity = strings.Split(RightColumn, ", ")
					}

					// Premium
					if LeftColumn == "Premium" {
						if RightColumn == "yes" {
							SpellsInfoPremium = true
						} else {
							SpellsInfoPremium = false
						}
					}

					// Magic level
					if LeftColumn == "Mag Lvl" {
						RuneInfoMagicLevel = TibiaDataStringToIntegerV3(RightColumn)
					}
				}

				return true
			})
		})

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	// Getting the description
	InnerTableContainerTMPB := ReaderHTML.Find(".BoxContent").Text()
	subma4 := SpellDescriptionRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)
	if len(subma4) > 0 {
		SpellDescription = subma4[0][1] + "."
	}

	//
	// Build the data-blob
	return &SpellInformationResponse{
		SpellsContainer{
			SpellData{
				Name:                SpellName,
				Spell:               strings.ToLower(SpellName),
				ImageURL:            SpellImageURL,
				Description:         SpellDescription,
				HasSpellInformation: SpellsHasSpellSection,
				SpellInformation: SpellInformation{
					Formula:       SpellsInfoFormula,
					Vocation:      SpellsInfoVocation,
					GroupAttack:   SpellsInfoGroupAttack,
					GroupHealing:  SpellsInfoGroupHealing,
					GroupSupport:  SpellsInfoGroupSupport,
					TypeInstant:   SpellsInfoTypeInstant,
					TypeRune:      SpellsInfoTypeRune,
					DamageType:    SpellsInfoDamageType,
					CooldownAlone: SpellsInfoCooldownAlone,
					CooldownGroup: SpellsInfoCooldownGroup,
					SoulPoints:    SpellsInfoSoulPoints,
					Amount:        SpellsInfoAmount,
					Level:         SpellsInfoLevel,
					Mana:          SpellsInfoMana,
					Price:         SpellsInfoPrice,
					City:          SpellsInfoCity,
					Premium:       SpellsInfoPremium,
				},
				HasRuneInformation: SpellsHasRuneSection,
				RuneInformation: RuneInformation{
					Vocation:     RuneInfoVocation,
					GroupAttack:  RuneInfoGroupAttack,
					GroupHealing: RuneInfoGroupHealing,
					GroupSupport: RuneInfoGroupSupport,
					DamageType:   RuneInfoDamageType,
					Level:        RuneInfoLevel,
					MagicLevel:   RuneInfoMagicLevel,
				},
			},
		},
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
