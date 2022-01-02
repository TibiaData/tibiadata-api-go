package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaCreaturesCreatureV3 func
func TibiaCreaturesCreatureV3(c *gin.Context) {

	// local strings used in this function
	var localDamageString = " damage"

	// getting params from URL
	race := c.Param("race")

	// Child of JSONData
	type Creature struct {
		Name             string   `json:"name"`
		Race             string   `json:"race"`
		ImageURL         string   `json:"image_url"`
		Description      string   `json:"description"`
		Behaviour        string   `json:"behaviour"`
		Hitpoints        int      `json:"hitpoints"`
		ImmuneTo         []string `json:"immune"`
		StrongAgainst    []string `json:"strong"`
		WeaknessAgainst  []string `json:"weakness"`
		BeParalysed      bool     `json:"be_paralysed"`
		BeSummoned       bool     `json:"be_summoned"`
		SummonMana       int      `json:"summoned_mana"`
		BeConvinced      bool     `json:"be_convinced"`
		ConvincedMana    int      `json:"convinced_mana"`
		SeeInvisible     bool     `json:"see_invisible"`
		ExperiencePoints int      `json:"experience_points"`
		IsLootable       bool     `json:"is_lootable"`
		LootList         []string `json:"loot_list"`
		Featured         bool     `json:"featured"`
	}

	//
	// The base includes two levels: Creature and Information
	type JSONData struct {
		Creature    Creature    `json:"creature"`
		Information Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/library/?subtopic=creatures&race=" + TibiadataQueryEscapeStringV3(race))

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Getting data
	InnerTableContainerTMP1, err := ReaderHTML.Find(".BoxContent div").First().NextAll().Html()
	if err != nil {
		log.Fatal(err)
	}

	// Regex to get data
	regex1 := regexp.MustCompile(`.*;">(.*)<\/h2> <img src="(.*)"\/>.*<p>(.*)<\/p> <p>(.*)<\/p> <p>(.*)<\/p>.*`)
	subma1 := regex1.FindAllStringSubmatch(InnerTableContainerTMP1, -1)

	// Preparing vars
	var CreatureDescription, CreatureBehaviour string
	var CreatureLootList, CreatureImmuneTo, CreatureStrongAgainst, CreatureWeaknessAgainst []string
	var CreatureHitpoints, CreatureSummonedMana, CreatureConvincedMana, CreatureExperiencePoints int
	var CreatureBeParalysed, CreatureBeSummoned, CreatureBeConvinced, CreatureSeeInvisible, CreatureIsLootable bool

	// Preparing data for JSONData
	if len(subma1) > 0 {

		// Description (and unescape hmtl string)
		CreatureDescription = strings.ReplaceAll(subma1[0][3], "<br/>", "\n")
		CreatureDescription = TibiaDataSanitizeEscapedString(CreatureDescription)

		// Behaviour
		// Regex to get data..
		regex2 := regexp.MustCompile(`.*have (.*) hitpoints. (.*)`)
		subma2 := regex2.FindAllStringSubmatch(subma1[0][4], -1)
		// Add data to vars
		CreatureHitpoints = TibiadataStringToIntegerV3(subma2[0][1])
		CreatureBehaviour = subma2[0][2]
		if !strings.Contains(subma1[0][4], "cannot be paralysed") {
			CreatureBeParalysed = true
		}
		if strings.Contains(subma1[0][4], "sense invisible creatures") {
			CreatureSeeInvisible = true
		}
		if strings.Contains(subma1[0][4], " are immune to ") {
			regex21 := regexp.MustCompile(`.*are immune to (.*)`)
			subma21 := regex21.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureImmuneToTmp := strings.Split(subma21[0][1], localDamageString)
			CreatureImmuneTo = strings.Split(strings.Replace(CreatureImmuneToTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], " are strong against ") {
			regex22 := regexp.MustCompile(`.*are strong against (.*)`)
			subma22 := regex22.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureStrongAgainstTmp := strings.Split(subma22[0][1], localDamageString)
			CreatureStrongAgainst = strings.Split(strings.Replace(CreatureStrongAgainstTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], " are weak against ") {
			regex23 := regexp.MustCompile(`.*are weak against (.*)`)
			subma23 := regex23.FindAllStringSubmatch(subma1[0][4], -1)
			CreatureWeaknessAgainstTmp := strings.Split(subma23[0][1], localDamageString)
			CreatureWeaknessAgainst = strings.Split(strings.Replace(CreatureWeaknessAgainstTmp[0], " and ", ", ", 1), ", ")
		}
		if strings.Contains(subma1[0][4], "It takes ") && strings.Contains(subma1[0][4], " mana to ") {
			regex24 := regexp.MustCompile(`.*It takes (.*) mana to (.*)`)
			subma24 := regex24.FindAllStringSubmatch(subma1[0][4], -1)
			subma24_0_2 := subma24[0][2]
			if strings.Contains(subma24_0_2, "convince these creatures but they cannot be") {
				CreatureBeConvinced = true
				CreatureConvincedMana = TibiadataStringToIntegerV3(subma24[0][1])
			} else if strings.Contains(subma24_0_2, "summon or convince these creatures") {
				CreatureBeSummoned = true
				CreatureSummonedMana = TibiadataStringToIntegerV3(subma24[0][1])
				CreatureBeConvinced = true
				CreatureConvincedMana = TibiadataStringToIntegerV3(subma24[0][1])
			}
		}

		// Loot
		// Regex to get loot information
		regex3 := regexp.MustCompile(`.*yield (.*) experience.*carry (.*)with them.`)
		subma3 := regex3.FindAllStringSubmatch(subma1[0][5], -1)
		// Adding data to vars
		CreatureExperiencePoints = TibiadataStringToIntegerV3(subma3[0][1])
		if subma3[0][2] != "nothing" {
			CreatureIsLootable = true
			CreatureLootListTmp := strings.Split(strings.Replace(strings.Replace(subma3[0][2], "items ", "", 1), " and sometimes other ", "", 1), ", ")
			for _, str := range CreatureLootListTmp {
				if str != "" {
					CreatureLootList = append(CreatureLootList, str)
				}
			}
		}
	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Creature{
			Name:             TibiaDataSanitizeEscapedString(subma1[0][1]),
			Race:             race,
			ImageURL:         subma1[0][2],
			Description:      CreatureDescription,
			Behaviour:        CreatureBehaviour,
			Hitpoints:        CreatureHitpoints,
			ImmuneTo:         CreatureImmuneTo,
			StrongAgainst:    CreatureStrongAgainst,
			WeaknessAgainst:  CreatureWeaknessAgainst,
			BeParalysed:      CreatureBeParalysed,
			BeSummoned:       CreatureBeSummoned,
			SummonMana:       CreatureSummonedMana,
			BeConvinced:      CreatureBeConvinced,
			ConvincedMana:    CreatureConvincedMana,
			SeeInvisible:     CreatureSeeInvisible,
			ExperiencePoints: CreatureExperiencePoints,
			IsLootable:       CreatureIsLootable,
			LootList:         CreatureLootList,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaCreaturesCreatureV3", jsonData)
}
