package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Child of KillStatistics
type Entry struct {
	Race                    string `json:"race"`                     // The name of the creature/race.
	LastDayKilledPlayers    int    `json:"last_day_players_killed"`  // Number of players killed by this race in the last day.
	LastDayKilledByPlayers  int    `json:"last_day_killed"`          // Number of creatures of this race killed in the last day.
	LastWeekKilledPlayers   int    `json:"last_week_players_killed"` // Number of players killed by this race in the last week.
	LastWeekKilledByPlayers int    `json:"last_week_killed"`         // Number of creatures of this race killed in the last week.
}

// Child of KillStatistics
type Total struct {
	LastDayKilledPlayers    int `json:"last_day_players_killed"`  // Total number of players killed in total in the last day.
	LastDayKilledByPlayers  int `json:"last_day_killed"`          // Total number of creatures in total killed in the last day.
	LastWeekKilledPlayers   int `json:"last_week_players_killed"` // Total number of players killed in total in the last week.
	LastWeekKilledByPlayers int `json:"last_week_killed"`         // Total number of creatures in total killed in the last week.
}

// Child of JSONData
type KillStatistics struct {
	World   string  `json:"world"`   // The world the statistics belong to.
	Entries []Entry `json:"entries"` // List of killstatistic.
	Total   Total   `json:"total"`   // List of total kills.
}

// The base includes two levels: KillStatistics and Information
type KillStatisticsResponse struct {
	KillStatistics KillStatistics `json:"killstatistics"`
	Information    Information    `json:"information"`
}

func TibiaKillstatisticsImpl(world string, BoxContentHTML string) (KillStatisticsResponse, error) {
	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		return KillStatisticsResponse{}, fmt.Errorf("[error] TibiaKillstatisticsImpl failed at goquery.NewDocumentFromReader, err: %s", err)
	}

	// Creating empty KillStatisticsData var
	var (
		KillStatisticsData                                                                                               []Entry
		TotalLastDayKilledPlayers, TotalLastDayKilledByPlayers, TotalLastWeekKilledPlayers, TotalLastWeekKilledByPlayers int
	)

	// Running query over each div
	ReaderHTML.Find("#KillStatisticsTable .TableContent tr.Odd,tr.Even").Each(func(index int, s *goquery.Selection) {
		DataColumns := s.Find("td").Nodes

		KillStatisticsLastDayKilledPlayers := TibiaDataStringToInteger(DataColumns[1].FirstChild.Data)
		TotalLastDayKilledPlayers += KillStatisticsLastDayKilledPlayers
		KillStatisticsLastDayKilledByPlayers := TibiaDataStringToInteger(DataColumns[2].FirstChild.Data)
		TotalLastDayKilledByPlayers += KillStatisticsLastDayKilledByPlayers
		KillStatisticsLastWeekKilledPlayers := TibiaDataStringToInteger(DataColumns[3].FirstChild.Data)
		TotalLastWeekKilledPlayers += KillStatisticsLastWeekKilledPlayers
		KillStatisticsLastWeekKilledByPlayers := TibiaDataStringToInteger(DataColumns[4].FirstChild.Data)
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
			APIDetails: TibiaDataAPIDetails,
			Timestamp:  TibiaDataDatetime(""),
			Link:       "https://www.tibia.com/community/?subtopic=killstatistics&world=" + TibiaDataQueryEscapeString(world),
			Status: Status{
				HTTPCode: http.StatusOK,
			},
		},
	}, nil
}
