package main

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	requestBodyReadError    = "Failed to read request body"
)

// GetRouter returns the Gin router.
func GetRouter(port string, dbConnection *DBConnection) *gin.Engine {
	h := &TokenHandler{
		dbConn: dbConnection,
	}

	r := gin.Default()
	r.POST("/use_token", h.postUseToken)
	r.GET("/test_get_token", h.getTestToken)
	return r
}

// TokenHandler implements the handler functions for the API endpoints.
// It also holds the database connection that's used by the handler functions.
type TokenHandler struct {
	dbConn *DBConnection
}

func (h *TokenHandler) postUseToken(c *gin.Context) {
	body := c.Request.Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		c.String(http.StatusBadRequest, requestBodyReadError)
		return
	}

	valid, err := h.dbConn.checkAndRemoveToken(string(data))
	
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(strconv.FormatBool(valid)))
}

func (h *TokenHandler) getTestToken(c *gin.Context) {
	token, err := h.dbConn.createToken()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(token))
}
