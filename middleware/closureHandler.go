package middleware

import (
	"awesomeProject/util"
	"bytes"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClosureHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch errV := err.(type) {
			case util.HttpExceptionWithLog:
				util.LogContextTrace(c, errV.Log, errV.Error)
				util.FailApiResponse(c, errV.StatusCode, errV.Error)
			case util.HttpException:
				util.FailApiResponse(c, errV.StatusCode, errV.Error)
			default:
				util.FailApiResponse(c, http.StatusInternalServerError, nil)
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

func RequestResponseTransformer(ctx *gin.Context) {
	w := &ResponseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = w

	ctx.Next()

	var xmlMap gin.H
	decoder := xml.NewDecoder(bytes.NewReader(w.body.Bytes()))
	if err := decoder.Decode(&xmlMap); err != nil {
		util.LogError("middleware::RequestResponseTransformer - Unhandled error during XML decode", err)
		return
	}
	if j, err := xml.Marshal(xmlMap); err == nil {
		w.body.Reset()
		w.body.Write(j)
	}
}

type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	n, err := r.body.Write(b)
	return n, err
}

func (r *ResponseBodyWriter) GetBody() *bytes.Buffer {
	return r.body
}
