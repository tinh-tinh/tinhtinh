package swagger

import "github.com/swaggo/swag"

type SpecBuilder struct {
	Title       string
	Description string
	Version     string
}

func NewSpec() *SpecBuilder {
	return &SpecBuilder{}
}

func (s *SpecBuilder) SetTitle(title string) *SpecBuilder {
	s.Title = title
	return s
}

func (s *SpecBuilder) SetDescription(description string) *SpecBuilder {
	s.Description = description
	return s
}

func (s *SpecBuilder) SetVersion(version string) *SpecBuilder {
	s.Version = version
	return s
}

func (s *SpecBuilder) Build() {
	swaggerInfo := &swag.Spec{
		Title:            s.Title,
		Description:      s.Description,
		Version:          s.Version,
		Schemes:          []string{"http", "https"},
		Host:             "",
		BasePath:         "",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
		LeftDelim:        "{{",
		RightDelim:       "}}",
	}
	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {}
}`
