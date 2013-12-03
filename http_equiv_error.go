package tigertonic

import "net/http"

// SnakeCaseHTTPEquivErrors being true will cause tigertonic.HTTPEquivError
// error responses to be written as (for example) "not_found" rather than
// "tigertonic.NotFound".
var SnakeCaseHTTPEquivErrors bool

// Err is an alias for the built-in error type so that it can be publicly
// exported when embedding.
type Err error

type HTTPEquivError interface {
	error
	Status() int
}

type Continue struct {
	Err
}

func (err Continue) Status() int { return http.StatusContinue }

type SwitchingProtocols struct {
	Err
}

func (err SwitchingProtocols) Status() int { return http.StatusSwitchingProtocols }

type OK struct {
	Err
}

func (err OK) Status() int { return http.StatusOK }

type Created struct {
	Err
}

func (err Created) Status() int { return http.StatusCreated }

type Accepted struct {
	Err
}

func (err Accepted) Status() int { return http.StatusAccepted }

type NonAuthoritativeInfo struct {
	Err
}

func (err NonAuthoritativeInfo) Status() int { return http.StatusNonAuthoritativeInfo }

type NoContent struct {
	Err
}

func (err NoContent) Status() int { return http.StatusNoContent }

type ResetContent struct {
	Err
}

func (err ResetContent) Status() int { return http.StatusResetContent }

type PartialContent struct {
	Err
}

func (err PartialContent) Status() int { return http.StatusPartialContent }

type MultipleChoices struct {
	Err
}

func (err MultipleChoices) Status() int { return http.StatusMultipleChoices }

type MovedPermanently struct {
	Err
}

func (err MovedPermanently) Status() int { return http.StatusMovedPermanently }

type Found struct {
	Err
}

func (err Found) Status() int { return http.StatusFound }

type SeeOther struct {
	Err
}

func (err SeeOther) Status() int { return http.StatusSeeOther }

type NotModified struct {
	Err
}

func (err NotModified) Status() int { return http.StatusNotModified }

type UseProxy struct {
	Err
}

func (err UseProxy) Status() int { return http.StatusUseProxy }

type TemporaryRedirect struct {
	Err
}

func (err TemporaryRedirect) Status() int { return http.StatusTemporaryRedirect }

type BadRequest struct {
	Err
}

func (err BadRequest) Status() int { return http.StatusBadRequest }

type Unauthorized struct {
	Err
}

func (err Unauthorized) Status() int { return http.StatusUnauthorized }

type PaymentRequired struct {
	Err
}

func (err PaymentRequired) Status() int { return http.StatusPaymentRequired }

type Forbidden struct {
	Err
}

func (err Forbidden) Status() int { return http.StatusForbidden }

type NotFound struct {
	Err
}

func (err NotFound) Status() int { return http.StatusNotFound }

type MethodNotAllowed struct {
	Err
}

func (err MethodNotAllowed) Status() int { return http.StatusMethodNotAllowed }

type NotAcceptable struct {
	Err
}

func (err NotAcceptable) Status() int { return http.StatusNotAcceptable }

type ProxyAuthRequired struct {
	Err
}

func (err ProxyAuthRequired) Status() int { return http.StatusProxyAuthRequired }

type RequestTimeout struct {
	Err
}

func (err RequestTimeout) Status() int { return http.StatusRequestTimeout }

type Conflict struct {
	Err
}

func (err Conflict) Status() int { return http.StatusConflict }

type Gone struct {
	Err
}

func (err Gone) Status() int { return http.StatusGone }

type LengthRequired struct {
	Err
}

func (err LengthRequired) Status() int { return http.StatusLengthRequired }

type PreconditionFailed struct {
	Err
}

func (err PreconditionFailed) Status() int { return http.StatusPreconditionFailed }

type RequestEntityTooLarge struct {
	Err
}

func (err RequestEntityTooLarge) Status() int { return http.StatusRequestEntityTooLarge }

type RequestURITooLong struct {
	Err
}

func (err RequestURITooLong) Status() int { return http.StatusRequestURITooLong }

type UnsupportedMediaType struct {
	Err
}

func (err UnsupportedMediaType) Status() int { return http.StatusUnsupportedMediaType }

type RequestedRangeNotSatisfiable struct {
	Err
}

func (err RequestedRangeNotSatisfiable) Status() int { return http.StatusRequestedRangeNotSatisfiable }

type ExpectationFailed struct {
	Err
}

func (err ExpectationFailed) Status() int { return http.StatusExpectationFailed }

type Teapot struct {
	Err
}

func (err Teapot) Status() int { return http.StatusTeapot }

type InternalServerError struct {
	Err
}

func (err InternalServerError) Status() int { return http.StatusInternalServerError }

type NotImplemented struct {
	Err
}

func (err NotImplemented) Status() int { return http.StatusNotImplemented }

type BadGateway struct {
	Err
}

func (err BadGateway) Status() int { return http.StatusBadGateway }

type ServiceUnavailable struct {
	Err
}

func (err ServiceUnavailable) Status() int { return http.StatusServiceUnavailable }

type GatewayTimeout struct {
	Err
}

func (err GatewayTimeout) Status() int { return http.StatusGatewayTimeout }

type HTTPVersionNotSupported struct {
	Err
}

func (err HTTPVersionNotSupported) Status() int { return http.StatusHTTPVersionNotSupported }
