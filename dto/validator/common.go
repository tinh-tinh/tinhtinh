package validator

import (
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsAlpha(str string) bool {
	return IsRegexMatch(`^[a-zA-Z]+$`, str)
}

func IsAlphanumeric(str string) bool {
	return IsRegexMatch(`^[a-zA-Z0-9]+$`, str)
}

func IsEmail(str string) bool {
	return IsRegexMatch(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, str)
}

func IsStrongPassword(str string) bool {
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

func IsUUID(str string) bool {
	_, err := uuid.Parse(str)

	return err == nil
}

func IsObjectId(str string) bool {
	_, err := primitive.ObjectIDFromHex(str)

	return err == nil
}

func IsRegexMatch(pattern string, str string) bool {
	regex := regexp.MustCompile(pattern)
	match := regex.MatchString(str)

	return match
}

// Numeric
func IsInt(str string) bool {
	_, err := strconv.Atoi(str)

	return err == nil
}

func IsFloat(str string) bool {
	_, err := strconv.ParseFloat(str, 64)

	return err == nil
}

func IsNumber(str string) bool {
	return IsInt(str) || IsFloat(str)
}

// Date time
func IsDateString(str string) bool {
	_, err := time.Parse("2006-01-02", str)

	return err == nil
}

// Boolean
func IsBool(str string) bool {
	_, err := strconv.ParseBool(str)

	return err == nil
}
