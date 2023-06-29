package main

import (
	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestForumsSectionWorldBoard(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/forums/section/worldboard.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boardsJson, err := TibiaForumSectionImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(90, len(boardsJson.Boards))

	adra := boardsJson.Boards[0]
	assert.Equal(146482, adra.ID)
	assert.Equal("Adra", adra.Name)
	assert.Equal("This board is for general discussions related to the game world Adra.", adra.Description)
	assert.Equal(388, adra.Threads)
	assert.Equal(1158, adra.Posts)
	assert.Equal(39395612, adra.LastPost.ID)
	assert.Equal("Rcdohl", adra.LastPost.CharacterName)
	assert.Equal("2023-06-02T23:19:00Z", adra.LastPost.PostedAt)

	alumbra := boardsJson.Boards[1]
	assert.Equal(147016, alumbra.ID)
	assert.Equal("Alumbra", alumbra.Name)
	assert.Equal("This board is for general discussions related to the game world Alumbra.", alumbra.Description)
	assert.Equal(563, alumbra.Threads)
	assert.Equal(1011, alumbra.Posts)
	assert.Equal(39395777, alumbra.LastPost.ID)
	assert.Equal("Mad Mustazza", alumbra.LastPost.CharacterName)
	assert.Equal("2023-06-04T15:51:13Z", alumbra.LastPost.PostedAt)
}

func TestForumsSectionSupportBoards(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/forums/section/supportboards.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boardsJson, err := TibiaForumSectionImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(3, len(boardsJson.Boards))

	paymentSupport := boardsJson.Boards[0]
	assert.Equal(20, paymentSupport.ID)
	assert.Equal("Payment Support", paymentSupport.Name)
	assert.Equal("Here you can ask questions on orders, payments and other payment issues.", paymentSupport.Description)
	assert.Equal(14468, paymentSupport.Threads)
	assert.Equal(56032, paymentSupport.Posts)
	assert.Equal(39400410, paymentSupport.LastPost.ID)
	assert.Equal("Dorieta", paymentSupport.LastPost.CharacterName)
	assert.Equal("2023-06-29T11:24:33Z", paymentSupport.LastPost.PostedAt)

	technicalSupport := boardsJson.Boards[1]
	assert.Equal(13, technicalSupport.ID)
	assert.Equal("Technical Support", technicalSupport.Name)
	assert.Equal("Here you can ask for help if you have a technical problem that is related to Tibia.", technicalSupport.Description)
	assert.Equal(64722, technicalSupport.Threads)
	assert.Equal(301828, technicalSupport.Posts)
	assert.Equal(39400354, technicalSupport.LastPost.ID)
	assert.Equal("Dio Sorcer", technicalSupport.LastPost.CharacterName)
	assert.Equal("2023-06-29T01:47:59Z", technicalSupport.LastPost.PostedAt)

	help := boardsJson.Boards[2]
	assert.Equal(113174, help.ID)
	assert.Equal("Help", help.Name)
	assert.Equal("Here you can ask other players all kind of questions about Tibia. Note that members of the CipSoft team will usually not reply here.", help.Description)
	assert.Equal(22788, help.Threads)
	assert.Equal(106713, help.Posts)
	assert.Equal(39400248, help.LastPost.ID)
	assert.Equal("Dekraken durinars", help.LastPost.CharacterName)
	assert.Equal("2023-06-28T17:17:45Z", help.LastPost.PostedAt)
}

func TestForumsSectionCommunityBoards(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/forums/section/communityboards.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boardsJson, err := TibiaForumSectionImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(6, len(boardsJson.Boards))

	gameplay := boardsJson.Boards[0]
	assert.Equal(120835, gameplay.ID)
	assert.Equal("Gameplay (English Only)", gameplay.Name)
	assert.Equal("Here is the place to talk about quests, achievements and the general gameplay of Tibia.", gameplay.Description)
	assert.Equal(35364, gameplay.Threads)
	assert.Equal(335328, gameplay.Posts)
	assert.Equal(39400492, gameplay.LastPost.ID)
	assert.Equal("Dreamsphere", gameplay.LastPost.CharacterName)
	assert.Equal("2023-06-29T19:59:41Z", gameplay.LastPost.PostedAt)

	auditorium := boardsJson.Boards[5]
	assert.Equal(89516, auditorium.ID)
	assert.Equal("Auditorium (English Only)", auditorium.Name)
	assert.Equal("Meet Tibia's community managers and give feedback on news articles and Tibia related topics. State your opinion and see what others think and have to say.", auditorium.Description)
}

func TestForumsSectionTradeBoards(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/forums/section/tradeboards.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	boardsJson, err := TibiaForumSectionImpl(string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal(90, len(boardsJson.Boards))

	adra := boardsJson.Boards[0]
	assert.Equal(146485, adra.ID)
	assert.Equal("Adra - Trade", adra.Name)
	assert.Equal("Use this board to make Tibia related trade offers on the game world Adra. Note that all trades must conform to the Tibia Rules.", adra.Description)
	assert.Equal(540, adra.Threads)
	assert.Equal(599, adra.Posts)
	assert.Equal(39400121, adra.LastPost.ID)
	assert.Equal("Arge Reotona", adra.LastPost.CharacterName)
	assert.Equal("2023-06-28T09:15:31Z", adra.LastPost.PostedAt)
}
