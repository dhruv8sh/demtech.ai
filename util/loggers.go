package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
)

var logger *log.Logger

const logPath = "~/GolandProjects/awesomeProject/app.log"

func init() {
	logFile, err := os.OpenFile(
		filepath.Join(logPath),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)
	if err != nil {
		panic("Error while initializing logging...")
	}
	logger = log.New(logFile, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

func LogContextTrace(c *gin.Context, msg string, v ...any) {
	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	headers := strings.Split(string(httpRequest), "\r\n")

	headers = append(headers, "ClientIP: "+c.ClientIP())
	headersToStr := strings.Join(headers, "\r\n")

	LogCritical(fmt.Sprintf("\n%s\n%s\n", headersToStr, msg), v...)
}

func LogInfo(msg string, v ...any) {
	logger.Println("[INFO]", msg, v)
}
func LogCritical(msg string, v ...any) {
	logger.Println("[CRITICAL]", msg, v)
}
func LogError(msg string, v ...any) {
	logger.Println("[INFO]", msg, v)
}
func LogDebug(msg string, v ...any) {
	logger.Println("[Debug]", msg, v)
}
