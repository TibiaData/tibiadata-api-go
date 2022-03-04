package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of JSONData
type NewsItem struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	News     string `json:"news"`
	Category string `json:"category"`
	Type     string `json:"type"`
	TibiaURL string `json:"url"`
	ApiURL   string `json:"url_api,omitempty"`
}

//
// The base
type NewsListResponse struct {
	News        []NewsItem  `json:"news"`
	Information Information `json:"information"`
}

func TibiaNewslistV3Impl(days int, BoxContentHTML string) NewsListResponse {
	// Declaring vars for later use..
	var NewsListData []NewsItem

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	ReaderHTML.Find(".Odd,.Even").Each(func(index int, s *goquery.Selection) {
		var OneNews NewsItem

		// getting category by image src
		CategoryImg, _ := s.Find("img").Attr("src")
		OneNews.Category = TibiaDataGetNewsCategory(CategoryImg)

		// getting type from headline
		NewsType := s.Nodes[0].FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild.Data
		OneNews.Type = TibiaDataGetNewsType(TibiaDataSanitizeNbspSpaceString(NewsType))

		// getting date from headline
		OneNews.Date = TibiaDataDateV3(s.Nodes[0].FirstChild.NextSibling.FirstChild.Data)
		OneNews.News = s.Find("a").Text()

		// getting remaining things as URLs
		NewsURL, _ := s.Find("a").Attr("href")
		p, _ := url.Parse(NewsURL)
		NewsID := p.Query().Get("id")
		NewsSplit := strings.Split(NewsURL, NewsID)
		OneNews.ID = TibiaDataStringToIntegerV3(NewsID)
		OneNews.TibiaURL = NewsSplit[0] + NewsID

		if TibiaDataHost != "" {
			OneNews.ApiURL = "https://" + TibiaDataHost + "/v3/news/id/" + NewsID
		}

		// add to NewsListData for response
		NewsListData = append(NewsListData, OneNews)
	})

	//
	// Build the data-blob
	return NewsListResponse{
		NewsListData,
		Information{
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
		},
	}
}
