package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"rpcdescriptors/excelizer"
	"rpcdescriptors/gen"
	_ "rpcdescriptors/gen"
)

func main() {
	var headers excelizer.Columns
	flag.Var(&headers, "col", "Добавьте строку в массив (можно использовать несколько раз)")
	flag.Parse()

	excel := excelizer.NewXLSX(headers)
	convertor := gen.New(excel)
	err := convertor.OpenDescriptor("gen/descriptor.pb")
	if err != nil {
		panic(err)
	}

	m := convertor.CollectEnumValues()

	v, _ := json.MarshalIndent(m, "", " ")
	fmt.Println(string(v))
}
