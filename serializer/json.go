package serializer

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON() converts protocol buffer message to JSON string
func ProtobufToJSON(message proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		UseEnumNumbers:    false,
		EmitDefaultValues: true,
		Indent:            " ",
		UseProtoNames:     false,
	}

	msgBytes, err := marshaler.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal message into json bytes: %w", err)
	}

	return msgBytes, nil
}
