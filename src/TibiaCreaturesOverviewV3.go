package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaCreaturesOverviewV3 func
func TibiaCreaturesOverviewV3(c *gin.Context) {

	// Child of Creatures (used for list of creatures and boosted section)
	type Creature struct {
		Name     string `json:"name"`
		Race     string `json:"race"`
		ImageURL string `json:"image_url"`
		Featured bool   `json:"featured"`
	}

	// Child of JSONData
	type Creatures struct {
		Boosted   Creature   `json:"boosted"`
		Creatures []Creature `json:"creature_list"`
	}

	//
	// The base includes two levels: Creatures and Information
	type JSONData struct {
		Creatures   Creatures   `json:"creatures"`
		Information Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest = map[string]map[string]string{
		"request": {
			"method": "GET",
			"url":    "https://www.tibia.com/library/?subtopic=creatures",
		}}
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find(".InnerTableContainer p").First().Html()
	if err != nil {
		log.Fatal(err)
	}

	// Regex to get data for name and race param for boosted creature
	regex1b := regexp.MustCompile(`<a.*race=(.*)".*?>(.*)</a>`)
	subma1b := regex1b.FindAllStringSubmatch(InnerTableContainerTMPB, -1)
	// Settings vars for usage in JSONData
	BoostedCreatureName := subma1b[0][2]
	BoostedCreatureRace := subma1b[0][1]

	// Regex to get image of boosted creature
	regex2b := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	subma2b := regex2b.FindAllStringSubmatch(InnerTableContainerTMPB, -1)
	// Settings vars for usage in JSONData
	BoostedCreatureImage := subma2b[0][1]

	// Creating empty CreaturesData var
	var CreaturesData []Creature

	// Running query over each div
	ReaderHTML.Find(".BoxContent div div").Each(func(index int, s *goquery.Selection) {

		// Storing HTML into CreatureDivHTML
		CreatureDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex to get data for name, race and img src param for creature
		regex1 := regexp.MustCompile(`.*race=(.*)"><img src="(.*)" border.*div>(.*)<\/div>`)
		subma1 := regex1.FindAllStringSubmatch(CreatureDivHTML, -1)

		// check if regex return length is over 0 and the match of name is over 1
		if len(subma1) > 0 && len(subma1[0][3]) > 1 {

			// Adding bool to indicate features in creature_list
			FeaturedRace := false
			if subma1[0][1] == BoostedCreatureRace {
				FeaturedRace = true
			}

			// Creating data block to return
			CreaturesData = append(CreaturesData, Creature{
				Name:     TibiaDataSanitizeEscapedString(subma1[0][3]),
				Race:     subma1[0][1],
				ImageURL: subma1[0][2],
				Featured: FeaturedRace,
			})

		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		Creatures{
			Boosted: Creature{
				Name:     BoostedCreatureName,
				Race:     BoostedCreatureRace,
				ImageURL: BoostedCreatureImage,
				Featured: true,
			},
			Creatures: CreaturesData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaCreaturesOverviewV3", jsonData)
}
