package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/foods", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
	return r
}

func TestGetFoods(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foods", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
