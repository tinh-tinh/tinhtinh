package swagger

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func CreateDocument(server *http.ServeMux) {

	server.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
