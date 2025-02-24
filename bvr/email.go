package bvr

import (
	"awesomeProject/entity"
	"awesomeProject/repo"
	"awesomeProject/util"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	smtpHost = "www.example.com"
	smtpPort = "25"
	smtpUser = "exampleUser"
	smtpPass = "examplePass"
)

func AsyncSendEmail(messageId uint32, params *entity.SendEmailParams) {
	defer func() {
		if err := recover(); err != nil {
			repo.UpsertEmailAction(messageId, entity.Failure)
			util.LogError("bvr::AsyncSendEmail: panic recover", err)
		}
	}()

	simulateActionAfterDelay := entity.SendEmailAction(params.SimulateAction)
	if simulateActionAfterDelay.IsValid() {
		// Simulate action
		time.Sleep(10 * time.Second)
		repo.UpsertEmailAction(messageId, simulateActionAfterDelay)
		util.LogInfo("bvr::AsyncSendEmail: SimulatingAction", simulateActionAfterDelay)
	} else if prepareAndSendMockMail(params) {
		repo.UpsertEmailAction(messageId, entity.Success)
		util.LogInfo("bvr::AsyncSendEmail: mockEmailSend - Success")
	} else {
		repo.UpsertEmailAction(messageId, entity.Failure)
		util.LogInfo("bvr::AsyncSendEmail: mockEmailSend - Failure")
	}
}

// ------------------------------------------------------------
// Copied from the internet and then adjusted to look coherent|
// ------------------------------------------------------------
func prepareAndSendMockMail(params *entity.SendEmailParams) bool {
	var allRecipients []string
	allRecipients = append(allRecipients, params.Destination.ToAddresses...)
	allRecipients = append(allRecipients, params.Destination.CcAddresses...)
	allRecipients = append(allRecipients, params.Destination.BccAddresses...)

	// Construct the email message
	var msgBuilder strings.Builder
	msgBuilder.WriteString(fmt.Sprintf("From: %s\r\n", params.Source))
	msgBuilder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(params.Destination.ToAddresses, ", ")))
	if len(params.Destination.CcAddresses) > 0 {
		msgBuilder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(params.Destination.CcAddresses, ", ")))
	}
	msgBuilder.WriteString(fmt.Sprintf("Subject: %s\r\n", params.Message.Subject.Data))
	msgBuilder.WriteString("MIME-Version: 1.0\r\n")

	if params.Message.Body.Html != nil {
		msgBuilder.WriteString("Content-Type: text/html; charset=\"" + params.Message.Body.Html.Charset + "\"\r\n\r\n")
		msgBuilder.WriteString(params.Message.Body.Html.Data + "\r\n")
	} else if params.Message.Body.Text != nil {
		msgBuilder.WriteString("Content-Type: text/plain; charset=\"" + params.Message.Body.Text.Charset + "\"\r\n\r\n")
		msgBuilder.WriteString(params.Message.Body.Text.Data + "\r\n")
	}

	// Connect to SMTP server using TLS
	serverAddr := smtpHost + ":" + smtpPort
	conn, err := tls.Dial("tcp", serverAddr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	})
	if err != nil {
		util.LogError("bvr::AsyncSendEmail: connect to server failed", err)
		return false
	}
	defer conn.Close()

	// Create an SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		util.LogError("bvr::AsyncSendEmail: create client failed", err)
		return false
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	if err := client.Auth(auth); err != nil {
		util.LogError("bvr::AsyncSendEmail: auth failed", err)
		return false
	}

	// Set the sender
	if err := client.Mail(params.Source); err != nil {
		util.LogError("bvr::AsyncSendEmail: failed to set sender", err)
		return false
	}

	// Set recipients
	for _, recipient := range allRecipients {
		if err := client.Rcpt(recipient); err != nil {
			util.LogError("bvr::AsyncSendEmail: failed to set recipient", err)
			return false
		}
	}

	// Write the email content
	wc, err := client.Data()
	if err != nil {
		util.LogError("bvr::AsyncSendEmail: failed to open data writer", err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(msgBuilder.String()))
	if err != nil {
		util.LogError("bvr::AsyncSendEmail: failed to write message", err)
		return false
	}

	// Send the QUIT signal and close the connection
	err = client.Quit()
	if err != nil {
		util.LogError("bvr::AsyncSendEmail: Failed to quit smtp session", err)
		return false
	}

	util.LogInfo("bvr::AsyncSendEmail: Sent email successfully")
	return true
}
