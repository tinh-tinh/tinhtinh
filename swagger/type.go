package swagger

// -------- Info Object --------
type InfoObject struct {
	title          string
	description    string
	version        string
	termsOfService string
	contact        ContactInfoObject
	license        LicenseInfoObject
}

type ContactInfoObject struct {
	name  string
	url   string
	email string
}

type LicenseInfoObject struct {
	name string
	url  string
}

// -------- Path Object --------
type PathObject map[string]PathItemObject

// -------- Path Item Object --------
type PathItemObject struct {
	Ref        string
	Post       OperationObject
	Get        OperationObject
	Put        OperationObject
	Delete     OperationObject
	Parameters []ParameterObject
}

// -------- Operation Object --------
type OperationObject struct {
	Tags        []string
	Summary     string
	Description string
	OperationID string
	Consumes    []string
	Produces    []string
	Parameters  []ParameterObject
	Schemes     []string
	Deprecated  bool
	Security    []SecuritySchemeObject
	Responses   map[string]ResponseObject
}

// -------- Parameter Object --------
type ParameterObject struct {
	Name             string
	In               string
	Description      string
	Required         bool
	Type             string
	Items            map[string]string
	Format           string
	CollectionFormat string
}

// -------- Definition Object --------
type DefinitionSwagger struct {
	Type       string
	Required   []string
	Properties map[string]SchemaObject
}

// -------- Schema Object --------
type SchemaObject struct {
	Type     string
	Ref      string
	Format   string
	Required string
	Enum     []string
	Items    ItemsObject
}

// -------- Response Object --------
type ResponseObject struct {
	Description string
	Schema      SchemaObject
}

// -------- Items Object --------
type ItemsObject struct {
	Type     string
	Format   string
	Required string
	Enum     []string
}

// -------- Security Scheme Object --------
type SecuritySchemeObject struct {
	Type             string
	Description      string
	Name             string
	In               string
	Flow             string
	AuthorizationUrl string
	Token            string
	Scopes           map[string]string
}

// -------- Header Object --------
type HeaderObject struct {
	Description string
	Type        string
	Format      string
	Enum        []string
}

type SpecBuilder struct {
	swagger     string
	info        InfoObject
	schemes     []string
	Produces    []string
	Consumes    []string
	Host        string
	BasePath    string
	Paths       PathObject
	Definitions map[string]DefinitionSwagger
}
