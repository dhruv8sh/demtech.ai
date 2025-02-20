package controller

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func SendEmail(c *gin.Context) {
	var params struct {
		XMLName     xml.Name `xml:"SendEmailRequest"`
		Xmlns       string   `xml:"xmlns,attr"`
		Source      string   `xml:"Source"`
		Destination struct {
			ToAddresses  []string `xml:"ToAddresses>member,omitempty"`
			CcAddresses  []string `xml:"CcAddresses>member,omitempty"`
			BccAddresses []string `xml:"BccAddresses>member,omitempty"`
		} `xml:"Destination"`
		Message struct {
			Subject struct {
				Data    string `xml:"Data"`
				Charset string `xml:"Charset,omitempty"`
			} `xml:"Subject"`
			Body struct {
				Text *struct {
					Data    string `xml:"Data"`
					Charset string `xml:"Charset,omitempty"`
				} `xml:"Text,omitempty"`
				Html *struct {
					Data    string `xml:"Data"`
					Charset string `xml:"Charset,omitempty"`
				} `xml:"Html,omitempty"`
			} `xml:"Body"`
		} `xml:"Message"`
		ReplyToAddresses []string `xml:"ReplyToAddresses>member,omitempty"`
		ReturnPath       string   `xml:"ReturnPath,omitempty"`
		ConfigurationSet string   `xml:"ConfigurationSetName,omitempty"`
	}
}

func GetEmailStatus(c *gin.Context) {

}

func GetQuota(c *gin.Context) {

}

func GetMetrics(c *gin.Context) {

}
