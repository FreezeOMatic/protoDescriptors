package gen

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"regexp"
	"rpcdescriptors/excelizer"
	"strconv"
)

type Convertor struct {
	FDS       descriptorpb.FileDescriptorSet
	Excelizer *excelizer.XLSX
}

func init() {

}

func New(excelizer *excelizer.XLSX) *Convertor {
	return &Convertor{Excelizer: excelizer}
}

func (c *Convertor) OpenDescriptor(filename string) error {
	fdsBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("open descriptor: %v", err)
	}

	if err = proto.Unmarshal(fdsBytes, &c.FDS); err != nil {
		return fmt.Errorf("unmarshal descriptor: %v", err)
	}

	return nil
}

func (c *Convertor) CollectEnumValues() []map[string]interface{} {

	typesMap := make([]map[string]interface{}, 0)

	for _, file := range c.FDS.GetFile() {
		parentFile := file.GetName()
		fmt.Printf("Converting file %s\n", parentFile)

		if len(file.GetEnumType()) > 0 {
			typesMap = append(typesMap, c.ProcessEnum(file)...)
		}

		if len(file.GetMessageType()) > 0 {
			typesMap = append(typesMap, c.ProcessMessageWithEnum(file)...)
		}
	}

	return typesMap
}

func (c *Convertor) ProcessEnum(file *descriptorpb.FileDescriptorProto) []map[string]interface{} {
	extensionRef := make(map[int]string)
	for _, v := range file.GetExtension() {
		extensionRef[int(v.GetNumber())] = v.GetName()
	}

	typesMap := make([]map[string]interface{}, 0)

	for _, enum := range file.GetEnumType() {
		parentMessage := ""

		parentEnumName := enum.GetName()

		for _, enumValue := range enum.GetValue() {
			referenceMap := make(map[string]interface{})

			referenceMap["parent.message"] = parentMessage
			referenceMap["parent.enum.name"] = parentEnumName

			referenceMap["parent.enum.name.valueName"] = enumValue.GetName()
			referenceMap["parent.enum.name.valueID"] = enumValue.GetNumber()

			re := regexp.MustCompile(`(\d+):"([^"]+)"`)
			options := re.FindAllStringSubmatch(enumValue.GetOptions().String(), -1)
			for _, option := range options {
				optionCode, _ := strconv.Atoi(option[1])
				optionValue := option[2]

				referenceMap[fmt.Sprintf("parent.enum.name.extention_%s", extensionRef[optionCode])] = optionValue
			}
			typesMap = append(typesMap, referenceMap)
		}
	}
	return typesMap
}

func (c *Convertor) ProcessMessageWithEnum(file *descriptorpb.FileDescriptorProto) []map[string]interface{} {
	extensionRef := make(map[int]string)
	for _, v := range file.GetExtension() {
		extensionRef[int(v.GetNumber())] = v.GetName()
	}

	typesMap := make([]map[string]interface{}, 0)

	for _, msg := range file.GetMessageType() {
		for _, enum := range msg.GetEnumType() {
			parentMessage := msg.GetName()

			parentEnumName := enum.GetName()

			for _, enumValue := range enum.GetValue() {
				referenceMap := make(map[string]interface{})

				referenceMap["parent.message"] = parentMessage
				referenceMap["parent.enum.name"] = parentEnumName

				referenceMap["parent.enum.name.valueName"] = enumValue.GetName()
				referenceMap["parent.enum.name.valueID"] = enumValue.GetNumber()

				re := regexp.MustCompile(`(\d+):"([^"]+)"`)
				options := re.FindAllStringSubmatch(enumValue.GetOptions().String(), -1)
				for _, option := range options {
					optionCode, _ := strconv.Atoi(option[1])
					optionValue := option[2]

					referenceMap[fmt.Sprintf("parent.enum.name.extention_%s", extensionRef[optionCode])] = optionValue
				}
				typesMap = append(typesMap, referenceMap)
			}
		}
	}
	return typesMap
}
