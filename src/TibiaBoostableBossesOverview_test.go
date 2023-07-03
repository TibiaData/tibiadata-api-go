package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestBoostableBossesOverview(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boostableBossesJson, _ := TibiaBoostableBossesOverviewImpl(string(data))
	assert := assert.New(t)
	boosted := boostableBossesJson.BoostableBosses.Boosted
	bosses := boostableBossesJson.BoostableBosses.BoostableBosses

	assert.Equal(91, len(bosses))
	assert.Equal("Sharpclaw", boosted.Name)
	assert.Equal(
		"https://static.tibia.com/images/global/header/monsters/sharpclaw.gif",
		boosted.ImageURL,
	)

	for _, tc := range []struct {
		idx      int
		name     string
		featured bool
		imageURL string
	}{
		{
			idx:      19,
			name:     "Gnomevil",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/gnomehorticulist.gif",
		},
		{
			idx:      24,
			name:     "Goshnar's Malice",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/goshnarsmalice.gif",
		},
		{
			idx:      52,
			name:     "Sharpclaw",
			featured: true,
			imageURL: "https://static.tibia.com/images/library/sharpclaw.gif",
		},
		{
			idx:      75,
			name:     "The Pale Worm",
			featured: false,
			imageURL: "https://static.tibia.com/images/library/paleworm.gif",
		},
	} {
		boss := bosses[tc.idx]
		assert.Equal(
			tc.name, boss.Name,
			"Wrong name\nidx: %d (%s)\nwant: %s\ngot: %s",
			tc.idx, tc.name, tc.name, boss.Name,
		)
		assert.Equal(
			tc.featured, boss.Featured,
			"Wrong featured status\nidx: %d (%s)\nwant: %v\ngot: %v",
			tc.idx, tc.name, tc.featured, boss.Featured,
		)
		assert.Equal(
			tc.imageURL, boss.ImageURL,
			"Wrong image URL\nidx: %d (%s)\nwant: %s\ngot: %s",
			tc.idx, tc.name, tc.imageURL, boss.ImageURL,
		)
	}
}

var bossSink BoostableBossesOverviewResponse

func BenchmarkTibiaBoostableBossesOverviewImpl(b *testing.B) {
	file, err := static.TestFiles.Open("testdata/boostablebosses/boostablebosses.html")
	if err != nil {
		b.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	rawData, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("File reading error: %s", err)
	}
	data := string(rawData)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bossSink, _ = TibiaBoostableBossesOverviewImpl(data)
	}
}
