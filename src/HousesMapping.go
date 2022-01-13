package main

import (
	"encoding/json"
	"io/ioutil"
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
	// load content from file into variable file
	file, _ := ioutil.ReadFile("houses_mapping.json")

	// loading json and mapping it into the data variable
	data := HousesMapping{}
	_ = json.Unmarshal([]byte(file), &data)

	// storing data so it's accessible from other places
	TibiadataHousesMapping = data
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
