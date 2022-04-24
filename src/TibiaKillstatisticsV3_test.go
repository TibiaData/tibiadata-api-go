package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestAntica(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/killstatistics/Antica.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	anticaJson, err := TibiaKillstatisticsV3Impl("Antica", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("Antica", anticaJson.KillStatistics.World)
	assert.Equal(1159, len(anticaJson.KillStatistics.Entries))

	elementalForces := anticaJson.KillStatistics.Entries[0]
	assert.Equal("(elemental forces)", elementalForces.Race)
	assert.Equal(6, elementalForces.LastDayKilledPlayers)
	assert.Equal(0, elementalForces.LastDayKilledByPlayers)
	assert.Equal(103, elementalForces.LastWeekKilledPlayers)
	assert.Equal(0, elementalForces.LastWeekKilledByPlayers)

	caveRats := anticaJson.KillStatistics.Entries[386]
	assert.Equal("cave rats", caveRats.Race)
	assert.Equal(2, caveRats.LastDayKilledPlayers)
	assert.Equal(1618, caveRats.LastDayKilledByPlayers)
	assert.Equal(78, caveRats.LastWeekKilledPlayers)
	assert.Equal(12876, caveRats.LastWeekKilledByPlayers)
}

func BenchmarkAntica(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/killstatistics/Antica.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}

	b.ReportAllocs()

	assert := assert.New(b)

	for i := 0; i < b.N; i++ {
		anticaJson, err := TibiaKillstatisticsV3Impl("Antica", string(data))
		if err != nil {
			b.Fatal(err)
		}

		assert.Equal("Antica", anticaJson.KillStatistics.World)
	}
}
