package client

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"log"
	"rpcdescriptors/gen/test"
)

func ServeMessage(messageBytes []byte) {
	var envelope test.Envelope
	if err := proto.Unmarshal(messageBytes, &envelope); err != nil {
		log.Fatalf("Failed to unmarshal Envelope: %v", err)
	}

	//fmt.Printf("%v\n", envelope.Descriptor_.Field)

	file := &descriptorpb.FileDescriptorProto{
		Name:        proto.String("received.proto"),
		Syntax:      proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{envelope.GetDescriptor_()},
	}

	files, err := protodesc.NewFiles(
		&descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{file},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create file registry: %v", err)
	}

	msgDescriptor, err := files.FindDescriptorByName(protoreflect.FullName(envelope.GetDescriptor_().GetName()))
	if err != nil {
		log.Fatalf("Failed to find message descriptor: %v", err)
	}

	msg := dynamicpb.NewMessage(msgDescriptor.(protoreflect.MessageDescriptor))

	if err = proto.Unmarshal(envelope.GetData(), msg); err != nil {
		log.Fatalf("Failed to unmarshal into dynamic message: %v", err)
	}

	fmt.Printf("Message type: %s\n", envelope.GetDescriptor_().GetName())
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fmt.Printf("  Field %s: %v\n", fd.Name(), v.Interface())
		return true
	})
}

func ServeEnumMessage(messageBytes []byte) {
	var m test.MessageWithEnum
	if err := proto.Unmarshal(messageBytes, &m); err != nil {
		log.Fatalf("Failed to unmarshal Envelope: %v", err)
	}

	enumDescValues := m.GetEnumSignal().Descriptor().Values()
	for i := 0; i < enumDescValues.Len(); i++ {
		descriptor := GetReliedDescriptor(enumDescValues.Get(i))

		message, err := GetConcreteMessageFromDescriptor(descriptor)
		if err != nil {
			log.Fatalf("Failed to get concrete message: %v", err)
		}

		switch message.(type) {
		case *test.Undefined:
			fmt.Printf("message: %v\n", message)
			message.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
				fmt.Printf("  Field %s: %v\n", fd.Name(), v.Interface())
				return true
			})
		}
	}
}

func GetReliedDescriptor(enumValue protoreflect.EnumValueDescriptor) *descriptorpb.DescriptorProto {
	options := enumValue.Options()

	messageDescriptor := proto.GetExtension(options, test.E_Descriptor).(*descriptorpb.DescriptorProto)
	if messageDescriptor == nil {
		fmt.Println("Message descriptor not found for enum value:", enumValue.Name())
		return nil
	}

	fmt.Printf("Enum value: %s\n", enumValue.Name())
	fmt.Printf("Message name: %s\n", messageDescriptor.GetName())
	fmt.Println("Fields:")
	for _, field := range messageDescriptor.GetField() {
		fmt.Printf("  - %s (type: %v, number: %d)\n", field.GetName(), field.GetType(), field.GetNumber())
	}
	fmt.Println("---------------------------")
	messageAbbreviation := proto.GetExtension(options, test.E_Abbr).(string)
	fmt.Printf("Additional info: %s\n", messageAbbreviation)

	return messageDescriptor
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

//file := &descriptorpb.FileDescriptorProto{
//		Name:     proto.String("enums.proto"),
//		Syntax:   proto.String("proto3"),
//		EnumType: []*descriptorpb.EnumDescriptorProto{envelope.GetDescriptor_()},
//	}
//
//	files, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{file}})
//	if err != nil {
//		log.Fatalf("Failed to create file registry: %v", err)
//	}
//
//	enumDescriptor, err := files.FindDescriptorByName(protoreflect.FullName(envelope.GetDescriptor_().GetName()))
//	if err != nil {
//		log.Fatalf("Failed to find enum descriptor: %v", err)
//	}
//
//	msg := dynamicpb.NewEnumType(enumDescriptor.(protoreflect.EnumDescriptor))
//
//	if err := proto.Unmarshal(envelope.EnumData, msg); err != nil {
//		log.Fatalf("Failed to unmarshal into dynamic message: %v", err)
//	}
//
//	fmt.Printf("Message type: %s\n", envelope.GetDescriptor_().GetName())
//	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
//		fmt.Printf("  Field %s: %v\n", fd.Name(), v.Interface())
//		return true
//	})
