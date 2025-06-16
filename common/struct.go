package common

import (
	"reflect"
	"runtime"
)

func GetStructName(str interface{}) string {
	name := ""
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return name
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func IsPrimitiveType(fieldKind reflect.Kind) bool {
	if fieldKind == reflect.Bool ||
		fieldKind == reflect.Int ||
		fieldKind == reflect.Int8 ||
		fieldKind == reflect.Int16 ||
		fieldKind == reflect.Int32 ||
		fieldKind == reflect.Int64 ||
		fieldKind == reflect.Uint ||
		fieldKind == reflect.Uint8 ||
		fieldKind == reflect.Uint16 ||
		fieldKind == reflect.Uint32 ||
		fieldKind == reflect.Uint64 ||
		fieldKind == reflect.Float32 ||
		fieldKind == reflect.Float64 ||
		fieldKind == reflect.String {
		return true
	}
	return false
}

func PartialStruct[T any](input T) any {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	fields := make([]reflect.StructField, inputType.NumField())
	for i := range inputType.NumField() {
		field := inputType.Field(i)
		if IsPrimitiveType(field.Type.Kind()) {
			fields[i] = reflect.StructField{
				Name: field.Name,
				Type: reflect.PointerTo(field.Type),
				Tag:  field.Tag,
			}
		} else {
			fields[i] = field
		}
	}

	newType := reflect.StructOf(fields)
	outputValue := reflect.New(newType).Elem()

	for i := range inputValue.NumField() {
		field := inputValue.Field(i)
		if IsPrimitiveType(field.Type().Kind()) {
			ptr := reflect.New(field.Type())
			ptr.Elem().Set(field)
			outputValue.Field(i).Set(ptr)
		} else {
			outputValue.Field(i).Set(field)
		}
	}

	return outputValue.Interface()
}

func PickStruct[T any](input T, fields []string) any {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	// Create a map of requested field names for easy lookup
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	// Collect only the requested fields
	var selectedFields []reflect.StructField
	for i := range inputType.NumField() {
		field := inputType.Field(i)
		if fieldMap[field.Name] {
			selectedFields = append(selectedFields, field)
		}
	}

	// Create new struct type with selected fields
	newType := reflect.StructOf(selectedFields)
	outputValue := reflect.New(newType).Elem()

	// Copy values for selected fields
	fieldIndex := 0
	for i := range inputType.NumField() {
		field := inputType.Field(i)
		if fieldMap[field.Name] {
			outputValue.Field(fieldIndex).Set(inputValue.Field(i))
			fieldIndex++
		}
	}

	return outputValue.Interface()
}

func OmitStruct[T any](input T, fields []string) any {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	// Create a map of fields to exclude for easy lookup
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	// Collect fields that are not in the exclude list
	var selectedFields []reflect.StructField
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if !fieldMap[field.Name] {
			selectedFields = append(selectedFields, field)
		}
	}

	// Create new struct type with remaining fields
	newType := reflect.StructOf(selectedFields)
	outputValue := reflect.New(newType).Elem()

	// Copy values for remaining fields
	fieldIndex := 0
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if !fieldMap[field.Name] {
			outputValue.Field(fieldIndex).Set(inputValue.Field(i))
			fieldIndex++
		}
	}

	return outputValue.Interface()
}

func MergeStruct[T any](input ...T) T {
	var result T
	resultVal := reflect.ValueOf(&result).Elem()

	for _, item := range input {
		inputVal := reflect.ValueOf(item)
		for i := range inputVal.NumField() {
			field := inputVal.Field(i)
			resultField := resultVal.Field(i)

			// Set field if result's field is zero and input's field is not zero
			if resultField.IsZero() && !field.IsZero() {
				resultField.Set(field)
			}
		}
	}

	return result
}
