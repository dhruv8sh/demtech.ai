package util

// HttpException is used to transport HttpStatus and Error
type HttpException struct {
	StatusCode int
	Error      any
}

// HttpExceptionWithLog is used to transport HttpStatus, Error with a log
type HttpExceptionWithLog struct {
	StatusCode int
	Error      any
	Log        string
}

func newException(httpStatusCode int, errors any) HttpException {
	v, ok := errors.(string)
	if ok {
		errors = []string{v}
	}
	return HttpException{
		StatusCode: httpStatusCode,
		Error:      errors,
	}
}
func newExceptionWithLog(httpStatusCode int, errors any, msg string) HttpExceptionWithLog {
	v, ok := errors.(string)
	if ok {
		errors = []string{v}
	}
	return HttpExceptionWithLog{
		StatusCode: httpStatusCode,
		Error:      errors,
		Log:        msg,
	}
}

// HttpFail panics to be caught by the ClosureHandler()
func HttpFail(httpStatusCode int, error any) {
	panic(newException(httpStatusCode, error))
}

// HttpFailCustom panics with CustomHttpError
func HttpFailCustom(httpError CustomHttpError) {
	panic(newException(httpError.GetHttpStatus(), httpError.Error()))
}

// HttpFailWithLog is HttpFail but with log
func HttpFailWithLog(httpStatusCode int, error any, log string) {
	panic(newExceptionWithLog(httpStatusCode, error, log))
}
