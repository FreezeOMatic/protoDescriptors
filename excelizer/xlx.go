package excelizer

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type Columns struct {
	Flags     Flags
	ColumnMap map[string]ColumnParams
	NamesMap  map[string]string
	Style     int
}

type ColumnParams struct {
	Width int    // Ширина колонки
	Order string // Буква колонки
}

func NewColumns() *Columns {
	return &Columns{
		ColumnMap: make(map[string]ColumnParams),
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
	if len(f.Flags) == 0 {
		f.ColumnMap["Field A"] = ColumnParams{
			Width: 10,
			Order: "A",
		}
		f.ColumnMap["Field B"] = ColumnParams{
			Width: 45,
			Order: "B",
		}
		f.ColumnMap["Field C"] = ColumnParams{
			Width: 40,
			Order: "C",
		}
		f.ColumnMap["Field D"] = ColumnParams{
			Width: 20,
			Order: "D",
		}
		f.ColumnMap["Field E"] = ColumnParams{
			Width: 60,
			Order: "E",
		}
		f.ColumnMap["Field F"] = ColumnParams{
			Width: 40,
			Order: "F",
		}
		f.ColumnMap["Field G"] = ColumnParams{
			Width: 40,
			Order: "G",
		}
		f.ColumnMap["Field H"] = ColumnParams{
			Width: 15,
			Order: "H",
		}
		f.ColumnMap["Field I"] = ColumnParams{
			Width: 15,
			Order: "I",
		}
		f.ColumnMap["Field J"] = ColumnParams{
			Width: 12,
			Order: "J",
		}
		f.ColumnMap["Field K"] = ColumnParams{
			Width: 12,
			Order: "K",
		}
		f.ColumnMap["Field L"] = ColumnParams{
			Width: 25,
			Order: "L",
		}

		return
	}

	for _, val := range f.Flags {
		colAndName := strings.Split(val, ":")
		if len(colAndName) != 2 {
			continue
		}
		col := colAndName[0]
		name := colAndName[1]

		f.ColumnMap[name] = ColumnParams{
			Width: 20,
			Order: col,
		}
	}
}

func (f *Columns) MapInternalNamesToFlagNames() {
	fmt.Println(f.Flags)
	// map columns with internal names/ Columns in right
	f.NamesMap["valueName"] = "Value"
}

type XLSX struct {
	file    *excelize.File
	headers Columns
	rows    RowsParams
}

type RowsParams struct {
	style int
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

	columnsStyle, err := x.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#156082"}, // Синий цвет фона
		},
		Font: &excelize.Font{
			Bold:   true,           // Жирный шрифт
			Color:  "#FFFFFF",      // Чёрный цвет текста
			Family: "Aptos Narrow", // Заголовочный шрифт
			Size:   11,             // Увеличенный размер шрифта
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center", // Горизонтальное выравнивание по центру
			Vertical:   "center", // Вертикальное выравнивание по центру
		},
	})
	if err != nil {
		return fmt.Errorf("error creating style: %e", err)
	}
	x.headers.Style = columnsStyle

	rowsStyle, err := x.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#CCECFF"}, // Голубой цвет фона
		},
		Font: &excelize.Font{
			Color:  "#000000",
			Family: "Aptos Narrow",
			Size:   11,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating style: %e", err)
	}
	x.rows.style = rowsStyle

	return nil
}
func (x *XLSX) CreateSheet(name string) error {
	i, err := x.file.NewSheet(name)
	if err != nil {
		return fmt.Errorf("error creating sheet: %e", err)
	}
	x.file.SetActiveSheet(i)

	return x.file.SetPanes(name, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      1,
		YSplit:      1,
		TopLeftCell: "B2",
		ActivePane:  "bottomRight",
	})
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
		err := x.file.SetColWidth(sheet, col.Order, col.Order, float64(col.Width))
		if err != nil {
			fmt.Println("Failed to set column width:", err)
		}

		err = x.file.SetCellStyle(sheet, fmt.Sprintf("%s1", col.Order), "A1", x.headers.Style)
		if err != nil {
			return fmt.Errorf("error cetting cell style: %e", err)
		}

		err = x.file.SetCellValue(sheet, fmt.Sprintf("%s1", col.Order), name)
		if err != nil {
			return fmt.Errorf("error writing header for %s: %e", name, err)
		}
		rowB.WriteString(fmt.Sprintf("%s1", col.Order) + ":" + name + " | ")
	}

	fmt.Println(rowB.String())

	for i, row := range values {
		rowS := strings.Builder{}
		for key, value := range row {
			cell := fmt.Sprintf("%s%d", x.headers.ColumnMap[x.headers.NamesMap[key]].Order, i+2)

			pValue := value
			if v, ok := value.(string); ok {
				pValue = strings.ReplaceAll(v, "|", "\n")
			}

			err := x.file.SetCellValue(sheet, cell, pValue)
			if err != nil {
				return fmt.Errorf("error setting cell %d:%s to %s: %e", i, key, cell, err)
			}
			err = x.file.SetCellStyle(sheet, cell, cell, x.rows.style)
			if err != nil {
				return fmt.Errorf("error cetting cell style: %e", err)
			}

			rowS.WriteString(cell + ":" + fmt.Sprintf("%v", pValue) + " | ")
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
