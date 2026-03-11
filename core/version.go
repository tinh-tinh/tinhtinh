package core

import (
	"net/http"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type VersionType string

const (
	URIVersion       VersionType = "uri"
	HeaderVersion    VersionType = "header"
	MediaTypeVersion VersionType = "mediaType"
	CustomVersion    VersionType = "custom"
)

type VersionOptions struct {
	Type      VersionType
	Header    string
	Key       string
	Extractor func(*http.Request) string
}

type Version struct {
	Type      VersionType
	Header    string
	Key       string
	Extractor func(*http.Request) string
}

// EnableVersioning enables versioning on the API server. The passed in options are used
// to configure the versioning middleware.
//
// The versioning middleware will extract the version from the request based on the
// Type field of the options. The version will then be used to select the correct
// router for the request.
//
// The following types are supported:
//
// - "uri": The version is extracted from the URI path.
// - "header": The version is extracted from the specified header.
// - "mediaType": The version is extracted from the specified media type.
// - "custom": The version is extracted using the specified extractor function.
func (app *App) EnableVersioning(opt VersionOptions) *App {
	app.version = &Version{
		Type: opt.Type,
	}

	switch opt.Type {
	case HeaderVersion:
		app.version.Header = opt.Header
	case MediaTypeVersion:
		app.version.Key = opt.Key
	case CustomVersion:
		app.version.Extractor = opt.Extractor
	}

	return app
}

// Get returns the version of the request. The version is extracted based on the
// Type field of the Version object.
//
// The following methods are supported:
//
//   - HeaderVersion: The version is extracted from the header specified in the
//     Header field.
//   - MediaTypeVersion: The version is extracted from the media type specified in
//     the Key field.
//   - CustomVersion: The version is extracted using the extractor function specified
//     in the Extractor field.
//
// If the type is not supported, an empty string is returned.
func (v *Version) Get(r *http.Request) string {
	switch v.Type {
	case HeaderVersion:
		return r.Header.Get(v.Header)
	case MediaTypeVersion:
		metadata := r.Header.Get("Accept")
		idx := strings.Index(metadata, v.Key)

		start := idx + len(v.Key)
		end := strings.Index(metadata[start:], ";")

		if end == -1 {
			return metadata[start:]
		} else {
			return metadata[start : start+end]
		}
	case CustomVersion:
		return v.Extractor(r)
	default:
		return ""
	}
}

// versionMiddleware returns a middleware that dispatches the request to the correct
// router based on the version of the request. The version is determined by the
// Version object of the App.
//
// If the version of the request is empty or there is only one router, the request
// is dispatched to the first router. Otherwise, the router with the matching
// version is selected. If no router matches the version, a 500 error is returned.
func (app *App) versionMiddleware(routers []*Router) http.Handler {
	for _, r := range routers {
		r.finalHandler = r.getHandler(app)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var version string
		if app.version != nil {
			version = app.version.Get(r)
		}
		if version == "" || len(routers) == 1 {
			routers[0].finalHandler.ServeHTTP(w, r)
		} else {
			idxVersion := slices.IndexFunc(routers, func(e *Router) bool {
				return e.Version == version
			})

			if idxVersion == -1 {
				common.InternalServerException(w, "version not found")
				return
			}

			routers[idxVersion].finalHandler.ServeHTTP(w, r)
		}
	})
}
