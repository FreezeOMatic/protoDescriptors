package utils

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func GetMessageDescriptor(msg proto.Message) protoreflect.MessageDescriptor {
	return msg.ProtoReflect().Descriptor()
}

func GetFieldAnnotations(desc protoreflect.FieldDescriptor) (semanticCode, valueType string) {
	//opts := desc.Options().(*descriptorpb.FieldOptions)

	//if proto.HasExtension(opts, test.E_FieldCode) {
	//	semanticCode = proto.GetExtension(opts, test.E_FieldCode).(string)
	//}
	//
	//if proto.HasExtension(opts, test.E_BitPosition) {
	//	valueType = proto.GetExtension(opts, test.E_BitPosition).(string)
	//}
	//
	return semanticCode, valueType
}

func GetReliedDescriptor(enumValue protoreflect.EnumValueDescriptor) *descriptorpb.DescriptorProto {
	//options := enumValue.Options()
	//
	//messageDescriptor := proto.GetExtension(options, test.E_Descriptor).(*descriptorpb.DescriptorProto)
	//if messageDescriptor == nil {
	//	fmt.Println("Message descriptor not found for enum value:", enumValue.Name())
	//	return nil
	//}
	//
	//fmt.Printf("Enum value: %s\n", enumValue.Name())
	//fmt.Printf("Message name: %s\n", messageDescriptor.GetName())
	//fmt.Println("Fields:")
	//for _, field := range messageDescriptor.GetField() {
	//	fmt.Printf("  - %s (type: %v, number: %d)\n", field.GetName(), field.GetType(), field.GetNumber())
	//}
	//fmt.Println("---------------------------")
	//messageAbbreviation := proto.GetExtension(options, test.E_Abbr).(string)
	//fmt.Printf("Additional info: %s\n", messageAbbreviation)

	return nil
}

func GetConcreteMessageFromDescriptor(desc *descriptorpb.DescriptorProto) (proto.Message, error) {
	// Получаем полное имя сообщения (например, "example.Undefined")
	fullName := protoreflect.FullName(desc.GetName())

	// Ищем тип сообщения в глобальном реестре типов
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to find message type: %v", err)
	}

	// Создаем новый экземпляр сообщения
	concreteMsg := messageType.New().Interface()

	return concreteMsg, nil
}

func GetMessageFromDescriptor(desc *descriptorpb.DescriptorProto) (*dynamicpb.Message, error) {
	fullName := protoreflect.FullName(desc.GetName())
	//protoregistry.GlobalTypes.RangeMessages(func(messageType protoreflect.MessageType) bool {
	//	fmt.Println(messageType.Descriptor().FullName())
	//	return true
	//})

	messageType, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to find message type: %v", err)
	}

	message := dynamicpb.NewMessage(messageType.Descriptor())

	return message, nil
}
