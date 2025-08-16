package validator

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

func typeof(v any) string {
	return fmt.Sprintf("%T", v)
}

func IsAlpha(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsAlpha(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	if typeof(input) != "string" {
		return false
	}

	return IsRegexMatch(`^[a-zA-Z]+$`, input)
}

func IsAlphanumeric(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsAlphanumeric(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	if typeof(input) != "string" {
		return false
	}
	return IsRegexMatch(`^[a-zA-Z0-9]+$`, input)
}

func IsEmail(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsEmail(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	if typeof(input) != "string" {
		return false
	}
	return IsRegexMatch(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, input)
}

func IsStrongPassword(str any) bool {
	if typeof(str) != "string" {
		return false
	}
	if len(str.(string)) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range str.(string) {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func IsUUID(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsUUID(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	if typeof(input) != "string" {
		return false
	}

	uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
	return IsRegexMatch(uuidPattern, input)
}

func IsObjectId(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsObjectId(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	if typeof(input) != "string" {
		return false
	}

	objectIdPattern := `^[a-f0-9]{24}$`
	return IsRegexMatch(objectIdPattern, input)
}

func IsRegexMatch(pattern string, str any) bool {
	if typeof(str) != "string" {
		return false
	}
	regex := regexp.MustCompile(pattern)
	match := regex.MatchString(str.(string))

	return match
}

// Numeric
func IsInt(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsInt(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	typeInt := reflect.TypeOf(input)
	switch typeInt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		_, err := strconv.Atoi(input.(string))
		return err == nil
	default:
		log.Printf("%v is not be integer\n", typeInt.Kind())
		return false
	}
}

func IsFloat(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsFloat(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	typeFloat := reflect.TypeOf(input)
	switch typeFloat.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		_, err := strconv.ParseFloat(input.(string), 64)

		return err == nil
	default:
		log.Printf("%v is not be float\n", typeFloat.Kind())
		return false
	}
}

func IsNumber(str any) bool {
	return IsInt(str) || IsFloat(str)
}

// Date time
func IsDate(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsDate(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	switch v := input.(type) {
	case time.Time:
		return true
	case string:
		_, err := time.Parse("2006-01-02", input.(string))
		return err == nil
	default:
		log.Printf("%v is not be date\n", v)
		return false
	}
}

// Boolean
func IsBool(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		for i := range value.Len() {
			if !IsBool(value.Index(i).Interface()) { // Recursive
				return false
			}
		}
		return true
	}

	typeBool := reflect.TypeOf(input)
	switch typeBool.Kind() {
	case reflect.Bool:
		return true
	case reflect.String:
		_, err := strconv.ParseBool(input.(string))

		return err == nil
	default:
		return false
	}
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Pointer:
		return val.IsNil()
	case reflect.Array, reflect.Slice:
		return val.Len() == 0
	case reflect.Map:
		return val.Len() == 0
	default:
		return val.IsZero()
	}
}

func MinLength(input any, min int) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		return value.Len() >= min
	}

	if typeof(input) != "string" {
		return false
	}

	return len(input.(string)) >= min
}

func MaxLength(input any, max int) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	t := value.Type()

	if t.Kind() == reflect.Slice {
		return value.Len() <= max
	}

	if typeof(input) != "string" {
		return false
	}

	return len(input.(string)) <= max
}
