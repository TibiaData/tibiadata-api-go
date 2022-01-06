package main

import (
	"github.com/gin-gonic/gin"
)

// TibiaNewsV3 func
func TibiaNewsV3(c *gin.Context) {

	//
	// The base
	type JSONData struct{}

	//
	// Build the data-blob
	jsonData := JSONData{}

	TibiaDataAPIHandleSuccessResponse(c, "TibiaNewsV3", jsonData)
}
