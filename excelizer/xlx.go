package excelizer

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type Columns []string

func (c *Columns) String() string {
	return strings.Join(*c, ", ")
}

// Set добавляет значение в массив
func (c *Columns) Set(value string) error {
	*c = append(*c, value)
	return nil
}

type XLSX struct {
	file    *excelize.File
	headers Columns
}

func NewXLSX(headers Columns) *XLSX {
	return &XLSX{headers: headers}
}

func (x *XLSX) OpenXLSX(filename string) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("error opening xl file: %e", err)
	}
	x.file = f
	return nil
}

func (x *XLSX) ReadRow(sheet string, row int) ([]string, error) {
	var result []string
	cols, err := x.file.GetCols(sheet)
	if err != nil {
		return nil, err
	}
	for _, col := range cols {
		if len(col) >= row {
			result = append(result, col[row-1])
		} else {
			result = append(result, "")
		}
	}
	return result, nil
}

func (x *XLSX) WriteRow(sheet string, row int, data []interface{}) error {
	for col, value := range data {
		cell, err := excelize.CoordinatesToCellName(col+1, row)
		if err != nil {
			return err
		}
		x.file.SetCellValue(sheet, cell, value)
	}
	return nil
}

func (x *XLSX) WriteCell(sheet string, row, col int, data interface{}) error {
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}

	return x.file.SetCellValue(sheet, cell, data)
}

func (x *XLSX) Save(filename string) error {
	return x.file.SaveAs(filename)
}
