package server

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"math/rand"
	"rpcdescriptors/gen/test"
	"strings"
	"time"
)

var LocalSignalMapper LocalSignals

func init() {
	LocalSignalMapper = make(LocalSignals)
}

type LocalSignals map[string][]LocalSignal
type LocalSignal struct {
	ID                  int32  `json:"id"`
	Name                string `json:"name"`
	Value               any    `json:"value"`
	NewMessageFieldName string `json:"newMessageFieldName"`
}

func (l LocalSignals) AddByKey(key string, signal LocalSignal) {
	if _, ok := l[key]; !ok {
		l[key] = make([]LocalSignal, 0)
	}
	l[key] = append(l[key], signal)
}

func ReadReceivedPackage(receivedPackSignals []test.Signals, receivedPackData []float32) ([]*dynamicpb.Message, error) {
	reg, err := NewLocalRegistry()
	if err != nil {
		return nil, err
	}

	for i, signal := range receivedPackSignals {
		signalDescriptor := getEnumDescriptorProto(signal)
		opts := signalDescriptor.Options()
		if proto.HasExtension(opts, test.E_ParentMessage) {
			parentExtension := proto.GetExtension(opts, test.E_ParentMessage)
			parent, ok := parentExtension.(string)
			if !ok || !strings.Contains(parent, ".") {
				return nil, fmt.Errorf("invalid parent_message format")
			}

			messageAnnotations := strings.Split(parentExtension.(string), ".")
			parentName := messageAnnotations[0]
			newMessageField := messageAnnotations[1]

			LocalSignalMapper.AddByKey(parentName, LocalSignal{
				ID:                  int32(signal.Number()),
				Name:                string(signalDescriptor.Name()),
				Value:               receivedPackData[i],
				NewMessageFieldName: newMessageField,
			})
		}
	}

	messages := make([]*dynamicpb.Message, 0)
	for k, signals := range LocalSignalMapper {
		fmt.Printf("LocalSignalMapperParentName: %v\n", k)
		// Создаем дескриптор сообщения
		desc := reg.CreateDescriptorForSignalPack(k, signals)
		// регистрируем дескриптор в localRegistry
		messageDescriptor, err := reg.RegisterDescriptor(desc)
		if err != nil {
			return nil, err
		}
		// создаём dynamic сообщение
		msg := dynamicpb.NewMessage(messageDescriptor)
		// заполняем значениями новое сообщение
		for _, signal := range signals {
			fd := messageDescriptor.Fields().ByName(protoreflect.Name(signal.NewMessageFieldName))
			if fd == nil {
				continue
			}
			msg.Set(fd, detectValueType(signal.Value))
		}
		fmt.Printf("Final Message %s: %v\n", msg.Descriptor().Name(), msg)
		messages = append(messages, msg)
	}

	return messages, nil
}

func PrepareForTransfer(msg *dynamicpb.Message) ([]byte, error) {
	msgData, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	msgDesc := msg.Descriptor()
	protoDesc := protodesc.ToDescriptorProto(msgDesc)

	bundle := &test.PackedBundle{
		Value:       msgData,
		Descriptor_: protoDesc,
	}

	return proto.Marshal(bundle)
}

func registerDescriptor(desc *descriptorpb.DescriptorProto) (protoreflect.MessageDescriptor, error) {
	fmt.Printf("RegisterDescriptor: %v\n", desc.GetName())

	fileDesc := &descriptorpb.FileDescriptorProto{
		Name:   proto.String(fmt.Sprintf("dynamic_%s.proto", desc.GetName())),
		Syntax: proto.String("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			desc, // Ваша функция создания дескриптора
		},
	}

	// 2. Конвертируем в FileDescriptor
	fd, err := protodesc.NewFile(fileDesc, protoregistry.GlobalFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create file descriptor: %w", err)
	}

	// 3. Получаем дескриптор сообщения
	msgDesc := fd.Messages().ByName(protoreflect.Name(desc.GetName()))
	if msgDesc == nil {
		return nil, fmt.Errorf("message descriptor not found")
	}

	return msgDesc, nil
}

func detectFieldType(v any) *descriptorpb.FieldDescriptorProto_Type {
	switch v.(type) {
	case float32:
		return descriptorpb.FieldDescriptorProto_TYPE_FLOAT.Enum()
	default:
		return descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum()
	}
}
func detectValueType(v any) protoreflect.Value {
	switch v.(type) {
	case float32:
		return protoreflect.ValueOfFloat32(v.(float32))
	default:
		return protoreflect.ValueOf(v)
	}
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

func getEnumDescriptorProto(enum protoreflect.Enum) protoreflect.EnumValueDescriptor {
	return enum.Descriptor().Values().ByNumber(enum.Number())
}

func createDescriptor() {

}

func randomFloat() float32 {
	return rand.Float32()*180 - 90 // Примерный диапазон для координат
}

func randomTime() uint64 {
	minV := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	maxV := time.Now().Unix()
	sec := rand.Int63n(maxV-minV) + minV
	return uint64(time.Unix(sec, 0).Unix())
}
