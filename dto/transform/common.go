package transform

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObjectID(str string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(str)
	return id
}

func ObjectIDToString(id primitive.ObjectID) string {
	return id.Hex()
}

func StringToBool(str string) bool {
	val, _ := strconv.ParseBool(str)
	return val
}

func StringToInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func StringToFloat(str string) float64 {
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

func StringToDate(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}
