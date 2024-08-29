package validation

import (
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// String
func IsAlpha(str string) bool {
	return IsRegexMatch(`[a-zA-Z]+`, str)
}

func IsAlphanumeric(str string) bool {
	return IsRegexMatch(`[a-zA-Z0-9]+`, str)
}

func IsEmail(str string) bool {
	return IsRegexMatch(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, str)
}

func IsStrongPassword(str string) bool {
	return IsRegexMatch("^(?=.*[A-Z].*[A-Z])(?=.*[!@#$&*])(?=.*[0-9].*[0-9])(?=.*[a-z].*[a-z].*[a-z]).{8}$", str)
}

func IsUUID(str string) bool {
	_, err := uuid.FromString(str)

	return err == nil
}

func IsObjectId(str string) bool {
	_, err := primitive.ObjectIDFromHex(str)

	return err == nil
}

func IsRegexMatch(pattern string, str string) bool {
	match, _ := regexp.MatchString(pattern, str)

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
