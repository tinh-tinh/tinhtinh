package swagger

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/tinh-tinh/tinhtinh/core"
)

func (spec *SpecBuilder) ParserPath(app *core.App) {
	mapperDoc := app.Module.MapperDoc

	groupRoute := make(map[string][]string)
	definitions := make(map[string]*DefinitionObject)
	for tag, mx := range mapperDoc {
		for k, v := range mx {
			router := core.ParseRoute(k)
			if groupRoute[router.Path] == nil {
				groupRoute[router.Path] = make([]string, 0)
			}
			path := fmt.Sprintf("%s_%s", tag, router.Method)
			groupRoute[router.Path] = append(groupRoute[router.Path], path)
			for _, p := range v {
				definitions[getNameStruct(p.Dto)] = parseDto(p.Dto)
			}
		}
	}

	pathObject := make(PathObject)
	for k, v := range groupRoute {
		itemObject := &PathItemObject{}
		for i := 0; i < len(v); i++ {
			route := strings.Split(v[i], "_")
			response := &ResponseObject{
				Description: "Ok",
			}
			res := map[string]*ResponseObject{"200": response}
			operation := &OperationObject{
				Tags:      []string{route[0]},
				Responses: res,
			}
			switch route[1] {
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

		pathObject[k] = itemObject
	}

	spec.Paths = pathObject
	spec.Definitions = definitions
}

type Mapper map[string]interface{}

func recursiveParsePath(val interface{}) Mapper {
	mapper := make(Mapper)

	if reflect.ValueOf(val).IsNil() {
		return nil
	}
	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		key := firstLetterToLower(field.Name)
		if ct.Field(i).Interface() == nil {
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			ptrVal := recursiveParsePath(ct.Field(i).Interface())
			if len(ptrVal) == 0 {
				continue
			}
			mapper[key] = ptrVal
		} else if field.Type.Kind() == reflect.Map {
			val := ct.Field(i).Interface()
			mapVal := reflect.ValueOf(val)
			subMapper := make(Mapper)
			for _, v := range mapVal.MapKeys() {
				subVal := recursiveParsePath(mapVal.MapIndex(v).Interface())
				if len(subVal) == 0 {
					continue
				}
				subMapper[v.String()] = subVal
			}
			mapper[key] = subMapper
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

func parseDto(dto interface{}) *DefinitionObject {
	properties := make(map[string]*SchemaObject)
	ct := reflect.ValueOf(dto).Elem()
	for i := 0; i < ct.NumField(); i++ {
		schema := &SchemaObject{
			Type:     ct.Field(i).Kind().String(),
			Required: "true",
			Example:  "abc",
		}
		properties[ct.Type().Field(i).Name] = schema
	}

	return &DefinitionObject{
		Type:       "object",
		Properties: properties,
	}
}

func getNameStruct(str interface{}) string {
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
