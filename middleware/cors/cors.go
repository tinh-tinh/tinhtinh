// See document about cors: https://fetch.spec.whatwg.org/#http-cors-protocol
package cors

import (
	"net/http"

	"github.com/tinh-tinh/tinhtinh/common"
)

type Options struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may be:
	// * A exact string origins
	// * A regex pattern of origin
	// Eg.
	//   AllowedOrigins: []string{"https://foo.com", "http://foo.com"}
	//   AllowedOrigins: []string{"foo.com", "foo.com:1234", "^https://foo.com$"}
	//   AllowedOrigins: []string{"*"}
	AllowedOrigins []string
	// AllowOriginFunc is a custom function for processing the Origin header.
	// It returns true if allowed or false otherwise. If no Origin header is
	// present, it returns true by default.
	// Eg.
	//   AllowedOriginFunc: func(r *http.Request) bool {
	//     origin := r.Header["Origin"][0]
	//     return origin == "https://foo.com"
	//   }
	AllowedOriginCtx func(r *http.Request) bool
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all methods will be allowed.
	// Eg.
	//   AllowedMethods: []string{"GET", "POST"}
	//   AllowedMethods: []string{"*"}
	AllowedMethods []string
	// AllowedHeaders is a list of non-simple headers the client is allowed to use
	// with cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Eg.
	//   AllowedHeaders: []string{"Authorization", "Content-Type"}
	//   AllowedHeaders: []string{"*"}
	AllowedHeaders []string
	Credentials    bool
	PassThrough    bool
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

// NewCors returns a new CORS middleware.
//
// If options.AllowedOrigins is empty or if the special "*" value is present in
// the list, all origins will be allowed.
//
// If options.AllowedOriginCtx is not nil, it will be used to process the Origin
// header. If no Origin header is present, it returns true by default.
//
// If options.AllowedMethods is empty or if the special "*" value is present in
// the list, all methods will be allowed.
//
// If options.AllowedHeaders is empty or if the special "*" value is present in
// the list, all headers will be allowed.
//
// If options.Credentials is true, the middleware will include the
// Access-Control-Allow-Credentials header in the response.
//
// If options.PassThrough is true, the middleware will not intercept OPTIONS
// requests.
func NewCors(options Options) *Cors {
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

	if len(options.AllowedHeaders) == 0 {
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

// Handler returns a http.Handler that adds CORS headers to the responses.
//
// It handles two types of requests: OPTIONS requests (preflight) and actual
// requests.
//
// For OPTIONS requests, it sets the CORS headers and returns a 200 status code.
// If the Passthrough option is set to true, the handler will also pass the
// request to the underlying handler.
//
// For actual requests, it sets the CORS headers and passes the request to the
// underlying handler if the request is allowed according to the CORS policy.
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
			pass := cors.handleActualReq(w, r)
			if pass {
				h.ServeHTTP(w, r)
			}
		}
	})
}

// handleActualReq validates and sets the CORS headers for an actual request.
//
// It returns true if the request is allowed according to the CORS policy,
// and false otherwise.
func (cors *Cors) handleActualReq(w http.ResponseWriter, r *http.Request) bool {
	headers := w.Header()
	// Validate and set origin
	if !cors.isOriginAllowed(r) {
		common.UnauthorizedException(w, "Origin not allowed")
		return false
	}
	if cors.allowedOriginsAll {
		headers["Access-Control-Allow-Origin"] = []string{"*"}
	} else {
		headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
	}

	if !cors.isMethodAllowed(r.Method) {
		common.NotAllowedException(w, "Method not allowed")
		return false
	}

	if cors.credential {
		headers["Access-Control-Allow-Credentials"] = []string{"true"}
	}

	return true
}

// handlePreflight validates and sets the CORS headers for a preflight request.
//
// It handles the following steps:
//
// 1. Validate and set origin
// 2. Validate and set method
// 3. Validate and set headers
//
// If any of the steps fail, it returns an error response.
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

// isOriginAllowed checks if the given request's origin is allowed
// according to the CORS policy. If the policy is set to allow all
// origins, it returns true. If the policy is set to allow a specific
// set of origins, it checks if the request's origin is in the list.
// If the policy is set to a custom function, it calls the function
// and returns the result.
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

// isMethodAllowed checks if the given request's method is allowed
// according to the CORS policy. If the policy is set to allow all
// methods, it returns true. If the policy is set to allow a specific
// set of methods, it checks if the request's method is in the list.
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

// isHeaderAllowed checks if the given request's headers are allowed
// according to the CORS policy. If the policy is set to allow all
// headers, it returns true. If the policy is set to allow a specific
// set of headers, it checks if the request's headers are in the list.
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
