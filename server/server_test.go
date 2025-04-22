package server

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"rpcdescriptors/client"
	"rpcdescriptors/gen/test"
	"testing"
)

func Test_PackGeo(t *testing.T) {
	receivedPackSignals := []test.Signals{ // имитируем доставку сигналов
		test.Signals_SignalLat,
		test.Signals_SignalLong,
		test.Signals_SignalDiff,
	}
	receivedPackData := []float32{ // как будто у нас есть значения для сигналов
		randomFloat(),
		randomFloat(),
		randomFloat(),
	}

	dynamicMessages, err := ReadReceivedPackage(receivedPackSignals, receivedPackData)
	if err != nil {
		t.Fatal(err)
	}

	for _, dynamicMessage := range dynamicMessages {
		bs, err := PrepareForTransfer(dynamicMessage)
		if err != nil {
			t.Fatal(err)
		}

		receivedMessage, err := client.ReceiveAndParse(bs)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Received message: %v", receivedMessage)
		t.Logf("Received message name: %v", receivedMessage.Descriptor().FullName())

		receivedMessage.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			fmt.Printf("%s: %v\n", fd.Name(), v)
			return true
		})
	}
}
