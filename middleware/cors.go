// See document about cors: https://fetch.spec.whatwg.org/#http-cors-protocol
package middleware

import (
	"net/http"

	"github.com/tinh-tinh/tinhtinh/common"
)

type CorsOptions struct {
	AllowedOrigins   []string
	AllowedOriginCtx func(r *http.Request) bool
	AllowedMethods   []string
	AllowedHeaders   []string
	Credentials      bool
	PassThrough      bool
}

type Cors struct {
	allowedOrigins   []string
	allowedOriginFnc func(r *http.Request) bool
	allowedMethods   []string
	allowedHeaders   map[string]interface{}
	credential       bool
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	optionPassthrough bool
}

func NewCors(options CorsOptions) *Cors {
	c := &Cors{
		allowedOrigins:    options.AllowedOrigins,
		allowedOriginFnc:  options.AllowedOriginCtx,
		allowedMethods:    options.AllowedMethods,
		credential:        options.Credentials,
		optionPassthrough: options.PassThrough,
		allowedHeaders:    make(map[string]interface{}),
	}

	if c.allowedOriginFnc == nil && len(c.allowedOrigins) == 0 {
		c.allowedOriginsAll = true
	}

	for _, origin := range c.allowedOrigins {
		if origin == "*" {
			c.allowedOriginsAll = true
			c.allowedOrigins = nil
			c.allowedOriginFnc = nil
			break
		}
	}

	if options.AllowedHeaders == nil || len(options.AllowedHeaders) == 0 {
		c.allowedHeadersAll = true
	}

	for _, header := range options.AllowedHeaders {
		if header == "*" {
			c.allowedHeadersAll = true
			c.allowedHeaders = nil
			break
		} else {
			c.allowedHeaders[header] = true
		}
	}

	if len(c.allowedMethods) == 0 {
		c.allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	}

	return c
}

func (cors Cors) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			cors.handlePreflight(w, r)
			if cors.optionPassthrough {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			cors.handleActualReq(w, r)
			h.ServeHTTP(w, r)
		}
	})
}

func (cors *Cors) handleActualReq(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	// Validate and set origin
	if !cors.isOriginAllowed(r) {
		common.UnauthorizedException(w, "Origin not allowed")
		return
	}
	if cors.allowedOriginsAll {
		headers["Access-Control-Allow-Origin"] = []string{"*"}
	} else {
		headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
	}

	if !cors.isMethodAllowed(r.Method) {
		common.NotAllowedException(w, "Method not allowed")
		return
	}

	if cors.credential {
		headers["Access-Control-Allow-Credentials"] = []string{"true"}
	}
}

func (cors *Cors) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	// Validate and set origin
	if !cors.isOriginAllowed(r) {
		common.UnauthorizedException(w, "Origin not allowed")
		return
	}
	if cors.allowedOriginsAll {
		headers["Access-Control-Allow-Origin"] = []string{"*"}
	} else {
		headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
	}

	// Validate and set method
	if r.Method != http.MethodOptions {
		common.NotAllowedException(w, "Method not allowed")
		return
	}

	method := r.Header.Get("Access-Control-Request-Method")
	if !cors.isMethodAllowed(method) {
		common.NotAllowedException(w, "Method not allowed")
		return
	}
	headers["Access-Control-Allow-Methods"] = r.Header["Access-Control-Request-Method"]

	// Validate and set headers
	if !cors.isHeaderAllowed(r) {
		common.NotAllowedException(w, "Header not allowed")
		return
	}
	headers["Access-Control-Allow-Headers"] = r.Header["Access-Control-Request-Headers"]

	if cors.credential {
		headers["Access-Control-Allow-Credentials"] = []string{"true"}
	}
}

func (cors *Cors) isOriginAllowed(r *http.Request) bool {
	if cors.allowedOriginFnc != nil {
		return cors.allowedOriginFnc(r)
	}

	if cors.allowedOriginsAll {
		return true
	}

	for _, origin := range cors.allowedOrigins {
		if origin == r.Header.Get("Origin") {
			return true
		}
	}

	return false
}

func (cors Cors) isMethodAllowed(method string) bool {
	if method == http.MethodOptions {
		return true
	}
	for _, m := range cors.allowedMethods {
		if m == method {
			return true
		}
	}

	return false
}

func (cors Cors) isHeaderAllowed(r *http.Request) bool {
	if cors.allowedHeadersAll {
		return true
	}

	reqHeaders, found := r.Header["Access-Control-Request-Headers"]
	if found {
		for _, header := range reqHeaders {
			if cors.allowedHeaders[header] == nil {
				return false
			}
		}
	}

	return true
}
