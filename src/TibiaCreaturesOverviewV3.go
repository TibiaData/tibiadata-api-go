package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of Creatures (used for list of creatures and boosted section)
type OverviewCreature struct {
	Name     string `json:"name"`
	Race     string `json:"race"`
	ImageURL string `json:"image_url"`
	Featured bool   `json:"featured"`
}

// Child of JSONData
type CreaturesContainer struct {
	Boosted   OverviewCreature   `json:"boosted"`
	Creatures []OverviewCreature `json:"creature_list"`
}

//
// The base includes two levels: Creatures and Information
type CreaturesOverviewResponse struct {
	Creatures   CreaturesContainer `json:"creatures"`
	Information Information        `json:"information"`
}

var (
	BoostedCreatureNameAndRaceRegex = regexp.MustCompile(`<a.*race=(.*)".*?>(.*)</a>`)
	BoostedCreatureImageRegex       = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	CreatureInformationRegex        = regexp.MustCompile(`.*race=(.*)"><img src="(.*)" border.*div>(.*)<\/div>`)
)

func TibiaCreaturesOverviewV3Impl(BoxContentHTML string) (*CreaturesOverviewResponse, error) {
	var (
		BoostedCreatureName, BoostedCreatureRace, BoostedCreatureImage string
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaCreaturesOverviewV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find(".InnerTableContainer p").First().Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaCreaturesOverviewV3Impl failed at ReaderHTML.Find, err: %s", err)
	}

	// Regex to get data for name and race param for boosted creature
	subma1b := BoostedCreatureNameAndRaceRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)

	if len(subma1b) > 0 {
		// Settings vars for usage in JSONData
		BoostedCreatureName = subma1b[0][2]
		BoostedCreatureRace = subma1b[0][1]
	}

	// Regex to get image of boosted creature
	subma2b := BoostedCreatureImageRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)

	if len(subma2b) > 0 {
		// Settings vars for usage in JSONData
		BoostedCreatureImage = subma2b[0][1]
	}

	var (
		// Creating empty CreaturesData var
		CreaturesData []OverviewCreature

		// Creating empty error var
		insideError error
	)

	// Running query over each div
	ReaderHTML.Find(".BoxContent div div").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		CreatureDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaCreaturesOverviewV3Impl failed at CreatureDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Regex to get data for name, race and img src param for creature
		subma1 := CreatureInformationRegex.FindAllStringSubmatch(CreatureDivHTML, -1)

		// check if regex return length is over 0 and the match of name is over 1
		if len(subma1) > 0 && len(subma1[0][3]) > 1 {
			// Adding bool to indicate features in creature_list
			FeaturedRace := false
			if subma1[0][1] == BoostedCreatureRace {
				FeaturedRace = true
			}

			// Creating data block to return
			CreaturesData = append(CreaturesData, OverviewCreature{
				Name:     TibiaDataSanitizeEscapedString(subma1[0][3]),
				Race:     subma1[0][1],
				ImageURL: subma1[0][2],
				Featured: FeaturedRace,
			})
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	// Build the data-blob
	return &CreaturesOverviewResponse{
		CreaturesContainer{
			Boosted: OverviewCreature{
				Name:     BoostedCreatureName,
				Race:     BoostedCreatureRace,
				ImageURL: BoostedCreatureImage,
				Featured: true,
			},
			Creatures: CreaturesData,
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
