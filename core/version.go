package core

import (
	"net/http"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/common"
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

func (app *App) versionMiddleware(routers []*Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var version string
		if app.version != nil {
			version = app.version.Get(r)
		}
		if version == "" || len(routers) == 1 {
			routers[0].getHandler(app).ServeHTTP(w, r)
		} else {
			idxVersion := slices.IndexFunc(routers, func(e *Router) bool {
				return e.Version == version
			})

			if idxVersion == -1 {
				common.InternalServerException(w, "version not found")
				return
			}

			routers[idxVersion].getHandler(app).ServeHTTP(w, r)
		}
	})
}
