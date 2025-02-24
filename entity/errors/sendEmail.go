package errEntity

import "awesomeProject/util"

// Source: https://docs.aws.amazon.com/ses/latest/APIReference/API_SendEmail.html

// SendEmailError is a CustomHttpError
type SendEmailError struct {
	message        string
	httpStatusCode int
}

func (e SendEmailError) GetHttpStatus() int {
	return e.httpStatusCode
}

func (e SendEmailError) Error() string {
	return e.message
}

var _ util.CustomHttpError = &SendEmailError{}

var (
	AccountSendingPaused = SendEmailError{
		"Email sending is disabled for your entire Amazon SES account. Enable it using UpdateAccountSendingEnabled.",
		400,
	}
	ConfigurationSetDoesNotExist = SendEmailError{
		"The configuration set does not exist.",
		400,
	}
	ConfigurationSetSendingPaused = SendEmailError{
		"Email sending is disabled for the configuration set. Enable it using UpdateConfigurationSetSendingEnabled.",
		400,
	}
	MailFromDomainNotVerified = SendEmailError{
		"Amazon SES could not read the MX record required to use the specified MAIL FROM domain.",
		400,
	}
	MessageRejected = SendEmailError{
		"The action failed, and the message could not be sent. Check the error stack for more details.",
		400,
	}
)
