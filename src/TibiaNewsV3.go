package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of JSONData
type News struct {
	ID          int    `json:"id"`
	Date        string `json:"date"`
	Title       string `json:"title,omitempty"`
	Category    string `json:"category"`
	Type        string `json:"type,omitempty"`
	TibiaURL    string `json:"url"`
	Content     string `json:"content"`
	ContentHTML string `json:"content_html"`
}

//
// The base
type NewsResponse struct {
	News        News        `json:"news"`
	Information Information `json:"information"`
}

var (
	martelRegex = regexp.MustCompile(`<img src=\"https:\/\/static\.tibia\.com\/images\/global\/letters\/letter_martel_(.)\.gif\" ([^\/>]+..)`)
)

func TibiaNewsV3Impl(NewsID int, rawUrl string, BoxContentHTML string) (*NewsResponse, error) {
	// Declaring vars for later use..
	var (
		NewsData    News
		tmp1        *goquery.Selection
		tmp2        string
		insideError error
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaNewsV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	NewsData.ID = NewsID
	NewsData.TibiaURL = rawUrl

	ReaderHTML.Find(".NewsHeadline").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// getting category by image src
		CategoryImg, _ := s.Find("img").Attr("src")
		NewsData.Category = TibiaDataGetNewsCategory(CategoryImg)

		// getting date from headline
		tmp1 = s.Find(".NewsHeadlineDate")
		tmp2, err = tmp1.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaNewsV3Impl failed at tmp2, err = tmp1.Html(), NewsHeadlineDate, err: %s", err)
			return false
		}

		NewsData.Date = TibiaDataDateV3(strings.ReplaceAll(tmp2, " - ", ""))

		// getting headline text (which could be title or also type)
		tmp1 = s.Find(".NewsHeadlineText")
		tmp2, err = tmp1.Html()
		if err != nil {
			insideError = fmt.Errorf("[error] TibiaNewsV3Impl failed at tmp2, err = tmp1.Html(), NewsHeadlineText, err: %s", err)
			return false
		}

		NewsData.Title = RemoveHtmlTag(tmp2)
		if NewsData.Title == "News Ticker" {
			NewsData.Type = "ticker"
			NewsData.Title = ""
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	ReaderHTML.Find(".NewsTableContainer").EachWithBreak(func(index int, s *goquery.Selection) bool {
		// checking if its a ticker..
		if NewsData.Type == "ticker" {
			tmp1 = s.Find("p")
			NewsData.Content = tmp1.Text()
			NewsData.ContentHTML, err = tmp1.Html()
			if err != nil {
				insideError = fmt.Errorf("[error] TibiaNewsV3Impl failed at NewsData.ContentHTML, err = tmp1.Html(), err: %s", err)
				return false
			}
		} else {
			// getting html
			tmp2, err = s.First().Html()
			if err != nil {
				insideError = fmt.Errorf("[error] TibiaNewsV3Impl failed at NewsData.ContentHTML, tmp2, err = s.First().Html(), err: %s", err)
				return false
			}

			// replacing martel letter in articles with real letters
			tmp2 = martelRegex.ReplaceAllString(tmp2, "$1")
			s.ReplaceWithHtml(tmp2)

			// storing html content
			NewsData.ContentHTML = tmp2

			// reading string again after replacing letters
			tmp1, err := goquery.NewDocumentFromReader(strings.NewReader(tmp2))
			if err != nil {
				insideError = fmt.Errorf("[error] TibiaNewsV3Impl failed attmp1, err := goquery.NewDocumentFromReader, err: %s", err)
				return false
			}

			// storing text content
			NewsData.Content = tmp1.Text()
		}

		return true
	})

	if insideError != nil {
		return nil, insideError
	}

	//
	// Build the data-blob
	return &NewsResponse{
		NewsData,
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
