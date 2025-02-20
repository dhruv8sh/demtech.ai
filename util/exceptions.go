package util

type HttpException struct {
	StatusCode int
	Error      any
}

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

func HttpFail(httpStatusCode int, error any) {
	panic(newException(httpStatusCode, error))
}

func HttpFailWithLog(httpStatusCode int, error any, log string) {
	panic(newExceptionWithLog(httpStatusCode, error, log))
}
