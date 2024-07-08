package dax

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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
		return nil, errors.New("empty input")
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
	badChars := ".,;':/\\*|?&%$!+=()[]{} <>"
	validTypes := []string{"STRING", "INTEGER", "DATETIME"}
	formatError := errors.New("incorrect format in input. expecting: table_name[column_name],type")
	invalidTypeError := errors.New("invalid type in input. expecting : STRING, INTEGER, DATETIME")
	badCharError := errors.New("invalid characters in input")
	var col Column
	tableColType := strings.Split(raw, ",")
	if len(tableColType) != 2 {
		return col, formatError
	}
	colType := tableColType[1]
	if !slices.Contains(validTypes, colType) {
		return col, invalidTypeError
	}
	fullCol := tableColType[0]
	tableCol := strings.Split(fullCol, "[")
	if len(tableCol) != 2 {
		return col, formatError
	}

	table := []rune(tableCol[0])
	if table[0] == '\'' && table[len(table)-1] == '\'' {
		table = table[1 : len(table)-1]
		if strings.ContainsAny(string(table), strings.ReplaceAll(badChars, " ", "")) {
			return col, badCharError
		}
	} else {
		if strings.ContainsAny(string(table), badChars) {
			return col, badCharError
		}
	}
	column := []rune(tableCol[1])

	if column[len(column)-1] != ']' {
		return col, formatError
	}
	column = column[:len(column)-1]
	if strings.ContainsAny(string(column), badChars) {
		return col, badCharError
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
		daxLink += fmt.Sprintf("V%d,", i)
		dax += fmt.Sprintf(
			`VAR V%d =
	IF(
		ISFILTERED( %s ) = TRUE(),
		VAR vals = VALUES( %s )
		VAR cnt = COUNTROWS(vals)
		VAR valsStr = CONCATENATEX(vals, %s, "%%27,%%20%%27")
		VAR fullStr = "%s/%s%%20" & IF(cnt = 1, "eq%%20%%27" & vals & "%%27", "in%%20(%%27" & valsStr & "%%27)")
		VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
		RETURN
			fullStrBlank,
		BLANK()
	)
`, i, col.fullName, col.fullName, col.fullName, strings.ReplaceAll(col.tableName, " ", "x0020"), strings.ReplaceAll(col.columnName, " ", "x0020"))
	}
	fullDax := fmt.Sprintf(
		`%s
VAR link = CONCATENATEX(FILTER({%s}, [Value] <> BLANK()), [Value], "%%20and%%20")
VAR fullLink = IF(link = BLANK(), BLANK(), [Dashboard Link] & "?rs:embed=true&filter=" & link)
RETURN
fullLink
`, dax, string([]rune(daxLink)[:len(daxLink)-1]))
	return fullDax, nil
}
