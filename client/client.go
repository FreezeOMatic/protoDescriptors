package client

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"rpcdescriptors/gen/test"
)

func ReceiveAndParse(data []byte) (*dynamicpb.Message, error) {
	// 1. Распаковываем контейнер
	var bundle test.PackedBundle
	if err := proto.Unmarshal(data, &bundle); err != nil {
		return nil, err
	}

	descriptorName := bundle.GetDescriptor_().Name
	fmt.Println("desc Name", *descriptorName)

	// 2. Создаем временный FileDescriptorProto
	fileDesc := &descriptorpb.FileDescriptorProto{
		Name:   proto.String("dynamic.proto"),
		Syntax: proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			bundle.GetDescriptor_(),
		},
	}

	// 3. Регистрируем тип
	fd, err := protodesc.NewFile(fileDesc, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("descriptor: ", fd.Messages().Get(0).FullName())

	// 4. Находим дескриптор
	desc := fd.Messages().ByName(protoreflect.Name(*descriptorName))

	// 5. Создаем и парсим сообщение
	msg := dynamicpb.NewMessage(desc)
	if err := proto.Unmarshal(bundle.GetValue(), msg); err != nil {
		return nil, err
	}

	return msg, nil

}
