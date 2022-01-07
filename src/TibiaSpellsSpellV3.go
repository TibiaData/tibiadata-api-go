package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"tibiadata-api-go/src/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaSpellsSpellV3 func
func TibiaSpellsSpellV3(c *gin.Context) {
	// getting params from URL
	spell := c.Param("spell")

	// The base includes two levels: Spell and Information
	type JSONData struct {
		Spell       structs.Spell       `json:"spells"`
		Information structs.Information `json:"information"`
	}

	// Setting spells string to lower chars
	spell = strings.ToLower(spell)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/library/?subtopic=spells&spell=" + TibiadataQueryEscapeStringV3(spell)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaSpellsSpellV3", gin.H{"error": err.Error()})
		return
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// creating empty vars for later use
	var (
		SpellsInfoVocation, SpellsInfoCity, RuneInfoVocation []string
		// SpellsInfoName, RuneInfoName string
		SpellInformationSection, SpellName, SpellImageURL, SpellDescription, SpellsInfoFormula, SpellsInfoDamageType, RuneInfoDamageType                                                                                                                  string
		SpellsInfoCooldownAlone, SpellsInfoCooldownGroup, SpellsInfoSoulPoints, SpellsInfoAmount, SpellsInfoLevel, SpellsInfoMana, SpellsInfoPrice, RuneInfoLevel, RuneInfoMagicLevel                                                                     int
		SpellsInfoGroupAttack, SpellsInfoGroupHealing, SpellsInfoGroupSupport, SpellsInfoTypeInstant, SpellsInfoTypeRune, RuneInfoGroupAttack, RuneInfoGroupHealing, RuneInfoGroupSupport, SpellsInfoPremium, SpellsHasSpellSection, SpellsHasRuneSection bool
	)

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

			switch WorldsInformationLeftColumn {
			case "Formula":
				SpellsInfoFormula = TibiaDataSanitizeDoubleQuoteString(WorldsInformationRightColumn)
			case "Vocation":
				switch SpellInformationSection {
				case "spell":
					SpellsInfoVocation = strings.Split(WorldsInformationRightColumn, ", ")
				case "rune":
					RuneInfoVocation = strings.Split(WorldsInformationRightColumn, ", ")
				}
			case "Group":
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
			case "Type":
				switch WorldsInformationRightColumn {
				case "Instant":
					SpellsInfoTypeInstant = true
				case "Rune":
					SpellsInfoTypeRune = true
				}
			case "Damage Type":
				switch SpellInformationSection {
				case "spell":
					SpellsInfoDamageType = strings.ToLower(WorldsInformationRightColumn)
				case "rune":
					RuneInfoDamageType = strings.ToLower(WorldsInformationRightColumn)
				}
			case "Cooldown":
				regex3 := regexp.MustCompile(`([0-9]+)s \(.*:.([0-9]+)s\)`)
				subma3 := regex3.FindAllStringSubmatch(SpellDivHTML, -1)
				if len(subma3) > 0 {
					SpellsInfoCooldownAlone = TibiadataStringToIntegerV3(subma3[0][1])
					SpellsInfoCooldownGroup = TibiadataStringToIntegerV3(subma3[0][2])
				}
			case "Soul Points":
				SpellsInfoSoulPoints = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			case "Amount":
				SpellsInfoAmount = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			case "Exp Lvl":
				switch SpellInformationSection {
				case "spell":
					SpellsInfoLevel = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				case "rune":
					RuneInfoLevel = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				}
			case "Mana":
				SpellsInfoMana = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
			case "Price":
				if WorldsInformationRightColumn == "free" {
					SpellsInfoPrice = 0
				} else {
					SpellsInfoPrice = TibiadataStringToIntegerV3(WorldsInformationRightColumn)
				}
			case "City":
				SpellsInfoCity = strings.Split(WorldsInformationRightColumn, ", ")
			case "Premium":
				if WorldsInformationRightColumn == "yes" {
					SpellsInfoPremium = true
				} else {
					SpellsInfoPremium = false
				}
			case "Mag Lvl":
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
		structs.Spell{
			Name:                SpellName,
			Spell:               spell,
			ImageURL:            SpellImageURL,
			Description:         SpellDescription,
			HasSpellInformation: SpellsHasSpellSection,
			SpellInformation: structs.SpellInformation{
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
			RuneInformation: structs.RuneInformation{
				Vocation:     RuneInfoVocation,
				GroupAttack:  RuneInfoGroupAttack,
				GroupHealing: RuneInfoGroupHealing,
				GroupSupport: RuneInfoGroupSupport,
				DamageType:   RuneInfoDamageType,
				Level:        RuneInfoLevel,
				MagicLevel:   RuneInfoMagicLevel,
			},
		},
		structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaSpellsSpellV3", jsonData)
}
