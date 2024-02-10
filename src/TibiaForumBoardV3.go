package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ForumBoardThreadLastPost struct {
	ID            int    `json:"id"`
	PostedAt      string `json:"posted_at"`
	CharacterName string `json:"character_name"`
}

type ForumBoardThread struct {
	ID            int                      `json:"id"`
	Title         string                   `json:"title"`
	IsHot         bool                     `json:"is_hot"`
	IsNew         bool                     `json:"is_new"`
	IsClosed      bool                     `json:"is_closed"`
	IsSticky      bool                     `json:"is_sticky"`
	Icon          string                   `json:"icon"`
	Pages         int                      `json:"pages"`
	Replies       int                      `json:"replies"`
	Views         int                      `json:"views"`
	CharacterName string                   `json:"character_name"`
	LastPost      ForumBoardThreadLastPost `json:"last_post"`
}

type ForumBoardPagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_records"`
}

type ForumBoard struct {
	ID                    int                  `json:"id"`
	SectionName           string               `json:"section_name"`
	BoardName             string               `json:"board_name"`
	ThreadsAge            string               `json:"threads_age"`
	ThreadsList           []ForumBoardThread   `json:"threads"`
	BoardsBoardPagination ForumBoardPagination `json:"pagination"`
}

type ForumBoardResponse struct {
	Board       ForumBoard  `json:"board"`
	Information Information `json:"information"`
}

var (
	threadIdAndTitleRegex = regexp.MustCompile(`threadid=([0-9]+)">(.*)<\/a>\S(<br\/>|<\/td>)`)
	statusIconRegex       = regexp.MustCompile(`.*<td class="HNCContainer"><img src="(.*)" border="0" width="22" height="22"/></td>`)
	threadIconRegex       = regexp.MustCompile(`.*<td class="HNCContainer"><img src=".*" border="0" width="15" height="15" alt="(.*)"/></td>`)
	pagesRegex            = regexp.MustCompile(`pagenumber=([0-9]+)"`)
	repliesViewsRegex     = regexp.MustCompile(`<td class="TextRight">([0-9,]+)</td><td class="TextRight">([0-9,]+)</td>`)
	threadStarterRegex    = regexp.MustCompile(`<div class="ThreadStarter">.*<a.*">(.*)</a></div>`)
	lastPostInfoRegex     = regexp.MustCompile(`<td><span class="LastPostInfo">.*postid=([0-9]+)#.*height="9"/></a>(.*)</span><span><font class="ff_info">.*name=.*">(.*)</a></span>`)
	paginationRegex       = regexp.MustCompile(`.*<b>.*Pages:\ ?(.*)<\/b>.*<b>.*Results:\ ?([0-9,]+)<\/b>.*`)
	boardBreadcrumbsRegex = regexp.MustCompile(`.*<div class=\"ForumBreadCrumbs\"><a.*">(.*)<\/a> \| <b>(.*)</b></div>.*`)
)

// TibiaForumBoardV3Impl func
func TibiaForumBoardV3Impl(BoardID string, BoxContentHTML string, currentPage int, threadsAge ThreadsAge) ForumBoardResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	var (
		ThreadsData              []ForumBoardThread
		totalPages, totalResults int
		sectionName, boardName   string
	)

	// getting amount of pages
	subma1 := boardBreadcrumbsRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		sectionName = subma1[0][1]
		boardName = subma1[0][2]
	}

	// getting amount of pages
	subma1 = paginationRegex.FindAllStringSubmatch(string(BoxContentHTML), 1)
	if len(subma1) > 0 {
		totalPages = strings.Count(subma1[0][1], "class=\"PageLink")
		totalResults = TibiaDataStringToIntegerV3(subma1[0][2])
	}

	// Running query over each div
	ReaderHTML.Find(".TableContentContainer .TableContent#ThreadOverview tbody tr").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CreatureDivHTML
		BoardDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		subma1 := threadIdAndTitleRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma1) == 0 {
			return
		}

		thread := ForumBoardThread{
			ID:    TibiaDataStringToIntegerV3(subma1[0][1]),
			Title: TibiaDataSanitizeStrings(subma1[0][2]),
		}

		subma2 := statusIconRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma2) > 0 {
			if strings.Contains(subma2[0][1], "hot") {
				thread.IsHot = true
			}

			if strings.Contains(subma2[0][1], "new") {
				thread.IsNew = true
			}

			if strings.Contains(subma2[0][1], "closed") {
				thread.IsClosed = true
			}

			if strings.Contains(subma2[0][1], "sticky") {
				thread.IsSticky = true
			}
		}

		subma3 := threadIconRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma3) > 0 {
			thread.Icon = subma3[0][1]
		}

		subma4 := pagesRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma4) > 0 {
			thread.Pages = TibiaDataStringToIntegerV3(subma4[len(subma4)-1][1])
		}

		subma5 := repliesViewsRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma5) > 0 {
			thread.Replies = TibiaDataStringToIntegerV3(subma5[0][1])
			thread.Views = TibiaDataStringToIntegerV3(subma5[0][2])
		}

		subma6 := threadStarterRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma6) > 0 {
			thread.CharacterName = subma6[0][1]
		}

		subma7 := lastPostInfoRegex.FindAllStringSubmatch(BoardDivHTML, -1)
		if len(subma7) > 0 {
			thread.LastPost = ForumBoardThreadLastPost{
				ID:            TibiaDataStringToIntegerV3(subma7[0][1]),
				PostedAt:      TibiaDataDatetimeV3(strings.Trim(TibiaDataSanitizeStrings(subma7[0][2]), " ")),
				CharacterName: subma7[0][3],
			}
		}

		ThreadsData = append(ThreadsData, thread)
	})

	threadsAgeString, _ := threadsAge.String()

	//
	// Build the data-blob
	return ForumBoardResponse{
		Board: ForumBoard{
			ID:          TibiaDataStringToIntegerV3(BoardID),
			SectionName: sectionName,
			BoardName:   boardName,
			ThreadsAge:  threadsAgeString,
			ThreadsList: ThreadsData,
			BoardsBoardPagination: ForumBoardPagination{
				CurrentPage:  currentPage,
				TotalPages:   totalPages,
				TotalResults: totalResults,
			},
		},
		Information: Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
