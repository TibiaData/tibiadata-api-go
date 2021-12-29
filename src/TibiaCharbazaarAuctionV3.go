package main

import (
	"encoding/json"
	"fmt"
	//"log"
	//"regexp"
	//"strings"

	//"github.com/PuerkitoBio/goquery"
)

// TibiaCharbazaarAuctionV3 func
func TibiaCharbazaarAuctionV3(id string) string {

	// Child of Auction
	type Outfit struct {
		Name	string	`json:"name"`
		Addon1	bool	`json:"addon1"`
		Addon2	bool	`json:"addon2"`
	}

	// Child of JSONData
	type Auction struct {
		Id		string	`json:"id"`
		CharacterName	string	`json:"characterName"`
		Level			int		`json:"level"`
		Vocation		string	`json:"vocation"`
	}

	// The base includes two levels: Auction and Information
	type JSONData struct {
		Auction		Auction		`json:"auction"`
		Information	Information	`json:"information"`
	}

	jsonData := JSONData{
		Auction: Auction{
			Id: id,
			CharacterName: "Pheizx",
			Level: 170,
			Vocation: "Sorcerer",
		},
		Information: Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	js, _ := json.Marshal(jsonData)
	if TibiadataDebug {
		fmt.Printf("%s\n", js)
	}

	return string(js)
}