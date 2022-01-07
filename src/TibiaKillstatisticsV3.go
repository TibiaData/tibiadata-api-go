package main

import (
	"log"
	"regexp"
	"strings"
	"tibiadata-api-go/src/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// TibiaKillstatisticsV3 func
func TibiaKillstatisticsV3(c *gin.Context) {
	// getting params from URL
	world := c.Param("world")

	// The base includes two levels: KillStatistics and Information
	type JSONData struct {
		KillStatistics structs.KillStatistics `json:"killstatistics"`
		Information    structs.Information    `json:"information"`
	}

	// Adding fix for First letter to be upper and rest lower
	world = TibiadataStringWorldFormatToTitleV3(world)

	// Getting data with TibiadataHTMLDataCollectorV3
	TibiadataRequest.URL = "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiadataQueryEscapeStringV3(world)
	BoxContentHTML := TibiadataHTMLDataCollectorV3(TibiadataRequest)

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty KillStatisticsData var
	var KillStatisticsData []structs.Entry
	var TotalLastDayKilledPlayers, TotalLastDayKilledByPlayers, TotalLastWeekKilledPlayers, TotalLastWeekKilledByPlayers int

	// Running query over each div
	ReaderHTML.Find("#KillStatisticsTable .TableContent tr").Each(func(index int, s *goquery.Selection) {
		// Storing HTML into CreatureDivHTML
		KillStatisticsDivHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Regex when highscore has 5 columns
		regex1 := regexp.MustCompile(`<td>(.*)<\/td><td.*>([0-9]+).*<\/td><td.*>([0-9]+).*<\/td><td.*>([0-9]+).*<\/td><td.*>([0-9]+).*<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(KillStatisticsDivHTML, -1)

		if len(subma1) > 0 {
			if strings.TrimSpace(subma1[0][1]) == "Total" {
				// we don't want to include the Total row
			} else {
				// Store the values..
				KillStatisticsLastDayKilledPlayers := TibiadataStringToIntegerV3(subma1[0][2])
				TotalLastDayKilledPlayers += KillStatisticsLastDayKilledPlayers
				KillStatisticsLastDayKilledByPlayers := TibiadataStringToIntegerV3(subma1[0][3])
				TotalLastDayKilledByPlayers += KillStatisticsLastDayKilledByPlayers
				KillStatisticsLastWeekKilledPlayers := TibiadataStringToIntegerV3(subma1[0][4])
				TotalLastWeekKilledPlayers += KillStatisticsLastWeekKilledPlayers
				KillStatisticsLastWeekKilledByPlayers := TibiadataStringToIntegerV3(subma1[0][5])
				TotalLastWeekKilledByPlayers += KillStatisticsLastWeekKilledByPlayers

				// Append new Entry item to KillStatisticsData
				KillStatisticsData = append(KillStatisticsData, structs.Entry{
					Race:                    TibiaDataSanitizeEscapedString(subma1[0][1]),
					LastDayKilledPlayers:    KillStatisticsLastDayKilledPlayers,
					LastDayKilledByPlayers:  KillStatisticsLastDayKilledByPlayers,
					LastWeekKilledPlayers:   KillStatisticsLastWeekKilledPlayers,
					LastWeekKilledByPlayers: KillStatisticsLastWeekKilledByPlayers,
				})
			}
		}
	})

	//
	// Build the data-blob
	jsonData := JSONData{
		structs.KillStatistics{
			World:   world,
			Entries: KillStatisticsData,
			Total: structs.Total{
				LastDayKilledPlayers:    TotalLastDayKilledPlayers,
				LastDayKilledByPlayers:  TotalLastDayKilledByPlayers,
				LastWeekKilledPlayers:   TotalLastWeekKilledPlayers,
				LastWeekKilledByPlayers: TotalLastWeekKilledByPlayers,
			},
		},
		structs.Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaKillstatisticsV3", jsonData)

}
