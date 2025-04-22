package server

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type LocalRegistry struct {
	files *protoregistry.Files
}

func NewLocalRegistry() (*LocalRegistry, error) {
	reg, err := createLocalRegistry()
	if err != nil {
		return nil, fmt.Errorf("error creating local registry: %v", err)
	}
	return &LocalRegistry{
		files: reg,
	}, nil
}
func createLocalRegistry() (*protoregistry.Files, error) {
	fileDesc := &descriptorpb.FileDescriptorProto{
		Name:        proto.String("dynamic_cache.proto"),
		Syntax:      proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{},
	}
	fds := descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{fileDesc}}

	return protodesc.NewFiles(&fds)
}

func (reg *LocalRegistry) CreateDescriptorForSignalPack(name string, signals []LocalSignal) *descriptorpb.DescriptorProto {
	desc := &descriptorpb.DescriptorProto{
		Name: proto.String(fmt.Sprintf("%sPacked", name)),
	}
	// Добавляем поля в дескриптор
	for i, signal := range signals {
		desc.Field = append(desc.Field, &descriptorpb.FieldDescriptorProto{
			Name:   proto.String(signal.NewMessageFieldName),
			Number: proto.Int32(int32(i + 1)), // Нумерация полей начинается с 1
			Type:   detectFieldType(signal.Value),
			Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(), // сделаем его optional
		})
	}

	return desc
}

func (reg *LocalRegistry) RegisterDescriptor(desc *descriptorpb.DescriptorProto) (protoreflect.MessageDescriptor, error) {
	fmt.Printf("RegisterDescriptor: %v\n", desc.GetName())
	fileDesc := &descriptorpb.FileDescriptorProto{
		Name:   proto.String(fmt.Sprintf("dynamic_%s.proto", desc.GetName())),
		Syntax: proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			desc,
		},
	}
	// Конвертируем в FileDescriptor
	fd, err := protodesc.NewFile(fileDesc, reg.files)
	if err != nil {
		return nil, fmt.Errorf("failed to create file descriptor: %w", err)
	}
	// Получаем дескриптор сообщения
	msgDesc := fd.Messages().ByName(protoreflect.Name(desc.GetName()))
	if msgDesc == nil {
		return nil, fmt.Errorf("message descriptor not found")
	}
	return msgDesc, nil
}
