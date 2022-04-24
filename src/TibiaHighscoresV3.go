package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/TibiaData/tibiadata-api-go/src/validation"
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
type HighscoresResponse struct {
	Highscores  Highscores  `json:"highscores"`
	Information Information `json:"information"`
}

var (
	HighscoresAgeRegex = regexp.MustCompile(`.*<div class="Text">Highscores.*Last Update: ([0-9]+) minutes ago.*`)
	SevenColumnRegex   = regexp.MustCompile(`<td>.*<\/td><td.*">(.*)<\/a><\/td><td.*>(.*)<\/td><td.*>(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
	SixColumnRegex     = regexp.MustCompile(`<td>.*<\/td><td.*">(.*)<\/a><\/td><td.*">(.*)<\/td><td>(.*)<\/td><td.*>(.*)<\/td><td.*>(.*)<\/td>`)
)

func TibiaHighscoresV3Impl(world string, category validation.HighscoreCategory, vocationName string, BoxContentHTML string) (*HighscoresResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaHighscoresV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty HighscoreData var
	var (
		HighscoreData                                                           []Highscore
		HighscoreDataVocation, HighscoreDataWorld, HighscoreDataTitle           string
		HighscoreDataRank, HighscoreDataLevel, HighscoreDataValue, HighscoreAge int
	)

	// getting age of data
	subma1 := HighscoresAgeRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)

	if len(subma1) > 0 {
		HighscoreAge = TibiaDataStringToIntegerV3(subma1[0][1])
	}

	var insideError error

	// Running query over each div
	ReaderHTML.Find(".TableContent tr").First().NextAll().EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		HighscoreDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaHighscoresV3Impl failed at HighscoreDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Regex the data table..
		var subma1 [][]string

		/*
			Tibia highscore table columns
			Achievment	=>	Rank		Name	Vocation	World		Level	Points
			Axe 		=>	Rank		Name	Vocation	World		Level	Skill Level
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

		if category == validation.HighScoreLoyaltypoints {
			subma1 = SevenColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		} else {
			subma1 = SixColumnRegex.FindAllStringSubmatch(HighscoreDivHTML, -1)
		}

		if len(subma1) > 0 {
			HighscoreDataRank++
			if category == validation.HighScoreLoyaltypoints {
				HighscoreDataTitle = subma1[0][2]
				HighscoreDataVocation = subma1[0][3]
				HighscoreDataWorld = subma1[0][4]
				HighscoreDataLevel = TibiaDataStringToIntegerV3(subma1[0][5])
				HighscoreDataValue = TibiaDataStringToIntegerV3(subma1[0][6])
			} else {
				HighscoreDataVocation = subma1[0][2]
				HighscoreDataWorld = subma1[0][3]
				HighscoreDataLevel = TibiaDataStringToIntegerV3(subma1[0][4])
				HighscoreDataValue = TibiaDataStringToIntegerV3(subma1[0][5])
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

		return true
	})

	categoryString, _ := category.String()

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &HighscoresResponse{
		Highscores{
			World:         strings.Title(strings.ToLower(world)),
			Category:      categoryString,
			Vocation:      vocationName,
			HighscoreAge:  HighscoreAge,
			HighscoreList: HighscoreData,
		},
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, err
}
