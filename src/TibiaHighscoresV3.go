package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaHighscoresV3 func
func TibiaHighscoresV3(c *gin.Context) {

	// getting params from URL
	world := c.Param("world")
	category := c.Param("category")
	vocation := c.Param("vocation")

	// do some validation of category and vocation
	// maybe return error on faulty value?!

	// Child of Highscores
	type Highscore struct {
		Rank     int    `json:"rank"`            // Rank column
		Name     string `json:"name"`            // Name column
		Vocation string `json:"vocation"`        // Vocation column
		World    string `json:"world"`           // World column
		Level    int    `json:"level"`           // Level column
		Value    int    `json:"value"`           // Points/SkillLevel column
		Title    string `json:"title,omitempty"` // Title column (when category: loyalty)
	}

	// Child of JSONData
	type Highscores struct {
		World         string      `json:"world"`
		Category      string      `json:"category"`
		Vocation      string      `json:"vocation"`
		HighscoreAge  int         `json:"highscore_age"`
		HighscoreList []Highscore `json:"highscore_list"`
	}

	//
	// The base includes two levels: Highscores and Information
	type JSONData struct {
		Highscores  Highscores  `json:"highscores"`
		Information Information `json:"information"`
	}

	// Adding fix for First letter to be upper and rest lower
	if strings.EqualFold(world, "all") {
		world = ""
	} else {
		world = TibiadataStringWorldFormatToTitleV3(world)
	}

	// Sanatize of category value
	category = strings.ToLower(category)
	var categoryid string
	categoryid = "6"
	if len(category) > 0 {
		switch category {
		case "achievements", "achievement":
			category = "achievements"
			categoryid = "1"
		case "axe", "axefighting":
			category = "axefighting"
			categoryid = "2"
		case "charm", "charms", "charmpoints":
			category = "charmpoints"
			categoryid = "3"
		case "club", "clubfighting":
			category = "clubfighting"
			categoryid = "4"
		case "distance", "distancefighting":
			category = "distancefighting"
			categoryid = "5"
		case "fishing":
			category = "fishing"
			categoryid = "7"
		case "fist", "fistfighting":
			category = "fistfighting"
			categoryid = "8"
		case "goshnar", "goshnars", "goshnarstaint":
			category = "goshnarstaint"
			categoryid = "9"
		case "loyalty", "loyaltypoints":
			category = "loyaltypoints"
			categoryid = "10"
		case "magic", "mlvl", "magiclevel":
			category = "magiclevel"
			categoryid = "11"
		case "shielding", "shield":
			category = "shielding"
			categoryid = "12"
		case "sword", "swordfighting":
			category = "swordfighting"
			categoryid = "13"
		case "drome", "dromescore":
			category = "dromescore"
			categoryid = "14"
		default:
			category = "experience"
		}
	} else {
		category = "experience"
	}

	// Sanitize of vocation input
	vocationName, vocationid := TibiaDataVocationValidator(vocation)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=highscores&world=" + TibiadataQueryEscapeStringV3(world) + "&category=" + TibiadataQueryEscapeStringV3(categoryid) + "&profession=" + TibiadataQueryEscapeStringV3(vocationid) + "&currentpage=400000000000000"
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g.1 for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusServiceUnavailable, "TibiaHighscoresV3", gin.H{"error": err.Error()})
		return
	}

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty HighscoreData var
	var HighscoreData []Highscore
	var HighscoreDataVocation, HighscoreDataWorld, HighscoreDataTitle string
	var HighscoreDataRank, HighscoreDataLevel, HighscoreDataValue, HighscoreAge int

	// getting age of data
	regex1 := regexp.MustCompile(`.*<div class="Text">Highscores.*Last Update: ([0-9]+) minutes ago.*`)
	subma1 := regex1.FindAllStringSubmatch(string(BoxContentHTML), 1)
	HighscoreAge = TibiadataStringToIntegerV3(subma1[0][1])

	// Running query over each div
	ReaderHTML.Find(".TableContent tr").First().NextAll().Each(func(index int, s *goquery.Selection) {

		// Storing HTML into CreatureDivHTML
		HighscoreDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex the data table..
		var subma1 [][]string

		/*
			Tibia highscore table columns
			Achievment	=>	Rank		Name	Vocation	World		Level	Points
			Axe			=>	Rank		Name	Vocation	World		Level	Skill Level
			Charm		=>	Rank		Name	Vocation	World		Level	Points
			Club		=>	Rank		Name	Vocation	World		Level	Skill Level
			Distance	=>	Rank		Name	Vocation	World		Level	Skill Level
			Drome		=>	Rank		Name	Vocation	World		Level	Score
			Exp			=>	Rank		Name	Vocation	World		Level	Points
			Fishing		=>	Rank		Name	Vocation	World		Level	Skill Level
			Fist		=>	Rank		Name	Vocation	World		Level	Skill Level
			Goshnar		=>	Rank		Name	Vocation	World		Level	Points
			Loyalty		=>	Rank		Name	Title		Vocation	World	Level			Points
			Magic lvl	=>	Rank		Name	Vocation	World		Level	Skill Level
			Shield		=>	Rank		Name	Vocation	World		Level	Skill Level
			Sword		=>	Rank		Name	Vocation	World		Level	Skill Level
		*/

		if category == "loyaltypoints" {
			// Regex when highscore has 7 columns
			regex1 := regexp.MustCompile(`<td>.*<\/td><td.*">(.*)<\/a><\/td><td.*>(.*)<\/td><td.*>(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
			subma1 = regex1.FindAllStringSubmatch(HighscoreDivHTML, -1)
		} else {
			// Regex when highscore has 6 columns (category except lojalty)
			regex1 := regexp.MustCompile(`<td>.*<\/td><td.*">(.*)<\/a><\/td><td.*">(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
			subma1 = regex1.FindAllStringSubmatch(HighscoreDivHTML, -1)
		}

		if len(subma1) > 0 {

			// Debugging of what is in which column
			if TibiadataDebug {
				log.Println("1 -> " + subma1[0][1])
				log.Println("2 -> " + subma1[0][2])
				log.Println("3 -> " + subma1[0][3])
				log.Println("4 -> " + subma1[0][4])
				log.Println("5 -> " + subma1[0][5])
				if category == "loyaltypoints" {
					log.Println("6 -> " + subma1[0][6])
				}
			}

			HighscoreDataRank++
			if category == "loyaltypoints" {
				HighscoreDataTitle = subma1[0][2]
				HighscoreDataVocation = subma1[0][3]
				HighscoreDataWorld = subma1[0][4]
				HighscoreDataLevel = TibiadataStringToIntegerV3(subma1[0][5])
				HighscoreDataValue = TibiadataStringToIntegerV3(subma1[0][6])
			} else {
				HighscoreDataVocation = subma1[0][2]
				HighscoreDataWorld = subma1[0][3]
				HighscoreDataLevel = TibiadataStringToIntegerV3(subma1[0][4])
				HighscoreDataValue = TibiadataStringToIntegerV3(subma1[0][5])
			}

			HighscoreData = append(HighscoreData, Highscore{
				Rank:     HighscoreDataRank,
				Name:     TibiaDataSanitizeEscapedString(subma1[0][1]),
				Vocation: HighscoreDataVocation,
				World:    HighscoreDataWorld,
				Level:    HighscoreDataLevel,
				Value:    HighscoreDataValue,
				Title:    HighscoreDataTitle,
			})

		}
	})

	// Printing the HighscoreData data to log
	if TibiadataDebug {
		log.Println(HighscoreData)
	}

	//
	// Build the data-blob
	jsonData := JSONData{
		Highscores{
			World:         strings.Title(strings.ToLower(world)),
			Category:      category,
			Vocation:      vocationName,
			HighscoreAge:  HighscoreAge,
			HighscoreList: HighscoreData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaHighscoresV3", jsonData)
}
