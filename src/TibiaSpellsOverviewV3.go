package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

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
type SpellsOverviewResponse struct {
	Spells      Spells      `json:"spells"`
	Information Information `json:"information"`
}

func TibiaSpellsOverviewV3(c *gin.Context) {
	// getting params from URL
	vocation := c.Param("vocation")
	if vocation == "" {
		vocation = TibiadataDefaultVoc
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
	TibiadataRequest.URL = "https://www.tibia.com/library/?subtopic=spells&vocation=" + TibiadataQueryEscapeStringV3(vocationName)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaSpellsOverviewV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaSpellsOverviewV3Impl(vocationName, BoxContentHTML)

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaSpellsOverviewV3", jsonData)
}

// TibiaSpellsOverviewV3 func
func TibiaSpellsOverviewV3Impl(vocationName string, BoxContentHTML string) SpellsOverviewResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	var SpellsData []Spell

	// Running query over each div
	ReaderHTML.Find("table.TableContent ~ table tr").Each(func(index int, s *goquery.Selection) {
		//Skip header row
		if index == 0 {
			return
		}

		spellBuilder := Spell{}

		s.Find("td").Each(func(index int, s2 *goquery.Selection) {
			selectionText := s2.Text()

			switch index {
			case 0:
				spellBuilder.Name = selectionText
				spellBuilder.Spell = selectionText[0:strings.Index(selectionText, " (")]
				spellBuilder.Formula = selectionText[strings.Index(selectionText, " (")+2 : strings.Index(selectionText, ")")]
			case 1:
				switch selectionText {
				case "Attack":
					spellBuilder.GroupAttack = true
				case "Healing":
					spellBuilder.GroupHealing = true
				case "Support":
					spellBuilder.GroupSupport = true
				}
			case 2:
				switch selectionText {
				case "Instant":
					spellBuilder.TypeInstant = true
				case "Rune":
					spellBuilder.TypeRune = true
				}
			case 3:
				spellBuilder.Level = TibiadataStringToIntegerV3(selectionText)
			case 4:
				mana := -1
				if selectionText != "var." {
					mana = TibiadataStringToIntegerV3(selectionText)
				}

				spellBuilder.Mana = mana
			case 5:
				price := 0
				if selectionText != "free" {
					price = TibiadataStringToIntegerV3(selectionText)
				}

				spellBuilder.Price = price
			case 6:
				if selectionText == "yes" {
					spellBuilder.PremiumOnly = true
				}
			}
		})

		SpellsData = append(SpellsData, spellBuilder)
	})

	// adding readable SpellsVocationFilter field
	if vocationName == "" {
		vocationName = "all"
	}

	// Build the data-blob
	return SpellsOverviewResponse{
		Spells{
			SpellsVocationFilter: vocationName,
			Spells:               SpellsData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
