package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	const expectedUserAgent = "TibiaData-API/v4 (release/unknown; build/manual; commit/-; edition/open-source; unittest.example.com)"

	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)
	assert.Equal(t, expectedUserAgent, TibiaDataUserAgent)
}

func TestTibiaDataInitializer(t *testing.T) {
	assert := assert.New(t)

	// Call the function to be tested
	TibiaDataInitializer()

	// Check that the variables have been set correctly
	assert.Equal("open-source", TibiaDataBuildEdition)
	assert.Equal("unittest.example.com", TibiaDataHost)
	assert.Equal("", TibiaDataProxyDomain)
}
