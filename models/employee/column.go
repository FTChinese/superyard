package employee

import (
	"fmt"
)

type Column int

const (
	ColumnStaffID Column = iota
	ColumnEmail
	ColumnUserName
)

var columnNames = [...]string{
	"staff_id",
	"email",
	"user_name",
}

var columnValues = map[string]Column{
	"id":    ColumnStaffID,
	"email": ColumnEmail,
	"name":  ColumnUserName,
}

func (c Column) String() string {
	if c >= ColumnStaffID && c <= ColumnUserName {
		return columnNames[c]
	}

	return ""
}

func ParseColumn(key string) (Column, error) {
	if c, ok := columnValues[key]; ok {
		return c, nil
	}

	return -1, fmt.Errorf("column for %s is not defined", key)
}
