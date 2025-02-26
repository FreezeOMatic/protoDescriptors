package client

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
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

	//println(*envelope.Descriptor_.Name)
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
