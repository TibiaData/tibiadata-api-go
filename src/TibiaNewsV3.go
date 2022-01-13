package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
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

// TibiaNewsV3 func
func TibiaNewsV3(c *gin.Context) {
	// getting params from URL
	NewsID := TibiadataStringToIntegerV3(c.Param("news_id"))

	// checking the NewsID provided
	if NewsID <= 0 {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadRequest, "TibiaNewsV3", gin.H{"error": "no valid news_id provided"})
		return
	}

	TibiadataRequest.URL = "https://www.tibia.com/news/?subtopic=newsarchive&id=" + strconv.Itoa(NewsID)

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaNewslistV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaNewsV3Impl(NewsID, TibiadataRequest.URL, BoxContentHTML)

	TibiaDataAPIHandleSuccessResponse(c, "TibiaNewsV3", jsonData)
}

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
		NewsData.Category = TibiadataGetNewsCategory(CategoryImg)

		// getting date from headline
		tmp1 = s.Find(".NewsHeadlineDate")
		tmp2, _ = tmp1.Html()
		NewsData.Date = TibiadataDateV3(strings.ReplaceAll(tmp2, " - ", ""))

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
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
