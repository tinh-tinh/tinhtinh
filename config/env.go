package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/utils"
)

func Scan(env interface{}) {
	ct := reflect.ValueOf(env).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("mapstructure")
		if tagVal != "" {
			val := os.Getenv(tagVal)
			if val == "" {
				continue
			}
			switch field.Type.Name() {
			case "string":
				ct.Field(i).SetString(val)
			case "int":
				ct.Field(i).SetInt(int64(utils.StringToInt(val)))
			case "bool":
				ct.Field(i).SetBool(utils.StringToBool(val))
			case "Duration":
				ct.Field(i).Set(reflect.ValueOf(utils.StringToTimeDuration(val)))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}
