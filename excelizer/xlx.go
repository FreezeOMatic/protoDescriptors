package excelizer

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type Columns struct {
	Flags     Flags
	ColumnMap map[string]string
	NamesMap  map[string]string
}

func NewColumns() *Columns {
	return &Columns{
		ColumnMap: make(map[string]string),
		NamesMap:  make(map[string]string),
	}
}

type Flags []string

func (c *Flags) String() string {
	return strings.Join(*c, ", ")
}

// Set добавляет значение в массив
func (c *Flags) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func (f *Columns) MapValuesToColumnsFromFlags() {
	for _, val := range f.Flags {
		colAndName := strings.Split(val, ":")
		if len(colAndName) != 2 {
			continue
		}
		col := colAndName[0]
		name := colAndName[1]

		f.ColumnMap[name] = col
	}
}

func (f *Columns) MapInternalNamesToFlagNames() {
	fmt.Println(f.Flags)

	f.NamesMap["parent.message"] = strings.Split(f.Flags[0], ":")[1]
	f.NamesMap["parent.enum.name"] = strings.Split(f.Flags[1], ":")[1]
	f.NamesMap["parent.enum.name.valueName"] = strings.Split(f.Flags[2], ":")[1]
	f.NamesMap["parent.enum.name.valueID"] = strings.Split(f.Flags[3], ":")[1]
	f.NamesMap["parent.enum.name.extention_info"] = strings.Split(f.Flags[4], ":")[1]
	f.NamesMap["parent.enum.name.extention_description"] = strings.Split(f.Flags[5], ":")[1]
	f.NamesMap["parent.enum.name.extention_measure"] = strings.Split(f.Flags[6], ":")[1]
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

func (x *XLSX) CreateXLSX() error {
	f := excelize.NewFile()
	x.file = f
	return nil
}
func (x *XLSX) CreateSheet(name string) error {
	i, err := x.file.NewSheet(name)
	if err != nil {
		return fmt.Errorf("error creating sheet: %e", err)
	}
	x.file.SetActiveSheet(i)
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

func (x *XLSX) WriteValuesMap(sheet string, values []map[string]interface{}) error {
	rowB := strings.Builder{}
	for name, col := range x.headers.ColumnMap {
		err := x.file.SetCellValue(sheet, fmt.Sprintf("%s1", col), name)
		if err != nil {
			return fmt.Errorf("error writing header for %s: %e", name, err)
		}
		rowB.WriteString(fmt.Sprintf("%s1", col) + ":" + name + " | ")
	}

	fmt.Println(rowB.String())

	for i, row := range values {
		rowS := strings.Builder{}
		for key, value := range row {
			cell := fmt.Sprintf("%s%d", x.headers.ColumnMap[x.headers.NamesMap[key]], i+2)
			err := x.file.SetCellValue(sheet, cell, value)
			if err != nil {
				return fmt.Errorf("error setting cell %d:%s to %s: %e", i, key, cell, err)
			}
			rowS.WriteString(cell + ":" + fmt.Sprintf("%v", value) + " | ")
		}
		fmt.Println(rowS.String())
	}

	return nil
}

func (x *XLSX) WriteRow(sheet string, row int, data []interface{}) error {
	for col, value := range data {
		cell, err := excelize.CoordinatesToCellName(col+1, row)
		if err != nil {
			return err
		}
		if err = x.file.SetCellValue(sheet, cell, value); err != nil {
			return fmt.Errorf("error writing row for %s: %e", cell, err)
		}
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
