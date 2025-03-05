package main

import (
	"flag"
	"rpcdescriptors/excelizer"
	"rpcdescriptors/gen"
	_ "rpcdescriptors/gen"
)

func main() {

	cols := excelizer.NewColumns()

	flag.Var(&cols.Flags, "col", "Добавьте строку в массив (можно использовать несколько раз)")
	flag.Parse()

	cols.MapValuesToColumnsFromFlags()
	cols.MapInternalNamesToFlagNames()

	excel := excelizer.NewXLSX(*cols)
	err := excel.CreateXLSX()
	if err != nil {
		panic(err)
	}
	err = excel.CreateSheet("vitrine_test")
	if err != nil {
		panic(err)
	}

	convertor := gen.New(excel)
	err = convertor.OpenDescriptor("gen/descriptor.pb")
	if err != nil {
		panic(err)
	}

	enums := convertor.CollectEnumValues()

	err = excel.WriteValuesMap("vitrine_test", enums)
	if err != nil {
		panic(err)
	}

	err = excel.Save("test.xlsx")
	if err != nil {
		panic(err)
	}
}
