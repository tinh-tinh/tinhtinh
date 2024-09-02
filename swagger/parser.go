package swagger

import "github.com/tinh-tinh/tinhtinh/core"

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
		itemObject := PathItemObject{}
		for i := 0; i < len(v); i++ {
			operation := OperationObject{}
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

	spec.Paths = pathObject
}
