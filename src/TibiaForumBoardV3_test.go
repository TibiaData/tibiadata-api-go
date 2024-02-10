package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBoardsBoard(t *testing.T) {
	data, err := os.ReadFile("../testdata/boards/board.html")
	if err != nil {
		t.Errorf("File reading error: %s", err)
		return
	}

	boardsJson := TibiaForumBoardV3Impl("144647", string(data), 1, lastDay)
	assert := assert.New(t)

	assert.Equal(144647, boardsJson.Board.ID)
	assert.Equal("World Boards", boardsJson.Board.SectionName)
	assert.Equal("Vunira", boardsJson.Board.BoardName)
	assert.Equal("lastDay", boardsJson.Board.ThreadsAge)

	assert.Equal(32, boardsJson.Board.BoardsBoardPagination.TotalResults)
	assert.Equal(2, boardsJson.Board.BoardsBoardPagination.TotalPages)
	assert.Equal(1, boardsJson.Board.BoardsBoardPagination.CurrentPage)

	assert.Equal(30, len(boardsJson.Board.ThreadsList))

	multilineTitleThread := boardsJson.Board.ThreadsList[0]
	assert.Equal(4950043, multilineTitleThread.ID)
	assert.Equal("TibiaMaps.io FERU HAT LOTTERY! 23/50 TICKETS SOLD", multilineTitleThread.Title)

	hotThread := boardsJson.Board.ThreadsList[0]
	assert.Equal(true, hotThread.IsHot)
	assert.Equal(false, hotThread.IsNew)
	assert.Equal(false, hotThread.IsClosed)
	assert.Equal(false, hotThread.IsSticky)

	iconExclamationThread := boardsJson.Board.ThreadsList[0]
	assert.Equal("Exclamation", iconExclamationThread.Icon)

	shortTitleThread := boardsJson.Board.ThreadsList[1]
	assert.Equal(4947560, shortTitleThread.ID)
	assert.Equal("Access quests", shortTitleThread.Title)

	newThread := boardsJson.Board.ThreadsList[1]
	assert.Equal(false, newThread.IsHot)
	assert.Equal(true, newThread.IsNew)
	assert.Equal(false, newThread.IsClosed)
	assert.Equal(false, newThread.IsSticky)

	lightbulbThread := boardsJson.Board.ThreadsList[1]
	assert.Equal("Lightbulb", lightbulbThread.Icon)

	closedThread := boardsJson.Board.ThreadsList[2]
	assert.Equal("Heart of Destruction and Thais Bosses", closedThread.Title)
	assert.Equal(false, closedThread.IsHot)
	assert.Equal(false, closedThread.IsNew)
	assert.Equal(true, closedThread.IsClosed)
	assert.Equal(false, closedThread.IsSticky)

	zeroPagesThread := boardsJson.Board.ThreadsList[2]
	assert.Equal(0, zeroPagesThread.Pages)

	stickyThread := boardsJson.Board.ThreadsList[3]
	assert.Equal("Friendly Neighborhood Imbuement Service", stickyThread.Title)
	assert.Equal(false, stickyThread.IsHot)
	assert.Equal(false, stickyThread.IsNew)
	assert.Equal(false, stickyThread.IsClosed)
	assert.Equal(true, stickyThread.IsSticky)

	fourPagesThread := boardsJson.Board.ThreadsList[3]
	assert.Equal(4, fourPagesThread.Pages)

	hotClosedStickyThread := boardsJson.Board.ThreadsList[4]
	assert.Equal("TIBIA MOTIVATIONAL VIDEO", hotClosedStickyThread.Title)
	assert.Equal(true, hotClosedStickyThread.IsHot)
	assert.Equal(false, hotClosedStickyThread.IsNew)
	assert.Equal(true, hotClosedStickyThread.IsClosed)
	assert.Equal(true, hotClosedStickyThread.IsSticky)

	noIconThread := boardsJson.Board.ThreadsList[4]
	assert.Equal("", noIconThread.Icon)

	lastPageThread := boardsJson.Board.ThreadsList[11]
	assert.Equal("Oriental ships sighted!", lastPageThread.Title)
	assert.Equal(19, lastPageThread.Pages)

	thousandViewsThread := boardsJson.Board.ThreadsList[11]
	assert.Equal(361, thousandViewsThread.Replies)
	assert.Equal(26732, thousandViewsThread.Views)

	characterWithSpacesThread := boardsJson.Board.ThreadsList[4]
	assert.Equal("Lokiec Amerykanska Ekspedycja", characterWithSpacesThread.CharacterName)
}
