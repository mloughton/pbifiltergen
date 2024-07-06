package dax

import (
	"errors"
	"strings"
)

type Column string

func ParseInput(input string) (*[]Column, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil, errors.New("empty input")
	}
	lines := strings.Split(input, "\n")
	columns := []Column{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		column, err := validateRawColumn(line)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return &columns, nil
}

func validateRawColumn(raw string) (Column, error) {
	// table_name[column_name]
	length := len(raw)
	if []rune(raw)[length-1] != ']' {
		return "", errors.New("malformed column: last bracket")
	}
	tableCol := strings.Split(string([]rune(raw)[:length-1]), "[")
	if len(tableCol) != 2 {
		return "", errors.New("malformed column: first bracket")
	}

	table := []rune(tableCol[0])
	badChars := ".,;':/\\*|?&%$!+=()[]{} <>"
	if table[0] == '\'' && table[len(table)-1] == '\'' {
		table = table[1 : len(table)-1]
		badChars = ".,;':/\\*|?&%$!+=()[]{}<>"
	}
	if strings.ContainsAny(string(table), badChars) {
		return "", errors.New("malformed column: table contains bad char")
	}
	column := tableCol[1]
	badChars = ".,;':/\\*|?&%$!+=()[]{}<>"
	if strings.ContainsAny(column, badChars) {
		return "", errors.New("malformed column: column contains bad char")
	}
	return Column(raw), nil
}
