package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	const (
		expectedUserAgent     = "TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source)"
		expectedUserAgentHost = "TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source; +https://unittest.example.com)"
	)

	TibiaDataHost = ""
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	assert.Equal(t, expectedUserAgent, TibiaDataUserAgent)

	TibiaDataHost = "unittest.example.com"
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	assert.Equal(t, expectedUserAgentHost, TibiaDataUserAgent)
}

func TestTibiaDataInitializer(t *testing.T) {
	assert := assert.New(t)

	TibiaDataHost = "unittest.example.com"

	// Call the function to be tested
	TibiaDataInitializer()

	// Check that the variables have been set correctly
	assert.Equal("open-source", TibiaDataBuildEdition)
	assert.Equal("https", TibiaDataProtocol)
	assert.Equal("unittest.example.com", TibiaDataHost)
}
