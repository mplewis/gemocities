package webproxys

import (
	"net/http"

	"git.sr.ht/~adnano/go-gemini"
)

var statusNames = map[gemini.Status]string{
	gemini.StatusInput:                    "Input required",
	gemini.StatusSensitiveInput:           "Sensitive input required",
	gemini.StatusSuccess:                  "Success",
	gemini.StatusRedirect:                 "Redirect",
	gemini.StatusPermanentRedirect:        "Permanent redirect",
	gemini.StatusTemporaryFailure:         "Temporary failure",
	gemini.StatusServerUnavailable:        "Server unavailable",
	gemini.StatusCGIError:                 "CGI error",
	gemini.StatusProxyError:               "Proxy error",
	gemini.StatusSlowDown:                 "Slow down",
	gemini.StatusPermanentFailure:         "Permanent failure",
	gemini.StatusNotFound:                 "Not found",
	gemini.StatusGone:                     "Gone",
	gemini.StatusProxyRequestRefused:      "Proxy request refused",
	gemini.StatusBadRequest:               "Bad request",
	gemini.StatusCertificateRequired:      "Certificate required",
	gemini.StatusCertificateNotAuthorized: "Certificate not authorized",
	gemini.StatusCertificateNotValid:      "Certificate not valid",
}

var statusToHTTP = map[gemini.Status]int{
	gemini.StatusInput:                    http.StatusBadRequest,
	gemini.StatusSensitiveInput:           http.StatusBadRequest,
	gemini.StatusSuccess:                  http.StatusOK,
	gemini.StatusRedirect:                 http.StatusFound,
	gemini.StatusPermanentRedirect:        http.StatusPermanentRedirect,
	gemini.StatusTemporaryFailure:         http.StatusInternalServerError,
	gemini.StatusServerUnavailable:        http.StatusServiceUnavailable,
	gemini.StatusCGIError:                 http.StatusInternalServerError,
	gemini.StatusProxyError:               http.StatusBadGateway,
	gemini.StatusSlowDown:                 http.StatusTooManyRequests,
	gemini.StatusPermanentFailure:         http.StatusInternalServerError,
	gemini.StatusNotFound:                 http.StatusNotFound,
	gemini.StatusGone:                     http.StatusGone,
	gemini.StatusProxyRequestRefused:      http.StatusForbidden,
	gemini.StatusBadRequest:               http.StatusBadRequest,
	gemini.StatusCertificateRequired:      http.StatusForbidden,
	gemini.StatusCertificateNotAuthorized: http.StatusForbidden,
	gemini.StatusCertificateNotValid:      http.StatusForbidden,
}
