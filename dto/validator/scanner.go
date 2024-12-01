package validator

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/dto/transform"
)

func Scanner(val interface{}) error {
	var errMsg []string
	if val == nil {
		panic(fmt.Sprintf("%v should be not nil", val))
	}
	if reflect.TypeOf(val).Kind() == reflect.Struct {
		panic(fmt.Sprintf("%v should be a value not struct", val))
	}

	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("validate")
		if tagVal == "" {
			continue
		}

		validators := strings.Split(tagVal, ",")
		value := ct.Field(i).Interface()

		defaultVal := field.Tag.Get("default")
		if defaultVal != "" && reflect.ValueOf(value).IsZero() {
			value = defaultVal
		}

		required := slices.IndexFunc(validators, func(v string) bool { return v == "required" })
		if required == -1 && value == "" {
			continue
		}
		for _, validate := range validators {
			switch validate {
			case "required":
				if IsNil(value) {
					errMsg = append(errMsg, field.Name+" is required")
				}
			case "isAlpha":
				if !IsAlpha(value) {
					errMsg = append(errMsg, field.Name+" is not a valid alpha")
				}
			case "isAlphaNumeric":
				if !IsAlphanumeric(value) {
					errMsg = append(errMsg, field.Name+" is not a valid alpha numeric")
				}
			case "isEmail":
				if !IsEmail(value) {
					errMsg = append(errMsg, field.Name+" is not a valid email")
				}
			case "isStrongPassword":
				if !IsStrongPassword(value) {
					errMsg = append(errMsg, field.Name+" is not a valid strong password")
				}
			case "isUUID":
				if !IsUUID(value) {
					errMsg = append(errMsg, field.Name+" is not a valid UUID")
				}
			case "isObjectId":
				if !IsObjectId(value) {
					errMsg = append(errMsg, field.Name+" is not a valid ObjectID")
				}
			case "isInt":
				if !IsInt(value) {
					errMsg = append(errMsg, field.Name+" is not a valid int")
				} else {
					ct.Field(i).Set(reflect.ValueOf(transform.ToInt(value)))
				}
			case "isFloat":
				if !IsFloat(value) {
					errMsg = append(errMsg, field.Name+" is not a valid float")
				} else {
					ct.Field(i).Set(reflect.ValueOf(transform.ToFloat(value)))
				}
			case "isNumber":
				if !IsNumber(value) {
					errMsg = append(errMsg, field.Name+" is not a valid number")
				} else {
					if IsInt(value) {
						ct.Field(i).Set(reflect.ValueOf(transform.ToInt(value)))
					} else {
						ct.Field(i).Set(reflect.ValueOf(transform.ToFloat(value)))
					}
				}
			case "isDateString":
				if !IsDateString(value) {
					errMsg = append(errMsg, field.Name+" is not a valid date time")
				} else {
					ct.Field(i).Set(reflect.ValueOf(transform.ToDate(value)))
				}
			case "isBool":
				if !IsBool(value) {
					errMsg = append(errMsg, field.Name+" is not a valid bool")
				} else {
					ct.Field(i).Set(reflect.ValueOf(transform.ToBool(value)))
				}
			case "nested":
				if field.Type.Kind() == reflect.Pointer {
					err := Scanner(ct.Field(i).Interface())
					if err != nil {
						errMsg = append(errMsg, err.Error())
					}
				} else if field.Type.Kind() == reflect.Slice {
					arrVal := reflect.ValueOf(ct.Field(i).Interface())
					if arrVal.IsValid() {
						for i := 0; i < arrVal.Len(); i++ {
							item := arrVal.Index(i)
							err := Scanner(item.Interface())
							if err != nil {
								errMsg = append(errMsg, err.Error())
							}
						}
					}
				}
			}
		}
	}

	if len(errMsg) == 0 {
		return nil
	}

	err := errors.New(errMsg[0])
	for i := 1; i < len(errMsg); i++ {
		err = errors.Join(err, errors.New(errMsg[i]))
	}

	return err
}
