package validator

import (
	"log"
	"reflect"
	"regexp"
	"unicode"
)

func IsAlpha(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isAlpha(value)
}

func isAlpha(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isAlpha(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	if str == "" {
		return false
	}
	for _, char := range str {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func IsAlphanumeric(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isAlphanumeric(value)
}

func isAlphanumeric(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isAlphanumeric(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	if str == "" {
		return false
	}
	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}

func IsEmail(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isEmail(value)
}

func isEmail(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isEmail(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	matched, _ := regexp.MatchString(emailPattern, str)
	return matched
}

func IsStrongPassword(str any) bool {
	if str == nil {
		return false
	}

	value := reflect.ValueOf(str)
	return isStrongPassword(value)
}

func isStrongPassword(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isStrongPassword(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	if len(str) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range str {
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
	return isUUID(value)
}

func isUUID(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isUUID(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	matched, _ := regexp.MatchString(uuidPattern, str)
	return matched
}

func IsObjectId(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isObjectId(value)
}

func isObjectId(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isObjectId(value.Index(i)) {
				return false
			}
		}
		return true
	}

	if value.Kind() != reflect.String {
		return false
	}

	str := value.String()
	matched, _ := regexp.MatchString(objectIdPattern, str)
	return matched
}

func IsRegexMatch(pattern string, str any) bool {
	if typeof(str) != "string" {
		return false
	}

	assertStr := str.(string)
	if len(assertStr) > MiB {
		log.Println("String too long for regex match")
		return false
	}
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	match := regex.MatchString(assertStr)

	return match
}
