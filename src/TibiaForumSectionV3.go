package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SectionBoardLastPost struct {
	ID            int    `json:"id"`
	PostedAt      string `json:"posted_at"`
	CharacterName string `json:"character_name"`
}

type SectionBoard struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Posts       int                  `json:"posts"`
	Threads     int                  `json:"threads"`
	LastPost    SectionBoardLastPost `json:"last_post"`
}

type ForumSectionResponse struct {
	Boards      []SectionBoard `json:"boards"`
	Information Information    `json:"information"`
}

var (
	boardInformationRegex      = regexp.MustCompile(`.*boardid=(.*)">(.*)<\/a><br\/><font class="ff_info">(.*)<\/font><\/td><td class="TextRight">(.*)<\/td><td class="TextRight">(.*)<\/td><td><span class="LastPostInfo">`)
	lastPostIdRegex            = regexp.MustCompile(`.*postid=(.*)#post`)
	lastPostPostedAtRegex      = regexp.MustCompile(`.*height="9"\/><\/a>(.*)<\/span><span><font class="ff_info">`)
	lastPostCharacterNameRegex = regexp.MustCompile(`.*subtopic=characters&amp;name=.*\">(.*)<\/a><\/span>`)
)

// TibiaForumSectionV3Impl func
func TibiaForumSectionV3Impl(BoxContentHTML string) ForumSectionResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	var (
		BoardsData                              []SectionBoard
		LastPostId                              int
		LastPostPostedAt, LastPostCharacterName string
	)

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CreatureDivHTML
		BoardsDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		subma1 := boardInformationRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma1) == 0 {
			return
		}

		subma2 := lastPostIdRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma2) > 0 {
			LastPostId = TibiaDataStringToIntegerV3(subma2[0][1])
		}

		subma3 := lastPostPostedAtRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma3) > 0 {
			LastPostPostedAt = TibiaDataDatetimeV3(strings.Trim(TibiaDataSanitizeStrings(subma3[0][1]), " "))
		}

		subma4 := lastPostCharacterNameRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma4) > 0 {
			LastPostCharacterName = TibiaDataSanitizeStrings(subma4[0][1])
		}

		BoardsData = append(BoardsData, SectionBoard{
			ID:          TibiaDataStringToIntegerV3(subma1[0][1]),
			Name:        subma1[0][2],
			Description: subma1[0][3],
			Posts:       TibiaDataStringToIntegerV3(subma1[0][4]),
			Threads:     TibiaDataStringToIntegerV3(subma1[0][5]),
			LastPost: SectionBoardLastPost{
				ID:            LastPostId,
				PostedAt:      LastPostPostedAt,
				CharacterName: LastPostCharacterName,
			},
		})
	})

	//
	// Build the data-blob
	return ForumSectionResponse{
		Boards: BoardsData,
		Information: Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
