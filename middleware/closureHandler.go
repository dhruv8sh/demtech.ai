package middleware

import (
	"awesomeProject/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClosureHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch errV := err.(type) {
				case util.HttpExceptionWithLog:

					util.LogContextTrace(c, errV.Log, errV.Error)

					util.FailApiResponse(c, errV.StatusCode,
						util.FailedApiRespError{
							FlagException: true,
							Errors:        errV.Error,
						})
				case util.HttpException:
					util.FailApiResponse(c, errV.StatusCode,
						util.FailedApiRespError{
							FlagException: true,
							Errors:        errV.Error,
						})
				default:
					util.FailApiResponse(c, http.StatusInternalServerError,
						util.FailedApiRespError{
							FlagPanic: true,
							Errors:    err,
						})
				}
			}

			switch c.Writer.(type) {
			case *ResponseBodyWriter:
				w := c.Writer.(*ResponseBodyWriter)
				if w.body.Len() > 0 {
					_, err := w.ResponseWriter.Write(w.body.Bytes())
					if err != nil {
					}
					w.body.Reset()
				}
			}
		}()

		c.Next()
	}
}
