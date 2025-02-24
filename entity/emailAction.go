package entity

import "encoding/xml"

type SendEmailAction string

const (
	InQueue SendEmailAction = "IN_QUEUE"
	Success SendEmailAction = "SUCCESS"
	Failure SendEmailAction = "FAILURE"
)

func (e SendEmailAction) IsValid() bool {
	return e == InQueue || e == Success || e == Failure
}

// SendEmailParams is the exact structure of AWS SES SendEmail
// https://docs.aws.amazon.com/ses/latest/APIReference/API_SendEmail.html
// Based off of XML, since AWS uses XML
type SendEmailParams struct {
	XMLName          xml.Name `xml:"SendEmailRequest" validate:"required"`
	Xmlns            string   `xml:"xmlns"`
	Source           string   `xml:"Source" validate:"required|email"`
	ReplyToAddresses []string `xml:"ReplyToAddresses" validate:"emails"`
	ReturnPath       string   `xml:"ReturnPath"`
	ConfigurationSet string   `xml:"ConfigurationSetName"`
	SimulateAction   string   `xml:"SimulateAction"`
	Destination      struct {
		ToAddresses  []string `xml:"ToAddresses" validate:"emails"`
		CcAddresses  []string `xml:"CcAddresses" validate:"emails"`
		BccAddresses []string `xml:"BccAddresses" validate:"emails"`
	} `xml:"Destination" validate:"required"`
	Message struct {
		Subject struct {
			Data    string `xml:"Data"`
			Charset string `xml:"Charset"`
		} `xml:"Subject"`
		Body struct {
			Text *struct {
				Data    string `xml:"Data"`
				Charset string `xml:"Charset"`
			} `xml:"Text,omitempty"`
			Html *struct {
				Data    string `xml:"Data"`
				Charset string `xml:"Charset"`
			} `xml:"Html"`
		} `xml:"Body"`
	} `xml:"Message" validate:"required"`
}
