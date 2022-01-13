package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
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

// TibiaNewslistV3 func
func TibiaNewslistV3(c *gin.Context) {
	// getting params from URL
	days := TibiadataStringToIntegerV3(c.Param("days"))
	if days == 0 {
		days = 90 // default for recent posts
	}

	// generating dates to pass to FormData
	DateBegin := time.Now().AddDate(0, 0, -days)
	DateEnd := time.Now()

	TibiadataRequest.Method = http.MethodPost
	TibiadataRequest.URL = "https://www.tibia.com/news/?subtopic=newsarchive"
	TibiadataRequest.FormData = map[string]string{
		"filter_begin_day":   strconv.Itoa(DateBegin.UTC().Day()),        // period
		"filter_begin_month": strconv.Itoa(int(DateBegin.UTC().Month())), // period
		"filter_begin_year":  strconv.Itoa(DateBegin.UTC().Year()),       // period
		"filter_end_day":     strconv.Itoa(DateEnd.UTC().Day()),          // period
		"filter_end_month":   strconv.Itoa(int(DateEnd.UTC().Month())),   // period
		"filter_end_year":    strconv.Itoa(DateEnd.UTC().Year()),         // period
		"filter_cipsoft":     "cipsoft",                                  // category
		"filter_community":   "community",                                // category
		"filter_development": "development",                              // category
		"filter_support":     "support",                                  // category
		"filter_technical":   "technical",                                // category
	}

	// getting type of news list
	switch tmp := strings.Split(c.Request.URL.Path, "/"); tmp[3] {
	case "newsticker":
		TibiadataRequest.FormData["filter_ticker"] = "ticker"
	case "latest":
		TibiadataRequest.FormData["filter_article"] = "article"
		TibiadataRequest.FormData["filter_news"] = "news"
	case "archive":
		TibiadataRequest.FormData["filter_ticker"] = "ticker"
		TibiadataRequest.FormData["filter_article"] = "article"
		TibiadataRequest.FormData["filter_news"] = "news"
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaNewslistV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaNewslistV3Impl(days, BoxContentHTML)

	TibiaDataAPIHandleSuccessResponse(c, "TibiaNewslistV3", jsonData)
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
		OneNews.Category = TibiadataGetNewsCategory(CategoryImg)

		// getting type from headline
		NewsType := s.Nodes[0].FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild.Data
		OneNews.Type = TibiadataGetNewsType(TibiaDataSanitizeNbspSpaceString(NewsType))

		// getting date from headline
		OneNews.Date = TibiadataDateV3(s.Nodes[0].FirstChild.NextSibling.FirstChild.Data)
		OneNews.News = s.Find("a").Text()

		// getting remaining things as URLs
		NewsURL, _ := s.Find("a").Attr("href")
		p, _ := url.Parse(NewsURL)
		NewsID := p.Query().Get("id")
		NewsSplit := strings.Split(NewsURL, NewsID)
		OneNews.ID = TibiadataStringToIntegerV3(NewsID)
		OneNews.TibiaURL = NewsSplit[0] + NewsID

		if TibiadataHost != "" {
			OneNews.ApiURL = "https://" + TibiadataHost + "/v3/news/id/" + NewsID
		}

		// add to NewsListData for response
		NewsListData = append(NewsListData, OneNews)
	})

	//
	// Build the data-blob
	return NewsListResponse{
		NewsListData,
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
