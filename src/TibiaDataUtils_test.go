package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTibiaCETDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T08:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 CET"))
}

func TestTibiaCESTDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T07:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 CEST"))
}

func TestTibiaUTCDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T09:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 UTC"))
}

func TestTibiaPSTDateFormat(t *testing.T) {
	assert.Equal(t, "2021-12-24T17:52:16Z", TibiadataDatetimeV3("Dec 24 2021, 09:52:16 PST"))
}
