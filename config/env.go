package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/dto/transform"
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
				ct.Field(i).SetInt(transform.StringToInt64(val))
			case "bool":
				ct.Field(i).SetBool(transform.StringToBool(val))
			case "Duration":
				ct.Field(i).Set(reflect.ValueOf(transform.StringToTimeDuration(val)))
			default:
				fmt.Println(field.Type.Name())
			}
		}
	}
}
