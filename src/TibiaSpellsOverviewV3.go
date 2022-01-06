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

	// Sanitize of vocation input
	vocationName, _ := TibiaDataVocationValidator(vocation)
	if vocationName == "all" || vocationName == "none" {
		vocationName = ""
	} else {
		// removes the last letter (s) from the string (required for spells page)
		vocationName = strings.TrimSuffix(vocationName, "s")
		// setting string to first upper case
		vocationName = strings.Title(vocationName)
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest = map[string]map[string]string{
		"request": {
			"method": "GET",
			"url":    "https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiadataQueryEscapeStringV3(vocationName),
		}}
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

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

			switch subma1[0][4] {
			case "Attack":
				GroupAttack = true
			case "Healing":
				GroupHealing = true
			case "Support":
				GroupSupport = true
			}

			// Type
			TypeInstant = false
			TypeRune = false

			switch subma1[0][5] {
			case "Instant":
				TypeInstant = true
			case "Rune":
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
				Formula:      TibiaDataSanitizeDoubleQuoteString(TibiaDataSanitizeEscapedString(subma1[0][3])),
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
	if vocationName == "" {
		vocationName = "all"
	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Spells{
			SpellsVocationFilter: vocationName,
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
