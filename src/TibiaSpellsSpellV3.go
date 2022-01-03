package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaSpellsSpellV3 func
func TibiaSpellsSpellV3(c *gin.Context) {

	// getting params from URL
	spell := c.Param("spell")

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

	// Child of JSONData
	type Spells struct {
		Spell Spell `json:"spell"`
	}

	//
	// The base includes two levels: Spell and Information
	type JSONData struct {
		Spells      Spells      `json:"spells"`
		Information Information `json:"information"`
	}

	// Setting spells string to lower chars
	spell = strings.ToLower(spell)

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/library/?subtopic=spells&spell=" + TibiadataQueryEscapeStringV3(spell))

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// creating empty vars for later use
	var SpellsInfoVocation, SpellsInfoCity, RuneInfoVocation []string
	// var SpellsInfoName, RuneInfoName string
	var SpellInformationSection, SpellName, SpellImageURL, SpellDescription, SpellsInfoFormula, SpellsInfoDamageType, RuneInfoDamageType string
	var SpellsInfoCooldownAlone, SpellsInfoCooldownGroup, SpellsInfoSoulPoints, SpellsInfoAmount, SpellsInfoLevel, SpellsInfoMana, SpellsInfoPrice, RuneInfoLevel, RuneInfoMagicLevel int
	var SpellsInfoGroupAttack, SpellsInfoGroupHealing, SpellsInfoGroupSupport, SpellsInfoTypeInstant, SpellsInfoTypeRune, RuneInfoGroupAttack, RuneInfoGroupHealing, RuneInfoGroupSupport, SpellsInfoPremium, SpellsHasSpellSection, SpellsHasRuneSection bool

	// Running query over each div
	ReaderHTML.Find(".BoxContent table tbody tr").Each(func(index int, s *goquery.Selection) {

		// Storing HTML into SpellDivHTML
		SpellDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for name, race and img src param for creature
		regex1 := regexp.MustCompile(`<td.*>(.*):<\/td><td.*>(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(SpellDivHTML, -1)

		// Get the name and image
		regex2 := regexp.MustCompile(`<td><img src="(.*)" width=.*<h2>(.*)<\/h2>.*`)
		subma2 := regex2.FindAllStringSubmatch(SpellDivHTML, -1)
		if len(subma2) > 0 {
			SpellName = subma2[0][2]
			SpellImageURL = subma2[0][1]
		}

		// Determine if this is the spell or rune section
		if strings.Contains(SpellDivHTML, "<b>Spell Information</b>") {
			SpellInformationSection = "spell"
			SpellsHasSpellSection = true
		} else if strings.Contains(SpellDivHTML, "<b>Rune Information</b>") {
			SpellInformationSection = "rune"
			SpellsHasRuneSection = true
		}

		// check if regex return length is over 0 and the match of name is over 1
		if len(subma1) > 0 {

			// Creating easy to use vars (and unescape hmtl right string)
			WorldsInformationLeftColumn := subma1[0][1]
			WorldsInformationRightColumn := TibiaDataSanitizeEscapedString(subma1[0][2])

			// Formula
			if WorldsInformationLeftColumn == "Formula" {
				SpellsInfoFormula = TibiaDataSanitizeDoubleQuoteString(WorldsInformationRightColumn)
			}

			// Vocation
			if WorldsInformationLeftColumn == "Vocation" {
				switch SpellInformationSection {
				case "spell":
					SpellsInfoVocation = strings.Split(WorldsInformationRightColumn, ", ")
				case "rune":
					RuneInfoVocation = strings.Split(WorldsInformationRightColumn, ", ")
				}
			}

			// Group information
			if WorldsInformationLeftColumn == "Group" {
				switch SpellInformationSection {
				case "spell":
					switch WorldsInformationRightColumn {
					case "Attack":
						SpellsInfoGroupAttack = true
					case "Healing":
						SpellsInfoGroupHealing = true
					case "Support":
						SpellsInfoGroupSupport = true
					}
				case "rune":
					switch WorldsInformationRightColumn {
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
			if WorldsInformationLeftColumn == "Type" {
				switch WorldsInformationRightColumn {
				case "Instant":
					SpellsInfoTypeInstant = true
				case "Rune":
					SpellsInfoTypeRune = true
				}
			}

			// Damage
			if WorldsInformationLeftColumn == "Damage Type" {
				switch SpellInformationSection {
				case "spell":
					SpellsInfoDamageType = strings.ToLower(WorldsInformationRightColumn)
				case "rune":
					RuneInfoDamageType = strings.ToLower(WorldsInformationRightColumn)
				}
			}

			// Cooldown
			if WorldsInformationLeftColumn == "Cooldown" {
				regex3 := regexp.MustCompile(`([0-9]+)s \(.*:.([0-9]+)s\)`)
				subma3 := regex3.FindAllStringSubmatch(SpellDivHTML, -1)
				if len(subma3) > 0 {
					SpellsInfoCooldownAlone = TibiadataStringToIntegerV3(subma3[0][1])
					SpellsInfoCooldownGroup = TibiadataStringToIntegerV3(subma3[0][2])
				}

			}

			// Soul Points
			if WorldsInformationLeftColumn == "Soul Points" {
				SpellsInfoSoulPoints = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			}

			// Amount
			if WorldsInformationLeftColumn == "Amount" {
				SpellsInfoAmount = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			}

			// Experience Level
			if WorldsInformationLeftColumn == "Exp Lvl" {
				switch SpellInformationSection {
				case "spell":
					SpellsInfoLevel = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				case "rune":
					RuneInfoLevel = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				}
			}

			// Mana
			if WorldsInformationLeftColumn == "Mana" {
				SpellsInfoMana = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			}

			// Price
			if WorldsInformationLeftColumn == "Price" {
				if WorldsInformationRightColumn == "free" {
					SpellsInfoPrice = 0
				} else {
					SpellsInfoPrice = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				}
			}

			// City
			if WorldsInformationLeftColumn == "City" {
				SpellsInfoCity = strings.Split(WorldsInformationRightColumn, ", ")
			}

			// Premium
			if WorldsInformationLeftColumn == "Premium" {
				if WorldsInformationRightColumn == "yes" {
					SpellsInfoPremium = true
				} else {
					SpellsInfoPremium = false
				}
			}

			// Magic level
			if WorldsInformationLeftColumn == "Mag Lvl" {
				RuneInfoMagicLevel = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			}

		}
	})

	// Getting the description
	InnerTableContainerTMPB := ReaderHTML.Find(".BoxContent").Text()
	regex4 := regexp.MustCompile(SpellName + `(.*)\.(Spell|Rune) InformationName:.*`)
	subma4 := regex4.FindAllStringSubmatch(InnerTableContainerTMPB, -1)
	if len(subma4) > 0 {
		SpellDescription = subma4[0][1] + "."
	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Spells{
			Spell{
				Name:                SpellName,
				Spell:               spell,
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
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaSpellsSpellV3", jsonData)
}
