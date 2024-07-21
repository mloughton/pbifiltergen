package dax

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	errEmptyInput  = errors.New("empty input")
	errFormat      = errors.New("incorrect format in input. expecting: table_name[column_name],type")
	errInvalidType = errors.New("invalid type in input. expecting : STRING, NUMBER, DATE")
	errBadChar     = errors.New("invalid characters in input")
	errINF         = errors.New("column/table name cannot begin with \"INF\"")
)

type Column struct {
	fullName   string
	tableName  string
	columnName string
	columnType string
}

func ParseInput(input string) (*[]Column, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil, errEmptyInput
	}
	lines := strings.Split(input, "\n")
	columns := []Column{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		column, err := createColumn(line)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return &columns, nil
}

func createColumn(raw string) (Column, error) {
	badChars := ".,;':/\\*|?&%$!+=()[]{}<>"
	validTypes := []string{"STRING", "NUMBER", "DATE"}
	var col Column
	tableColType := strings.Split(raw, ",")
	if len(tableColType) != 2 {
		return col, errFormat
	}
	colType := tableColType[1]
	if !slices.Contains(validTypes, colType) {
		return col, errInvalidType
	}
	fullCol := tableColType[0]
	tableCol := strings.Split(fullCol, "[")
	if len(tableCol) != 2 {
		return col, errFormat
	}

	table := []rune(tableCol[0])
	if string(table[:3]) == "INF" {
		return col, errINF
	}
	if table[0] == '\'' && table[len(table)-1] == '\'' {
		table = table[1 : len(table)-1]
		if strings.ContainsAny(string(table), badChars) {
			return col, errBadChar
		}
	} else {
		if strings.ContainsAny(string(table), badChars+" ") {
			return col, errBadChar
		}
	}

	column := []rune(tableCol[1])
	if string(column[:3]) == "INF" {
		return col, errINF
	}
	if column[len(column)-1] != ']' {
		return col, errFormat
	}
	column = column[:len(column)-1]
	if strings.ContainsAny(string(column), badChars) {
		return col, errBadChar
	}

	col.fullName = fullCol
	col.tableName = string(table)
	col.columnName = string(column)
	col.columnType = colType
	return col, nil
}

func GenerateDax(cols *[]Column) (string, error) {
	dax := ""
	daxLink := ""
	for i, col := range *cols {
		urlTable := strings.ReplaceAll(col.tableName, " ", "_x0020_")
		urlColumn := strings.ReplaceAll(col.columnName, " ", "_x0020_")
		daxLink += fmt.Sprintf("V%d,", i)
		switch col.columnType {
		case "STRING":
			dax += fmt.Sprintf(templateString, i, col.fullName, col.fullName, col.fullName, urlTable, urlColumn)
		case "DATE":
			dax += fmt.Sprintf(templateDate, i, col.fullName, col.fullName, urlTable, urlColumn, col.fullName, urlTable, urlColumn)
		case "NUMBER":
			dax += fmt.Sprintf(templateNum, i, col.fullName, col.fullName, col.fullName, col.fullName, urlTable, urlColumn, col.fullName, urlTable, urlColumn)
		default:
			return "", errInvalidType
		}

	}
	fullDax := fmt.Sprintf(templateEnd, dax, string([]rune(daxLink)[:len(daxLink)-1]))
	return fullDax, nil
}
