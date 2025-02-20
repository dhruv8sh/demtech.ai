package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FailApiResponse(c *gin.Context, httpStatus int, err interface{}) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	//var displayErrors interface{}
	switch errV := err.(type) {
	case error: // Log error and continue without doing anything...
		LogContextTrace(c, "Error", errV.Error())
	case string: // Direct error message is given; this can be returned
		c.Abort()
		c.PureJSON(httpStatus, map[string]interface{}{
			"status": 0,
			"errors": errV,
		})
		return
	}

	switch httpStatus {
	case http.StatusUnprocessableEntity:
		c.Abort()
		c.PureJSON(httpStatus, map[string]interface{}{
			"status": 0,
			"errors": "Invalid request",
		})
	default:
		c.AbortWithStatus(httpStatus)
	}
}

func SuccessApiResponse(c *gin.Context, message string, data interface{}) {
	// PureJSON, unlike JSON, does not replace special html characters with their unicode entities.
	c.PureJSON(http.StatusOK,
		map[string]interface{}{
			"status":  1,
			"message": message,
			"data":    data,
			"header":  getDataForHeader(c),
		})
}

func getDataForHeader(_ *gin.Context) map[string]interface{} {
	// Populate header data here
	return nil
}
