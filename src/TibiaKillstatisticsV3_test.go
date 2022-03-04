package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAntica(t *testing.T) {
	data, err := os.ReadFile("../testdata/killstatistics/Antica.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	anticaJson := TibiaKillstatisticsV3Impl("Antica", string(data))
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
	data, err := os.ReadFile("../testdata/killstatistics/Antica.html")
	if err != nil {
		b.Errorf("File reading error: %s", err)
		return
	}

	b.ReportAllocs()

	assert := assert.New(b)

	for i := 0; i < b.N; i++ {
		anticaJson := TibiaKillstatisticsV3Impl("Antica", string(data))

		assert.Equal("Antica", anticaJson.KillStatistics.World)
	}
}
