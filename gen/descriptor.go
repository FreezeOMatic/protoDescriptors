package gen

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
)

var FDS descriptorpb.FileDescriptorSet

func init() {
	fdsBytes, err := os.ReadFile("gen/descriptor.pb")
	if err != nil {
		panic(err)
	}

	if err = proto.Unmarshal(fdsBytes, &FDS); err != nil {
		panic(err)
	}

}
