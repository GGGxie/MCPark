package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// parsePaging reads page/page_size from the query string. page is clamped to
// [1, maxPage]; page_size falls back to defaultSize when missing or < 1 and is
// capped at maxSize.
func parsePaging(c *gin.Context, defaultSize, maxSize int) (page, pageSize int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	if page > maxPage {
		page = maxPage
	}
	pageSize, err = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultSize)))
	if err != nil || pageSize < 1 {
		pageSize = defaultSize
	}
	if pageSize > maxSize {
		pageSize = maxSize
	}
	return page, pageSize
}

// bindJSON decodes the request body into v. On failure it writes the error
// response (413 for oversized bodies, 400 otherwise) and returns false.
func bindJSON(c *gin.Context, v any) bool {
	err := c.ShouldBindJSON(v)
	if err == nil {
		return true
	}
	var tooBig *http.MaxBytesError
	if errors.As(err, &tooBig) {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request body too large"})
		return false
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return false
}
