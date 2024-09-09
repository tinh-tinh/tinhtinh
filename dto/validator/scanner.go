package validator

import (
	"errors"
	"reflect"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/dto/transform"
)

func Scanner(val interface{}, trans bool) error {
	var errMsg []string

	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tagVal := field.Tag.Get("validate")
		if tagVal == "" {
			continue
		}

		validators := strings.Split(tagVal, ",")
		value := ct.Field(i).Interface()

		required := slices.IndexFunc(validators, func(v string) bool { return v == "required" })
		if required == -1 && value == "" {
			continue
		}
		for _, validate := range validators {
			switch validate {
			case "required":
				if value == nil {
					errMsg = append(errMsg, field.Name+" is required")
				}
			case "isAlpha":
				if !IsAlpha(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid alpha")
				}
			case "isAlphaNumeric":
				if !IsAlphanumeric(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid alpha numeric")
				}
			case "isEmail":
				if !IsEmail(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid email")
				}
			case "isStrongPassword":
				if !IsStrongPassword(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid strong password")
				}
			case "isUUID":
				if !IsUUID(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid UUID")
				}
			case "isObjectId":
				if !IsObjectId(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid ObjectID")
				} else if trans {
					ct.Field(i).Set(reflect.ValueOf(transform.StringToObjectID(value.(string))))
				}
			case "isInt":
				if !IsInt(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid int")
				} else if trans {
					ct.Field(i).Set(reflect.ValueOf(transform.StringToInt(value.(string))))
				}
			case "isFloat":
				if !IsFloat(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid float")
				} else if trans {
					ct.Field(i).Set(reflect.ValueOf(transform.StringToFloat(value.(string))))
				}
			case "isNumber":
				if !IsNumber(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid number")
				} else if trans {
					if IsInt(value.(string)) {
						ct.Field(i).Set(reflect.ValueOf(transform.StringToInt(value.(string))))
					} else {
						ct.Field(i).Set(reflect.ValueOf(transform.StringToFloat(value.(string))))
					}
				}
			case "isDateString":
				if !IsDateString(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid date time")
				} else if trans {
					ct.Field(i).Set(reflect.ValueOf(transform.StringToDate(value.(string))))
				}
			case "isBool":
				if !IsBool(value.(string)) {
					errMsg = append(errMsg, field.Name+" is not a valid bool")
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
