package exception

import (
	"encoding/json"
	"net/http"
)

type Http struct {
	Status int
	Msg    string
}

func ThrowHttp(msg string, status int) Http {
	return Http{Status: status, Msg: msg}
}

func BadRequest(msg string) Http {
	return ThrowHttp(msg, http.StatusBadRequest)
}

func Unauthorized(msg string) Http {
	return ThrowHttp(msg, http.StatusUnauthorized)
}

func Forbidden(msg string) Http {
	return ThrowHttp(msg, http.StatusForbidden)
}

func NotFound(msg string) Http {
	return ThrowHttp(msg, http.StatusNotFound)
}

func MethodNotAllowed(msg string) Http {
	return ThrowHttp(msg, http.StatusMethodNotAllowed)
}

func NotAcceptable(msg string) Http {
	return ThrowHttp(msg, http.StatusNotAcceptable)
}

func RequestTimeout(msg string) Http {
	return ThrowHttp(msg, http.StatusRequestTimeout)
}

func Conflict(msg string) Http {
	return ThrowHttp(msg, http.StatusConflict)
}

func Gone(msg string) Http {
	return ThrowHttp(msg, http.StatusGone)
}

func LengthRequired(msg string) Http {
	return ThrowHttp(msg, http.StatusLengthRequired)
}

func PreconditionFailed(msg string) Http {
	return ThrowHttp(msg, http.StatusPreconditionFailed)
}

func ContentTooLarge(msg string) Http {
	return ThrowHttp(msg, http.StatusRequestEntityTooLarge)
}

func RequestURITooLong(msg string) Http {
	return ThrowHttp(msg, http.StatusRequestURITooLong)
}

func UnsupportedMediaType(msg string) Http {
	return ThrowHttp(msg, http.StatusUnsupportedMediaType)
}

func RequestedRangeNotSatisfiable(msg string) Http {
	return ThrowHttp(msg, http.StatusRequestedRangeNotSatisfiable)
}

func ExpectationFailed(msg string) Http {
	return ThrowHttp(msg, http.StatusExpectationFailed)
}

func IamATeapot(msg string) Http {
	return ThrowHttp(msg, http.StatusTeapot)
}

func UnprocessableEntity(msg string) Http {
	return ThrowHttp(msg, http.StatusUnprocessableEntity)
}

func Locked(msg string) Http {
	return ThrowHttp(msg, http.StatusLocked)
}

func FailedDependency(msg string) Http {
	return ThrowHttp(msg, http.StatusFailedDependency)
}

func UpgradeRequired(msg string) Http {
	return ThrowHttp(msg, http.StatusUpgradeRequired)
}

func PreconditionRequired(msg string) Http {
	return ThrowHttp(msg, http.StatusPreconditionRequired)
}

func TooManyRequests(msg string) Http {
	return ThrowHttp(msg, http.StatusTooManyRequests)
}

func InternalServer(msg string) Http {
	return ThrowHttp(msg, http.StatusInternalServerError)
}

func NotImplemented(msg string) Http {
	return ThrowHttp(msg, http.StatusNotImplemented)
}

func BadGateway(msg string) Http {
	return ThrowHttp(msg, http.StatusBadGateway)
}

func ServiceUnavailable(msg string) Http {
	return ThrowHttp(msg, http.StatusServiceUnavailable)
}

func GatewayTimeout(msg string) Http {
	return ThrowHttp(msg, http.StatusGatewayTimeout)
}

func HttpVersionNotSupported(msg string) Http {
	return ThrowHttp(msg, http.StatusHTTPVersionNotSupported)
}

func (e Http) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func AdapterHttpError(err error) Http {
	var e Http
	er := json.Unmarshal([]byte(err.Error()), &e)
	if er != nil {
		return Http{Status: http.StatusInternalServerError, Msg: err.Error()}
	}
	return e
}
