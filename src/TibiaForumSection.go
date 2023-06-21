package main

import (
	"fmt"
	"net/http"
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

// TibiaForumSectionImpl func
func TibiaForumSectionImpl(BoxContentHTML string) (*ForumSectionResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaForumSectionImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	var (
		BoardsData                              []SectionBoard
		LastPostId                              int
		LastPostPostedAt, LastPostCharacterName string

		insideError error
	)

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent tbody tr:not(.LabelH)").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// Storing HTML into CreatureDivHTML
		BoardsDivHTML, err := s.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaForumSectionImpl failed at BoardsDivHTML, err := s.Html(), err: %s", err)
			return false
		}

		subma1 := boardInformationRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma1) == 0 {
			return false
		}

		subma2 := lastPostIdRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma2) > 0 {
			LastPostId = TibiaDataStringToInteger(subma2[0][1])
		}

		subma3 := lastPostPostedAtRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma3) > 0 {
			LastPostPostedAt = TibiaDataDatetime(strings.Trim(TibiaDataSanitizeStrings(subma3[0][1]), " "))
		}

		subma4 := lastPostCharacterNameRegex.FindAllStringSubmatch(BoardsDivHTML, -1)
		if len(subma4) > 0 {
			LastPostCharacterName = TibiaDataSanitizeStrings(subma4[0][1])
		}

		BoardsData = append(BoardsData, SectionBoard{
			ID:          TibiaDataStringToInteger(subma1[0][1]),
			Name:        subma1[0][2],
			Description: subma1[0][3],
			Posts:       TibiaDataStringToInteger(subma1[0][4]),
			Threads:     TibiaDataStringToInteger(subma1[0][5]),
			LastPost: SectionBoardLastPost{
				ID:            LastPostId,
				PostedAt:      LastPostPostedAt,
				CharacterName: LastPostCharacterName,
			},
		})

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &ForumSectionResponse{
		Boards: BoardsData,
		Information: Information{
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
