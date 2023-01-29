package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of BoostableBoss (used for list of boostable bosses and boosted boss section)
type OverviewBoostableBoss struct {
	Name     string `json:"name"`      // The name of the boss.
	ImageURL string `json:"image_url"` // The URL to this boss’s image.
	Featured bool   `json:"featured"`  // Whether it is featured of not.
}

// Child of JSONData
type BoostableBossesContainer struct {
	Boosted         OverviewBoostableBoss   `json:"boosted"`             // The current boosted boss.
	BoostableBosses []OverviewBoostableBoss `json:"boostable_boss_list"` // The list of boostable bosses.
}

// The base includes two levels: BoostableBosses and Information
type BoostableBossesOverviewResponse struct {
	BoostableBosses BoostableBossesContainer `json:"boostable_bosses"`
	Information     Information              `json:"information"`
}

var (
	BoostedBossNameRegex          = regexp.MustCompile(`<b>(.*)</b>`)
	BoostedBossImageRegex         = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	BoostableBossInformationRegex = regexp.MustCompile(`<img src="(.*)" border.*div>(.*)<\/div>`)
)

func TibiaBoostableBossesOverviewImpl(BoxContentHTML string) (*BoostableBossesOverviewResponse, error) {
	// Creating empty vars
	var (
		BoostedBossName, BoostedBossImage string
	)
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaBoostableBossesOverviewImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Getting data from div.InnerTableContainer and then first p
	InnerTableContainerTMPB, err := ReaderHTML.Find(".InnerTableContainer p").First().Html()
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaBoostableBossesOverviewImpl failed at ReaderHTML.Find, error: %s", err)
	}

	// Regex to get data for name for boosted boss
	subma1b := BoostedBossNameRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)

	if len(subma1b) > 0 {
		// Settings vars for usage in JSONData
		BoostedBossName = subma1b[0][1]
	}

	// Regex to get image of boosted boss
	subma2b := BoostedBossImageRegex.FindAllStringSubmatch(InnerTableContainerTMPB, -1)

	if len(subma2b) > 0 {
		// Settings vars for usage in JSONData
		BoostedBossImage = subma2b[0][1]
	}

	// Creating empty BoostableBossesData var
	var BoostableBossesData []OverviewBoostableBoss

	var insideError error

	// Running query over each div
	ReaderHTML.Find(".BoxContent div div").EachWithBreak(func(index int, s *goquery.Selection) bool {

		// Storing HTML into BoostableBossDivHTML
		BoostableBossDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaBoostableBossesOverviewImpl failed at BoostableBossDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		// Regex to get data for name, race and img src param for creature
		subma1 := BoostableBossInformationRegex.FindAllStringSubmatch(BoostableBossDivHTML, -1)

		// check if regex return length is over 0 and the match of name is over 1
		if len(subma1) > 0 && len(subma1[0][2]) > 1 {
			// Adding bool to indicate features in boostable_boss_list
			FeaturedRace := false
			if subma1[0][2] == BoostedBossName {
				FeaturedRace = true
			}

			// Creating data block to return
			BoostableBossesData = append(BoostableBossesData, OverviewBoostableBoss{
				Name:     TibiaDataSanitizeEscapedString(subma1[0][2]),
				ImageURL: subma1[0][1],
				Featured: FeaturedRace,
			})
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	// Build the data-blob
	return &BoostableBossesOverviewResponse{
		BoostableBossesContainer{
			Boosted: OverviewBoostableBoss{
				Name:     TibiaDataSanitizeEscapedString(BoostedBossName),
				ImageURL: BoostedBossImage,
				Featured: true,
			},
			BoostableBosses: BoostableBossesData,
		},
		Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
