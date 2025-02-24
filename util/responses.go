package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CustomHttpError interface {
	GetHttpStatus() int
	Error() string
}

// FailApiResponse writes a response to the body with given HttpStatus
// Optionally logs to the app.log file
func FailApiResponse(c *gin.Context, httpStatus int, err any) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	switch errV := err.(type) {
	case error:
		LogContextTrace(c, "Error", errV.Error())
		c.Abort()
		c.XML(http.StatusInternalServerError, gin.H{
			"status": 0,
			"errors": errV.Error(),
		})
		return
	case string: // Direct error message is given; this can be returned
		c.Abort()
		c.XML(httpStatus, gin.H{
			"status": 0,
			"errors": []string{errV},
		})
		return
	case []string: // Direct error message is given; this can be returned
		c.Abort()
		c.XML(httpStatus, gin.H{
			"status": 0,
			"errors": errV,
		})
		return
	case CustomHttpError:
		c.Abort()
		c.XML(errV.GetHttpStatus(), gin.H{
			"status": 0,
			"errors": errV.Error(),
		})
		return
	}

	switch httpStatus {
	case http.StatusUnprocessableEntity:
		c.Abort()
		c.XML(httpStatus, gin.H{
			"status": 0,
			"errors": "Invalid request",
		})
	default:
		c.AbortWithStatus(httpStatus)
		c.XML(httpStatus, gin.H{
			"status": 0,
			"errors": "Unknown error",
		})
		return
	}
}

// SuccessApiResponse writes a response to the body with Status 200 (OK)
func SuccessApiResponse(c *gin.Context, message string, data any) {
	c.XML(http.StatusOK,
		gin.H{
			"status":  1,
			"message": message,
			"data":    data,
		})
}
