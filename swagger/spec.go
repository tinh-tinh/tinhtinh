package swagger

import (
	"encoding/json"

	"github.com/swaggo/swag"
	"github.com/tinh-tinh/tinhtinh/core"
)

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{}
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

	jsonBytes, _ := json.Marshal(spec)
	swaggerInfo := &swag.Spec{
		Version:         spec.Info.Version,
		Host:            spec.Host,
		BasePath:        spec.BasePath,
		Schemes:         spec.Schemes,
		Title:           spec.Info.Title,
		Description:     spec.Info.Description,
		SwaggerTemplate: string(jsonBytes),
		LeftDelim:       "{{",
		RightDelim:      "}}",
	}

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}
