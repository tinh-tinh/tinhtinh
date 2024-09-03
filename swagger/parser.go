package swagger

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/tinh-tinh/tinhtinh/core"
)

func (spec *SpecBuilder) ParserPath(app *core.App) {
	mux := app.Module.MapMux

	groupRoute := make(map[string][]string)
	for tag, mx := range mux {
		for k := range mx {
			router := core.ParseRoute(k)
			if groupRoute[router.Path] == nil {
				groupRoute[router.Path] = make([]string, 0)
			}
			path := fmt.Sprintf("%s_%s", tag, router.Method)
			groupRoute[router.Path] = append(groupRoute[router.Path], path)
		}
	}
	fmt.Println(groupRoute)

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
}

type Mapper map[string]interface{}

func recursiveParse(val interface{}) Mapper {
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
			ptrVal := recursiveParse(ct.Field(i).Interface())
			if len(ptrVal) == 0 {
				continue
			}
			mapper[key] = ptrVal
		} else if field.Type.Kind() == reflect.Map {
			val := ct.Field(i).Interface()
			mapVal := reflect.ValueOf(val)
			subMapper := make(Mapper)
			for _, v := range mapVal.MapKeys() {
				subVal := recursiveParse(mapVal.MapIndex(v).Interface())
				if len(subVal) == 0 {
					continue
				}
				subMapper[v.String()] = subVal
			}
			mapper[key] = subMapper
		} else {
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
