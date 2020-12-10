package protodet

import (
	"github.com/golang/protobuf/proto"
)

// Protobuf with deterministic marshalling

// Single shot encoding
func Marshal(msg proto.Message) ([]byte, error) {
	buf := proto.NewBuffer(nil)
	buf.SetDeterministic(true)
	err := buf.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Single shot decoding
func Unmarshal(bs []byte, msg proto.Message) error {
	buf := proto.NewBuffer(bs)
	buf.SetDeterministic(true)
	return buf.Unmarshal(msg)
}
