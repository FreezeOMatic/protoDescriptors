package server

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"log"
	"rpcdescriptors/gen/test"
)

func GetMessage() []byte {
	//person := &test.Person{
	//	Name: "John Doe",
	//	Age:  30,
	//}

	person := &test.NextPerson{
		Id:    "agsglg1p2g212",
		Order: 2,
	}

	data, err := proto.Marshal(person)
	if err != nil {
		log.Fatalf("Failed to marshal Envelope: %v", err)
	}

	descriptor := getDescriptorProto(person)

	envelope := &test.Envelope{
		Data:        data,
		Descriptor_: descriptor,
	}

	envdata, err := proto.Marshal(envelope)
	if err != nil {
		log.Fatalf("Failed to marshal Envelope: %v", err)
	}

	return envdata
}

func getDescriptorProto(msg protoreflect.ProtoMessage) *descriptorpb.DescriptorProto {
	fileDesc := msg.ProtoReflect().Descriptor().ParentFile()

	// Преобразуем FileDescriptor в FileDescriptorProto
	fileDescProto := protodesc.ToFileDescriptorProto(fileDesc)

	// Ищем DescriptorProto для конкретного сообщения
	for _, msgDescProto := range fileDescProto.MessageType {
		if msgDescProto.GetName() == string(msg.ProtoReflect().Descriptor().Name()) {
			return msgDescProto
		}
	}

	return nil
}
