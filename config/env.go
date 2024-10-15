package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"gopkg.in/yaml.v3"
)

func New[E any](path string) (*E, error) {
	if path == "" {
		path = ".env"
	}
	if strings.Contains(path, ".env") {
		err := godotenv.Load(path)
		if err != nil {
			return nil, err
		}

		var env E
		Scan(&env)
		return &env, nil
	} else if strings.Contains(path, ".yml") || strings.Contains(path, ".yaml") {
		var e E
		fmt.Println("here")

		file, err := os.ReadFile(path)
		if err != nil {
			log.Printf("yamlFile get error: %v\n", err)
			os.Exit(1)
		}

		err = yaml.Unmarshal(file, &e)
		if err != nil {
			log.Printf("get error: %v\n", err)
			os.Exit(1)
		}
		return &e, nil
	} else {
		log.Printf("not supported type: %v\n", path)
		return nil, errors.New("not support")
	}
}

func GetRaw(key string) string {
	return os.Getenv(key)
}

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
