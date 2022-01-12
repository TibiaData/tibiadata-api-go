package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Child of KillStatistics
type Entry struct {
	Race                    string `json:"race"`
	LastDayKilledPlayers    int    `json:"last_day_players_killed"`
	LastDayKilledByPlayers  int    `json:"last_day_killed"`
	LastWeekKilledPlayers   int    `json:"last_week_players_killed"`
	LastWeekKilledByPlayers int    `json:"last_week_killed"`
}

// Child of KillStatistics
type Total struct {
	LastDayKilledPlayers    int `json:"last_day_players_killed"`
	LastDayKilledByPlayers  int `json:"last_day_killed"`
	LastWeekKilledPlayers   int `json:"last_week_players_killed"`
	LastWeekKilledByPlayers int `json:"last_week_killed"`
}

// Child of JSONData
type KillStatistics struct {
	World   string  `json:"world"`
	Entries []Entry `json:"entries"`
	Total   Total   `json:"total"`
}

//
// The base includes two levels: KillStatistics and Information
type KillStatisticsResponse struct {
	KillStatistics KillStatistics `json:"killstatistics"`
	Information    Information    `json:"information"`
}

// TibiaKillstatisticsV3 func
func TibiaKillstatisticsV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiadataQueryEscapeStringV3(world)
	BoxContentHTML, err := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// return error (e.g. for maintenance mode)
	if err != nil {
		TibiaDataAPIHandleOtherResponse(c, http.StatusBadGateway, "TibiaKillstatisticsV3", gin.H{"error": err.Error()})
		return
	}

	jsonData := TibiaKillstatisticsV3Impl(world, BoxContentHTML)

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaKillstatisticsV3", jsonData)
}

func TibiaKillstatisticsV3Impl(world string, BoxContentHTML string) KillStatisticsResponse {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty KillStatisticsData var
	var (
		KillStatisticsData                                                                                               []Entry
		TotalLastDayKilledPlayers, TotalLastDayKilledByPlayers, TotalLastWeekKilledPlayers, TotalLastWeekKilledByPlayers int
	)

	// Running query over each div
	ReaderHTML.Find("#KillStatisticsTable .TableContent tr.Odd,tr.Even").Each(func(index int, s *goquery.Selection) {
		DataColumns := s.Find("td").Nodes

		KillStatisticsLastDayKilledPlayers := TibiadataStringToIntegerV3(DataColumns[1].FirstChild.Data)
		TotalLastDayKilledPlayers += KillStatisticsLastDayKilledPlayers
		KillStatisticsLastDayKilledByPlayers := TibiadataStringToIntegerV3(DataColumns[2].FirstChild.Data)
		TotalLastDayKilledByPlayers += KillStatisticsLastDayKilledByPlayers
		KillStatisticsLastWeekKilledPlayers := TibiadataStringToIntegerV3(DataColumns[3].FirstChild.Data)
		TotalLastWeekKilledPlayers += KillStatisticsLastWeekKilledPlayers
		KillStatisticsLastWeekKilledByPlayers := TibiadataStringToIntegerV3(DataColumns[4].FirstChild.Data)
		TotalLastWeekKilledByPlayers += KillStatisticsLastWeekKilledByPlayers

		// Append new Entry item to KillStatisticsData
		KillStatisticsData = append(KillStatisticsData, Entry{
			Race:                    TibiaDataSanitizeEscapedString(DataColumns[0].FirstChild.Data),
			LastDayKilledPlayers:    KillStatisticsLastDayKilledPlayers,
			LastDayKilledByPlayers:  KillStatisticsLastDayKilledByPlayers,
			LastWeekKilledPlayers:   KillStatisticsLastWeekKilledPlayers,
			LastWeekKilledByPlayers: KillStatisticsLastWeekKilledByPlayers,
		})
	})

	//
	// Build the data-blob
	return KillStatisticsResponse{
		KillStatistics{
			World:   world,
			Entries: KillStatisticsData,
			Total: Total{
				LastDayKilledPlayers:    TotalLastDayKilledPlayers,
				LastDayKilledByPlayers:  TotalLastDayKilledByPlayers,
				LastWeekKilledPlayers:   TotalLastWeekKilledPlayers,
				LastWeekKilledByPlayers: TotalLastWeekKilledByPlayers,
			},
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}
}
