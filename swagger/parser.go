package swagger

import (
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/tinh-tinh/tinhtinh/core"
)

func (spec *SpecBuilder) ParsePaths(app *core.App) {
	mapperDoc := app.Module.MapperDoc

	pathObject := make(PathObject)
	definitions := make(map[string]*DefinitionObject)

	for tag, mx := range mapperDoc {
		for path, dtos := range mx {
			router := core.ParseRoute(path)
			parametes := []*ParameterObject{}
			for _, p := range dtos {
				switch p.In {
				case core.InBody:
					definitions[getNameStruct(p.Dto)] = ParseDefinition(p.Dto)
					parametes = append(parametes, &ParameterObject{
						Name: getNameStruct(p.Dto),
						In:   string(p.In),
						Schema: &SchemaObject{
							Ref: "#/definitions/" + getNameStruct(p.Dto),
						},
					})
				case core.InQuery:
					parametes = append(parametes, ScanQuery(p.Dto, p.In)...)
				case core.InPath:
					parametes = append(parametes, ScanQuery(p.Dto, p.In)...)
				}
			}

			if pathObject[router.Path] == nil {
				pathObject[router.Path] = &PathItemObject{}
			}
			itemObject := pathObject[router.Path]
			response := &ResponseObject{
				Description: "Ok",
			}
			res := map[string]*ResponseObject{"200": response}
			operation := &OperationObject{
				Tags:       []string{tag},
				Responses:  res,
				Parameters: parametes,
			}
			switch router.Method {
			case "GET":
				itemObject.Get = operation
			case "POST":
				itemObject.Post = operation
			case "PUT":
				itemObject.Put = operation
			case "DELETE":
				itemObject.Delete = operation
			}
		}
	}

	spec.Definitions = definitions
	spec.Paths = pathObject
}

type Mapper map[string]interface{}

func recursiveParseStandardSwagger(val interface{}) Mapper {
	mapper := make(Mapper)

	if reflect.ValueOf(val).IsNil() {
		return nil
	}
	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		key := firstLetterToLower(field.Name)
		if key == "ref" {
			key = "$ref"
		}
		if ct.Field(i).Interface() == nil {
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			ptrVal := recursiveParseStandardSwagger(ct.Field(i).Interface())
			if len(ptrVal) == 0 {
				continue
			}
			mapper[key] = ptrVal
		} else if field.Type.Kind() == reflect.Map {
			val := ct.Field(i).Interface()
			mapVal := reflect.ValueOf(val)
			subMapper := make(Mapper)
			for _, v := range mapVal.MapKeys() {
				subVal := recursiveParseStandardSwagger(mapVal.MapIndex(v).Interface())
				if IsNil(subVal) {
					continue
				}
				subKey := firstLetterToLower(v.String())
				subMapper[subKey] = subVal
			}
			mapper[key] = subMapper
		} else if field.Type.Kind() == reflect.Slice {
			arrVal := reflect.ValueOf(ct.Field(i).Interface())
			if arrVal.IsValid() {
				arr := []interface{}{}
				for i := 0; i < arrVal.Len(); i++ {
					item := arrVal.Index(i)
					if item.Kind() == reflect.Pointer {
						arr = append(arr, recursiveParseStandardSwagger(item.Interface()))
					} else {
						arr = append(arr, item.Interface())
					}
				}
				if IsNil(arr) {
					continue
				}
				mapper[key] = arr
			}
		} else {
			val := ct.Field(i).Interface()
			if IsNil(val) {
				continue
			}
			mapper[key] = ct.Field(i).Interface()
		}
	}

	if len(mapper) == 0 {
		return nil
	}

	return mapper
}

func firstLetterToLower(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func IsNil(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []*interface{}:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case []*SecuritySchemeObject:
		return len(v) == 0
	case []*ParameterObject:
		return len(v) == 0
	default:
		return val == nil
	}
}

func ParseDefinition(dto interface{}) *DefinitionObject {
	properties := make(map[string]*SchemaObject)
	ct := reflect.ValueOf(dto).Elem()
	for i := 0; i < ct.NumField(); i++ {
		schema := &SchemaObject{
			Type: ct.Field(i).Kind().String(),
		}

		field := ct.Type().Field(i)
		validator := field.Tag.Get("validate")
		isRequired := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool { return v == "required" })
		if isRequired == -1 {
			schema.Required = false
		} else {
			schema.Required = true
		}
		example := field.Tag.Get("example")
		if example != "" {
			schema.Example = example
		}

		properties[ct.Type().Field(i).Name] = schema
	}

	return &DefinitionObject{
		Type:       "object",
		Properties: properties,
	}
}

func getNameStruct(str interface{}) string {
	name := ""
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return firstLetterToLower(name)
}

func ScanQuery(val interface{}, in core.InDto) []*ParameterObject {
	ct := reflect.ValueOf(val).Elem()

	params := []*ParameterObject{}
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)

		param := &ParameterObject{
			Name: field.Name,
			Type: field.Type.Name(),
			In:   string(in),
		}
		validator := field.Tag.Get("validate")
		isRequired := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool { return v == "required" })
		if isRequired == -1 {
			param.Required = false
		} else {
			param.Required = true
		}
		example := field.Tag.Get("example")
		if example != "" {
			param.Default = example
		}

		params = append(params, param)
	}

	return params
}
