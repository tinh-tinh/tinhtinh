package core

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"
)

// ParseCsv takes an interface{} and a slice of string headers and returns a
// [][]string that can be used with the encoding/csv package. The function
// expects the given interface{} to be a slice of structs, and it will extract
// the values of the structs and append them to the given headers.
//
// The function returns a [][]string where the first element is the given
// headers, and the rest of the elements are the extracted values of the given
// structs.
//
// If the given interface{} is not a slice, the function will return a [][]string
// with only the given headers.
func ParseCsv(data interface{}, headers []string) [][]string {
	parseBody := [][]string{}

	if headers == nil {
		return parseBody
	}
	parseBody = append(parseBody, headers)

	if data == nil || reflect.TypeOf(data).Kind() != reflect.Slice {
		return parseBody
	}

	arrVal := reflect.ValueOf(data)
	if arrVal.IsValid() {
		for i := 0; i < arrVal.Len(); i++ {
			item := arrVal.Index(i).Interface()
			ctItem := reflect.ValueOf(item).Elem()
			row := []string{}

			for i := 0; i < ctItem.NumField(); i++ {
				value := ctItem.Field(i).Interface()
				valStr := fmt.Sprintf("%v", value)

				row = append(row, valStr)
			}

			parseBody = append(parseBody, row)
		}
	}

	return parseBody
}

// checkFileName takes a filename as input and returns the same filename with
// the ".csv" extension added if it was not already present. It is used to
// ensure that a filename passed to ExportCSV is a valid CSV file name.
func checkFileName(name string) string {
	isCsv := strings.Contains(name, ".csv")

	if !isCsv {
		name = name + ".csv"
	}

	return name
}

// ExportCSV takes a filename and a body of data, and writes that data to the
// HTTP response as a CSV file. The filename is used to set the
// Content-Disposition header, and must have a ".csv" extension. The body
// should be a slice of slices of strings, where each inner slice is a row of
// data, and each string is a column in that row. The data is written to the
// response in the same order as it appears in the body.
func (ctx *Ctx) ExportCSV(name string, body [][]string) error {
	name = checkFileName(name)

	ctx.Res().Header().Set("Content-Type", "text/csv")
	ctx.Res().Header().Set("Content-Disposition", "attachment; filename="+name)

	writer := csv.NewWriter(ctx.Res())
	defer writer.Flush()

	for _, row := range body {
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}
