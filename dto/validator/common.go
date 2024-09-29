package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func IsAlpha(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}
	return IsRegexMatch(`^[a-zA-Z]+$`, str)
}

func IsAlphanumeric(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}
	return IsRegexMatch(`^[a-zA-Z0-9]+$`, str)
}

func IsEmail(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}
	return IsRegexMatch(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, str)
}

func IsStrongPassword(str interface{}) bool {
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

func IsUUID(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}

	uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
	return IsRegexMatch(uuidPattern, str)
}

func IsObjectId(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}

	objectIdPattern := `/^[a-f\d]{24}$/i`
	return IsRegexMatch(objectIdPattern, str)
}

func IsRegexMatch(pattern string, str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}
	regex := regexp.MustCompile(pattern)
	match := regex.MatchString(str.(string))

	return match
}

// Numeric
func IsInt(str interface{}) bool {
	switch v := str.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case string:
		_, err := strconv.Atoi(str.(string))
		return err == nil
	default:
		fmt.Println(v)
		return false
	}
}

func IsFloat(str interface{}) bool {
	switch v := str.(type) {
	case float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(str.(string), 64)

		return err == nil
	default:
		fmt.Println(v)
		return false
	}
}

func IsNumber(str interface{}) bool {
	return IsInt(str) || IsFloat(str)
}

// Date time
func IsDateString(str interface{}) bool {
	if typeof(str) != "string" {
		return false
	}
	_, err := time.Parse("2006-01-02", str.(string))

	return err == nil
}

// Boolean
func IsBool(str interface{}) bool {
	switch typeof(str) {
	case "bool":
		return str.(bool)
	case "string":
		_, err := strconv.ParseBool(str.(string))

		return err == nil
	default:
		return false
	}
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
	default:
		return val == nil
	}
}
