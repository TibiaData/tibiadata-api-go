package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreadAgeLastYearString(t *testing.T) {
	assert := assert.New(t)

	threadsAge := lastYear
	stringValue, err := threadsAge.String()

	assert.Nil(err)
	assert.Equal("lastYear", stringValue)
	assert.Equal(ThreadsAge(365), threadsAge)
}

func TestThreadsAgeLastDayFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(lastDay, ThreadsAgeFromString("lastDay"))
}

func TestThreadsAgeLast2DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last2Days, ThreadsAgeFromString("last2Days"))
}

func TestThreadsAgeLast5DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last5Days, ThreadsAgeFromString("last5Days"))
}

func TestThreadsAgeLast10DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last10Days, ThreadsAgeFromString("last10Days"))
}

func TestThreadsAgeLast20DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last20Days, ThreadsAgeFromString("last20Days"))
}

func TestThreadsAgeLast30DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last30Days, ThreadsAgeFromString("last30Days"))
}

func TestThreadsAgeLast45DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last45Days, ThreadsAgeFromString("last45Days"))
}

func TestThreadsAgeLast60DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last60Days, ThreadsAgeFromString("last60Days"))
}

func TestThreadsAgeLast75DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last75Days, ThreadsAgeFromString("last75Days"))
}

func TestThreadsAgeLast100DaysFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last100Days, ThreadsAgeFromString("last100Days"))
}

func TestThreadsAgeLastYearFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(lastYear, ThreadsAgeFromString("lastYear"))
}

func TestThreadsAgeAllFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(all, ThreadsAgeFromString("all"))
}

func TestThreadsAgeFromString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(last30Days, ThreadsAgeFromString("last30Days"))
}
