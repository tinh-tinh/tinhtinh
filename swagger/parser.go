package swagger

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/tinh-tinh/tinhtinh/core"
)

func (spec *SpecBuilder) ParserPath(app *core.App) {
	mux := app.Module.Mux

	groupRoute := make(map[string][]string)
	for k := range mux {
		router := core.ParseRoute(k)
		if groupRoute[router.Path] == nil {
			groupRoute[router.Path] = make([]string, 0)
		}
		groupRoute[router.Path] = append(groupRoute[router.Path], router.Method)
	}

	pathObject := make(PathObject)
	for k, v := range groupRoute {
		itemObject := &PathItemObject{}
		for i := 0; i < len(v); i++ {
			response := &ResponseObject{
				Description: "Ok",
			}
			res := map[string]*ResponseObject{"200": response}
			operation := &OperationObject{
				Responses: res,
			}
			switch v[i] {
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

	spec.Paths = &pathObject
}

type Mapper map[string]interface{}

func recursiveParse(val interface{}) Mapper {
	fmt.Printf("val is %v\n", val)
	mapper := make(Mapper)

	if reflect.TypeOf(val).Kind() == reflect.Map {
		fmt.Print("val is map\n")
		for k, v := range val.(map[string]interface{}) {
			mapper[k] = recursiveParse(v)
		}

		return mapper
	}

	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		key := firstLetterToLower(field.Name)
		if ct.Field(i).Interface() == nil {
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			mapper[key] = recursiveParse(ct.Field(i).Interface())
		} else {
			mapper[key] = ct.Field(i).Interface()
		}
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
