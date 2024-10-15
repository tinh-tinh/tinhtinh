package core

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"
)

func ParseCsv(data interface{}, headers []string) [][]string {
	parseBody := [][]string{}

	parseBody = append(parseBody, headers)

	if reflect.TypeOf(data).Kind() != reflect.Slice {
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

func checkFileName(name string) string {
	isCsv := strings.Contains(name, ".csv")

	if !isCsv {
		name = name + ".csv"
	}

	return name
}

func (ctx *Ctx) ExportCSV(name string, body [][]string) error {
	name = checkFileName(name)

	ctx.Res().Header().Set("Content-Type", "text/csv")
	ctx.Res().Header().Set("Content-Disposition", "attachment; filename="+name)

	writer := csv.NewWriter(ctx.Res())
	defer writer.Flush()

	for _, row := range body {
		writer.Write(row)
	}
	return nil
}
