package main

import "github.com/gin-gonic/gin"

// TibiaNewslistV3 func
func TibiaNewslistV3(c *gin.Context) {

	//
	// The base
	type JSONData struct{}

	//
	// Build the data-blob
	jsonData := JSONData{}

	TibiaDataAPIHandleSuccessResponse(c, "TibiaNewslistV3", jsonData)
}
