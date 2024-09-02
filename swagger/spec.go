package swagger

import (
	"encoding/json"

	"github.com/swaggo/swag"
	"github.com/tinh-tinh/tinhtinh/core"
)

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		Info: &InfoObject{
			Version:        "1.0",
			Title:          "Swagger Example API for Tinh Tinh",
			Description:    "This is a sample server Tinh tinh server.",
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
		Swagger:  "2.0",
		Schemes:  []string{"http", "https"},
		BasePath: "/v1",
		Host:     "tinhtinh.swagger.io",
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

func (spec *SpecBuilder) Build() *SpecBuilder {
	return spec
}

func SetUp(app *core.App, spec *SpecBuilder) {
	spec.ParserPath(app)
	mapper := recursiveParse(spec)
	jsonBytes, _ := json.Marshal(mapper)
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

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}
