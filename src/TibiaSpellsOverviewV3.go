package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// TibiaSpellsOverviewV3 func
func TibiaSpellsOverviewV3Impl(vocationName string, BoxContentHTML string) (*SpellsOverviewResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaSpellsOverviewV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
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
				spellBuilder.Level = TibiaDataStringToIntegerV3(selectionText)
			case 4:
				mana := -1
				if selectionText != "var." {
					mana = TibiaDataStringToIntegerV3(selectionText)
				}

				spellBuilder.Mana = mana
			case 5:
				price := 0
				if selectionText != "free" {
					price = TibiaDataStringToIntegerV3(selectionText)
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
	return &SpellsOverviewResponse{
		Spells{
			SpellsVocationFilter: vocationName,
			Spells:               SpellsData,
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
