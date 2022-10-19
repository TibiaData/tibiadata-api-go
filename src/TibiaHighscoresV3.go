package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

// Child of Highscore
type HighscorePage struct {
	CurrentPage     int `json:"current_page"`  // Current page
	TotalPages      int `json:"total_pages"`   // Total page count
	TotalHighscores int `json:"total_records"` // Total highscore records
}

// Child of JSONData
type Highscores struct {
	World         string        `json:"world"`
	Category      string        `json:"category"`
	Vocation      string        `json:"vocation"`
	HighscoreAge  int           `json:"highscore_age"`
	HighscoreList []Highscore   `json:"highscore_list"`
	HighscorePage HighscorePage `json:"highscore_page"`
}

// The base includes two levels: Highscores and Information
type HighscoresResponse struct {
	Highscores  Highscores  `json:"highscores"`
	Information Information `json:"information"`
}

var (
	HighscoresAgeRegex  = regexp.MustCompile(`.*<div class="Text">Highscores.*Last Update: ([0-9]+) minutes ago.*`)
	HighscoresPageRegex = regexp.MustCompile(`.*<b>\&raquo; Pages:\ ?(.*)<\/b>.*<b>\&raquo; Results:\ ?([0-9,]+)<\/b>.*`)
	SevenColumnRegex    = regexp.MustCompile(`<td>(.*)<\/td><td.*">(.*)<\/a><\/td><td.*>(.*)<\/td><td.*>(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
	SixColumnRegex      = regexp.MustCompile(`<td>(.*)<\/td><td.*">(.*)<\/a><\/td><td.*">(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
)

func TibiaHighscoresV3Impl(world string, category HighscoreCategory, vocationName string, currentPage int, BoxContentHTML string) HighscoresResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty HighscoreData var
	var (
		HighscoreData                                                                                                          []Highscore
		HighscoreDataVocation, HighscoreDataWorld, HighscoreDataTitle                                                          string
		HighscoreDataRank, HighscoreDataLevel, HighscoreDataValue, HighscoreAge, HighscoreTotalPages, HighscoreTotalHighscores int
	)

	// getting age of data
	subma1 := HighscoresAgeRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		HighscoreAge = TibiaDataStringToIntegerV3(subma1[0][1])
	}

	// getting amount of pages
	subma1 = HighscoresPageRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		HighscoreTotalPages = strings.Count(subma1[0][1], "class=\"PageLink")
		HighscoreTotalHighscores = TibiaDataStringToIntegerV3(subma1[0][2])
	}

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
			Axe 		=>	Rank		Name	Vocation	World		Level	Skill Level
			Boss		=>	Rank		Name	Vocation	World		Level	Points
			Charm		=>	Rank		Name	Vocation	World		Level	Points
			Club		=>	Rank		Name	Vocation	World		Level	Skill Level
			Distance	=>	Rank		Name	Vocation	World		Level	Skill Level
			Drome		=>	Rank		Name	Vocation	World		Level	Score
			Exp 		=>	Rank		Name	Vocation	World		Level	Points
			Fishing		=>	Rank		Name	Vocation	World		Level	Skill Level
			Fist		=>	Rank		Name	Vocation	World		Level	Skill Level
			Goshnar		=>	Rank		Name	Vocation	World		Level	Points
			Loyalty		=>	Rank		Name	Title		Vocation	World	Level		Points
			Magic lvl	=>	Rank		Name	Vocation	World		Level	Skill Level
			Shield		=>	Rank		Name	Vocation	World		Level	Skill Level
			Sword		=>	Rank		Name	Vocation	World		Level	Skill Level
		*/

		if category == loyaltypoints {
			subma1 = SevenColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		} else {
			subma1 = SixColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		}

		if len(subma1) > 0 {

			HighscoreDataRank = TibiaDataStringToIntegerV3(subma1[0][1])
			if category == loyaltypoints {
				HighscoreDataTitle = subma1[0][3]
				HighscoreDataVocation = subma1[0][4]
				HighscoreDataWorld = subma1[0][5]
				HighscoreDataLevel = TibiaDataStringToIntegerV3(subma1[0][6])
				HighscoreDataValue = TibiaDataStringToIntegerV3(subma1[0][7])
			} else {
				HighscoreDataVocation = subma1[0][3]
				HighscoreDataWorld = subma1[0][4]
				HighscoreDataLevel = TibiaDataStringToIntegerV3(subma1[0][5])
				HighscoreDataValue = TibiaDataStringToIntegerV3(subma1[0][6])
			}

			HighscoreData = append(HighscoreData, Highscore{
				Rank:     HighscoreDataRank,
				Name:     TibiaDataSanitizeEscapedString(subma1[0][2]),
				Vocation: HighscoreDataVocation,
				World:    HighscoreDataWorld,
				Level:    HighscoreDataLevel,
				Value:    HighscoreDataValue,
				Title:    HighscoreDataTitle,
			})

		}
	})

	categoryString, _ := category.String()

	//
	// Build the data-blob
	return HighscoresResponse{
		Highscores{
			World:         strings.Title(strings.ToLower(world)),
			Category:      categoryString,
			Vocation:      vocationName,
			HighscoreAge:  HighscoreAge,
			HighscoreList: HighscoreData,
			HighscorePage: HighscorePage{
				CurrentPage:     currentPage,
				TotalPages:      HighscoreTotalPages,
				TotalHighscores: HighscoreTotalHighscores,
			},
		},
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
