package main

import (
	"log"
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

func TibiaNewsV3Impl(NewsID int, rawUrl string, BoxContentHTML string) NewsResponse {
	// Declaring vars for later use..
	var (
		NewsData News
		tmp1     *goquery.Selection
		tmp2     string
	)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	NewsData.ID = NewsID
	NewsData.TibiaURL = rawUrl

	ReaderHTML.Find(".NewsHeadline").Each(func(index int, s *goquery.Selection) {

		// getting category by image src
		CategoryImg, _ := s.Find("img").Attr("src")
		NewsData.Category = TibiaDataGetNewsCategory(CategoryImg)

		// getting date from headline
		tmp1 = s.Find(".NewsHeadlineDate")
		tmp2, _ = tmp1.Html()
		NewsData.Date = TibiaDataDateV3(strings.ReplaceAll(tmp2, " - ", ""))

		// getting headline text (which could be title or also type)
		tmp1 = s.Find(".NewsHeadlineText")
		tmp2, _ = tmp1.Html()
		NewsData.Title = RemoveHtmlTag(tmp2)
		if NewsData.Title == "News Ticker" {
			NewsData.Type = "ticker"
			NewsData.Title = ""
		}
	})

	ReaderHTML.Find(".NewsTableContainer").Each(func(index int, s *goquery.Selection) {

		// checking if its a ticker..
		if NewsData.Type == "ticker" {
			tmp1 = s.Find("p")
			NewsData.Content = tmp1.Text()
			NewsData.ContentHTML, _ = tmp1.Html()
		} else {
			// getting html
			tmp2, _ = s.First().Html()
			// replacing martel letter in articles with real letters
			tmp2 = martelRegex.ReplaceAllString(tmp2, "$1")
			s.ReplaceWithHtml(tmp2)

			// storing html content
			NewsData.ContentHTML = tmp2

			// reading string again after replacing letters
			tmp1, _ := goquery.NewDocumentFromReader(strings.NewReader(tmp2))

			// storing text content
			NewsData.Content = tmp1.Text()
		}
	})

	//
	// Build the data-blob
	return NewsResponse{
		NewsData,
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
