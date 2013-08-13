package tigertonic

import "net/http"

// SnakeCaseHTTPEquivErrors being true will cause tigertonic.HTTPEquivError
// error responses to be written as (for example) "not_found" rather than
// "tigertonic.NotFound".
var SnakeCaseHTTPEquivErrors bool

type HTTPEquivError interface {
	error
	Status() int
}

type Continue struct {
	error
}

func (err Continue) Status() int { return http.StatusContinue }

type SwitchingProtocols struct {
	error
}

func (err SwitchingProtocols) Status() int { return http.StatusSwitchingProtocols }

type OK struct {
	error
}

func (err OK) Status() int { return http.StatusOK }

type Created struct {
	error
}

func (err Created) Status() int { return http.StatusCreated }

type Accepted struct {
	error
}

func (err Accepted) Status() int { return http.StatusAccepted }

type NonAuthoritativeInfo struct {
	error
}

func (err NonAuthoritativeInfo) Status() int { return http.StatusNonAuthoritativeInfo }

type NoContent struct {
	error
}

func (err NoContent) Status() int { return http.StatusNoContent }

type ResetContent struct {
	error
}

func (err ResetContent) Status() int { return http.StatusResetContent }

type PartialContent struct {
	error
}

func (err PartialContent) Status() int { return http.StatusPartialContent }

type MultipleChoices struct {
	error
}

func (err MultipleChoices) Status() int { return http.StatusMultipleChoices }

type MovedPermanently struct {
	error
}

func (err MovedPermanently) Status() int { return http.StatusMovedPermanently }

type Found struct {
	error
}

func (err Found) Status() int { return http.StatusFound }

type SeeOther struct {
	error
}

func (err SeeOther) Status() int { return http.StatusSeeOther }

type NotModified struct {
	error
}

func (err NotModified) Status() int { return http.StatusNotModified }

type UseProxy struct {
	error
}

func (err UseProxy) Status() int { return http.StatusUseProxy }

type TemporaryRedirect struct {
	error
}

func (err TemporaryRedirect) Status() int { return http.StatusTemporaryRedirect }

type BadRequest struct {
	error
}

func (err BadRequest) Status() int { return http.StatusBadRequest }

type Unauthorized struct {
	error
}

func (err Unauthorized) Status() int { return http.StatusUnauthorized }

type PaymentRequired struct {
	error
}

func (err PaymentRequired) Status() int { return http.StatusPaymentRequired }

type Forbidden struct {
	error
}

func (err Forbidden) Status() int { return http.StatusForbidden }

type NotFound struct {
	error
}

func (err NotFound) Status() int { return http.StatusNotFound }

type MethodNotAllowed struct {
	error
}

func (err MethodNotAllowed) Status() int { return http.StatusMethodNotAllowed }

type NotAcceptable struct {
	error
}

func (err NotAcceptable) Status() int { return http.StatusNotAcceptable }

type ProxyAuthRequired struct {
	error
}

func (err ProxyAuthRequired) Status() int { return http.StatusProxyAuthRequired }

type RequestTimeout struct {
	error
}

func (err RequestTimeout) Status() int { return http.StatusRequestTimeout }

type Conflict struct {
	error
}

func (err Conflict) Status() int { return http.StatusConflict }

type Gone struct {
	error
}

func (err Gone) Status() int { return http.StatusGone }

type LengthRequired struct {
	error
}

func (err LengthRequired) Status() int { return http.StatusLengthRequired }

type PreconditionFailed struct {
	error
}

func (err PreconditionFailed) Status() int { return http.StatusPreconditionFailed }

type RequestEntityTooLarge struct {
	error
}

func (err RequestEntityTooLarge) Status() int { return http.StatusRequestEntityTooLarge }

type RequestURITooLong struct {
	error
}

func (err RequestURITooLong) Status() int { return http.StatusRequestURITooLong }

type UnsupportedMediaType struct {
	error
}

func (err UnsupportedMediaType) Status() int { return http.StatusUnsupportedMediaType }

type RequestedRangeNotSatisfiable struct {
	error
}

func (err RequestedRangeNotSatisfiable) Status() int { return http.StatusRequestedRangeNotSatisfiable }

type ExpectationFailed struct {
	error
}

func (err ExpectationFailed) Status() int { return http.StatusExpectationFailed }

type Teapot struct {
	error
}

func (err Teapot) Status() int { return http.StatusTeapot }

type InternalServerError struct {
	error
}

func (err InternalServerError) Status() int { return http.StatusInternalServerError }

type NotImplemented struct {
	error
}

func (err NotImplemented) Status() int { return http.StatusNotImplemented }

type BadGateway struct {
	error
}

func (err BadGateway) Status() int { return http.StatusBadGateway }

type ServiceUnavailable struct {
	error
}

func (err ServiceUnavailable) Status() int { return http.StatusServiceUnavailable }

type GatewayTimeout struct {
	error
}

func (err GatewayTimeout) Status() int { return http.StatusGatewayTimeout }

type HTTPVersionNotSupported struct {
	error
}

func (err HTTPVersionNotSupported) Status() int { return http.StatusHTTPVersionNotSupported }
