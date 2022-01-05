package main

import "github.com/gin-gonic/gin"

// TibiaHousesHouseV3 func
func TibiaHousesHouseV3(c *gin.Context) {

	/*
		// getting params from URL
		world := c.Param("world")
		houseid := c.Param("houseid")
	*/

	//
	// The base
	type JSONData struct{}

	/*
		// Adding fix for First letter to be upper and rest lower
		world = TibiadataStringWorldFormatToTitleV3(world)
	*/

	//
	// Build the data-blob
	jsonData := JSONData{}

	// return jsonData
	TibiaDataAPIHandleSuccessResponse(c, "TibiaHousesHouseV3", jsonData)
}
