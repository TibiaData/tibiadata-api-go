package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDebugHandler(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	debugHandler(c)

	assert.Equal(http.StatusOK, w.Code)
}
