package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func TibiaKillstatisticsV3Impl(world string, BoxContentHTML string) (*KillStatisticsResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return nil, fmt.Errorf("[error] TibiaKillstatisticsV3Impl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty KillStatisticsData var
	var (
		KillStatisticsData                                                                                               []Entry
		TotalLastDayKilledPlayers, TotalLastDayKilledByPlayers, TotalLastWeekKilledPlayers, TotalLastWeekKilledByPlayers int
	)

	// Running query over each div
	ReaderHTML.Find("#KillStatisticsTable .TableContent tr.Odd,tr.Even").Each(func(index int, s *goquery.Selection) {
		DataColumns := s.Find("td").Nodes

		KillStatisticsLastDayKilledPlayers := TibiaDataStringToIntegerV3(DataColumns[1].FirstChild.Data)
		TotalLastDayKilledPlayers += KillStatisticsLastDayKilledPlayers
		KillStatisticsLastDayKilledByPlayers := TibiaDataStringToIntegerV3(DataColumns[2].FirstChild.Data)
		TotalLastDayKilledByPlayers += KillStatisticsLastDayKilledByPlayers
		KillStatisticsLastWeekKilledPlayers := TibiaDataStringToIntegerV3(DataColumns[3].FirstChild.Data)
		TotalLastWeekKilledPlayers += KillStatisticsLastWeekKilledPlayers
		KillStatisticsLastWeekKilledByPlayers := TibiaDataStringToIntegerV3(DataColumns[4].FirstChild.Data)
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
	return &KillStatisticsResponse{
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
			APIVersion: TibiaDataAPIversion,
			Timestamp:  TibiaDataDatetimeV3(""),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
