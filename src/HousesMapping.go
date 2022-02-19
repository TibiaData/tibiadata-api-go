package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	TibiadataHousesMapping HousesMapping
)

type AssetsHouse struct {
	HouseID   int    `json:"house_id"`
	Town      string `json:"town"`
	HouseType string `json:"type"`
}
type HousesMapping struct {
	Houses []AssetsHouse `json:"houses"`
}

// TibiaDataHousesMappingInitiator func - used to load data from local JSON file
func TibiaDataHousesMappingInitiator() {

	// Setting up resty client
	client := resty.New()

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   TibiadataUserAgent,
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	TibiadataAssetsURL := "https://assets.tibiadata.com/data.min.json"
	res, err := client.R().Get(TibiadataAssetsURL)

	switch res.StatusCode() {
	case http.StatusOK:
		// adding response into the data field
		data := HousesMapping{}
		err = json.Unmarshal([]byte(res.Body()), &data)

		if err != nil {
			log.Println("[error] TibiaData API failed to parse content from assets.tibiadata.com/data.min.json")
		} else {
			// storing data so it's accessible from other places
			TibiadataHousesMapping = data
		}

	default:
		log.Printf("[error] TibiaData API failed to load houses mapping. %s", err)
	}
}

// TibiaDataHousesMapResolver func - used to return both town and type
func TibiaDataHousesMapResolver(houseid int) (town string, housetype string) {
	for _, value := range TibiadataHousesMapping.Houses {
		if houseid == value.HouseID {
			return value.Town, value.HouseType
		}
	}
	return "", ""
}
