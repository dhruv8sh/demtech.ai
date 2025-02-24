package errEntity

import (
	"awesomeProject/util"
	"net/http"
)

// Source: https://docs.aws.amazon.com/ses/latest/APIReference/CommonErrors.html

type CommonError struct {
	message        string
	httpStatusCode int
}

func (e CommonError) GetHttpStatus() int {
	return e.httpStatusCode
}

func (e CommonError) Error() string {
	return e.message
}

var _ util.CustomHttpError = CommonError{}

func (e CommonError) TriggerCritical(msg string) {
	util.LogCritical(msg, e.message)
}

var AccessDeniedException = CommonError{
	"You do not have sufficient access to perform this action.",
	http.StatusBadRequest,
}
var IncompleteSignature = CommonError{
	"The request signature does not conform to AWS standards.",
	400,
}
var InternalFailure = CommonError{
	"The request processing has failed because of an unknown error, exception or failure.",
	500,
}
var InvalidAction = CommonError{
	"The action or operation requested is invalid. Verify that the action is typed correctly.",
	400,
}
var InvalidClientTokenId = CommonError{
	"The X.509 certificate or AWS access key ID provided does not exist in our records.",
	403,
}
var InvalidParameterCombination = CommonError{
	"Parameters that must not be used together were used together.",
	400,
}
var InvalidParameterValue = CommonError{
	"An invalid or out-of-range value was supplied for the input parameter.",
	400,
}
var InvalidQueryParameter = CommonError{
	"The AWS query string is malformed or does not adhere to AWS standards.",
	400,
}
var MalformedQueryString = CommonError{
	"The query string contains a syntax error.",
	404,
}
var MissingAction = CommonError{
	"The request is missing an action or a required parameter.",
	400,
}
var MissingAuthenticationToken = CommonError{
	"The request must contain either a valid (registered) AWS access key ID or X.509 certificate.",
	403,
}
var MissingParameter = CommonError{
	"A required parameter for the specified action is not supplied.",
	400,
}
var NotAuthorized = CommonError{
	"You do not have permission to perform this action.",
	400,
}
var OptInRequired = CommonError{
	"The AWS access key ID needs a subscription for the service.",
	403,
}
var RequestExpired = CommonError{
	"The request reached the service more than 15 minutes after the date stamp on the request or more than 15 minutes after the request expiration date (such as for pre-signed URLs), or the date stamp on the request is more than 15 minutes in the future.",
	400,
}
var ServiceUnavailable = CommonError{
	"The request has failed due to a temporary failure of the server.",
	503,
}
var ThrottlingException = CommonError{
	"The request was denied due to request throttling.",
	400,
}
var ValidationError = CommonError{
	"The input fails to satisfy the constraints specified by an AWS service.",
	400,
}
