package swagger

import (
	"encoding/json"
	"fmt"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"github.com/tinh-tinh/tinhtinh/core"
)

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		Info: &InfoObject{
			Version:        "1.0",
			Title:          "Swagger UI",
			Description:    "This is a sample server.",
			TermsOfService: "http://swagger.io/terms/",
			Contact: &ContactInfoObject{
				Name:  "API Support",
				Url:   "http://www.swagger.io/support",
				Email: "support@swagger.io",
			},
			License: &LicenseInfoObject{
				Name: "Apache 2.0",
				Url:  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		Swagger: "2.0",
		Schemes: []string{"http", "https"},
	}
}

func (spec *SpecBuilder) SetTitle(title string) *SpecBuilder {
	spec.Info.Title = title
	return spec
}

func (spec *SpecBuilder) SetDescription(description string) *SpecBuilder {
	spec.Info.Description = description
	return spec
}

func (spec *SpecBuilder) SetVersion(version string) *SpecBuilder {
	spec.Info.Version = version
	return spec
}

func (spec *SpecBuilder) SetHost(host string) *SpecBuilder {
	spec.Host = host
	return spec
}

func (spec *SpecBuilder) SetBasePath(basePath string) *SpecBuilder {
	spec.BasePath = basePath
	return spec
}

func (spec *SpecBuilder) Build() *SpecBuilder {
	return spec
}

func SetUp(path string, app *core.App, spec *SpecBuilder) {
	spec.ParserPath(app)
	mapper := recursiveParsePath(spec)
	jsonBytes, _ := json.Marshal(mapper)

	fmt.Println(string(jsonBytes))
	swaggerInfo := &swag.Spec{
		Version:          spec.Info.Version,
		Host:             spec.Host,
		BasePath:         spec.BasePath,
		Schemes:          spec.Schemes,
		Title:            spec.Info.Title,
		Description:      spec.Info.Description,
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(jsonBytes),
		LeftDelim:        "{{",
		RightDelim:       "}}",
	}

	route := fmt.Sprintf("%s%s", core.IfSlashPrefixString(app.Prefix), core.IfSlashPrefixString(path))

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
	app.Mux.Handle("GET "+route+"/*", httpSwagger.Handler(
		httpSwagger.URL(spec.Host+route+"/doc.json"),
	))
}
