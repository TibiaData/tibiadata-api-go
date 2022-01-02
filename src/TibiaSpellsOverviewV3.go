package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaSpellsOverviewV3 func
func TibiaSpellsOverviewV3(c *gin.Context) {

	// getting params from URL
	vocation := c.Param("vocation")
	if vocation == "" {
		vocation = TibiadataDefaultVoc
	}

	// Child of Spells
	type Spell struct {
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

	// Child of JSONData
	type Spells struct {
		SpellsVocationFilter string  `json:"spells_filter"`
		Spells               []Spell `json:"spell_list"`
	}

	//
	// The base includes two levels: Spells and Information
	type JSONData struct {
		Spells      Spells      `json:"spells"`
		Information Information `json:"information"`
	}

	// Sanatize of vocation value
	vocation = strings.ToLower(vocation)
	if len(vocation) > 0 {
		if strings.EqualFold(vocation, "knight") || strings.EqualFold(vocation, "knights") {
			vocation = strings.Title("knight")
		} else if strings.EqualFold(vocation, "paladin") || strings.EqualFold(vocation, "paladins") {
			vocation = strings.Title("paladin")
		} else if strings.EqualFold(vocation, "sorcerer") || strings.EqualFold(vocation, "sorcerers") {
			vocation = strings.Title("sorcerer")
		} else if strings.EqualFold(vocation, "druid") || strings.EqualFold(vocation, "druids") {
			vocation = strings.Title("druid")
		} else {
			vocation = ""
		}
	} else {
		vocation = ""
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiadataQueryEscapeStringV3(vocation))

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty SpellsData var
	var SpellsData []Spell
	var GroupAttack, GroupHealing, GroupSupport, TypeInstant, TypeRune, PremiumOnly bool

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer table tr").Each(func(index int, s *goquery.Selection) {

		// Storing HTML into SpellDivHTML
		SpellDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for name, race and img src param for creature
		regex1 := regexp.MustCompile(`<td>.*spell=(.*)&amp;voc.*">(.*)<\/a> \((.*)\)<\/td><td>(.*)<\/td><td>(.*)<\/td><td>([0-9]+)<\/td><td>([0-9]+)<\/td><td>([0-9]+)<\/td><td>(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(SpellDivHTML, 1)

		// check if regex return length is over 0 and the match of name is over 1
		if len(subma1) > 0 {

			// SpellGroup
			GroupAttack = false
			GroupHealing = false
			GroupSupport = false

			if subma1[0][4] == "Attack" {
				GroupAttack = true
			} else if subma1[0][4] == "Healing" {
				GroupHealing = true
			} else if subma1[0][4] == "Support" {
				GroupSupport = true
			}

			// Type
			TypeInstant = false
			TypeRune = false
			if subma1[0][5] == "Instant" {
				TypeInstant = true
			} else if subma1[0][5] == "Rune" {
				TypeRune = true
			}

			// PremiumOnly
			if subma1[0][9] == "yes" {
				PremiumOnly = true
			} else {
				PremiumOnly = false
			}

			// Creating data block to return
			SpellsData = append(SpellsData, Spell{
				Name:         subma1[0][2],
				Spell:        subma1[0][1],
				Formula:      TibiaDataSanitizeDoubleQuoteString(TibiadataUnescapeStringV3(subma1[0][3])),
				Level:        TibiadataStringToIntegerV3(subma1[0][6]),
				Mana:         TibiadataStringToIntegerV3(subma1[0][7]),
				Price:        TibiadataStringToIntegerV3(subma1[0][8]),
				GroupAttack:  GroupAttack,
				GroupHealing: GroupHealing,
				GroupSupport: GroupSupport,
				TypeInstant:  TypeInstant,
				TypeRune:     TypeRune,
				PremiumOnly:  PremiumOnly,
			})
		}

	})

	// adding readable SpellsVocationFilter field
	if vocation == "" {
		vocation = "none"
	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Spells{
			SpellsVocationFilter: vocation,
			Spells:               SpellsData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaSpellsOverviewV3", jsonData)
}
