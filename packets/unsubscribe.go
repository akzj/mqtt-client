package packets

import (
	"encoding/binary"
	"io"
	"unsafe"
)

type UnSubscribe struct {
	Header    Header
	MessageID uint16
	Topics    []string
}

func (s *UnSubscribe) Size() int {
	var size int
	size += 2 // message ID
	for _, topic := range s.Topics {
		size += 2 //topic length
		size += len(topic)
	}
	s.Header.RemainingLength = size
	return s.Header.Size() + size
}

func (s *UnSubscribe) Marshal() ([]byte, error) {
	buffer := make([]byte, s.Size())
	_ = s.MarshalTo(buffer)
	return buffer, nil
}

func (s *UnSubscribe) MarshalTo(buffer []byte) error {
	if len(buffer) < s.Size() {
		return io.ErrShortBuffer
	}
	_ = s.Header.MarshalTo(buffer)
	buffer = buffer[s.Header.Size():]

	binary.BigEndian.PutUint16(buffer, s.MessageID)
	buffer = buffer[2:]

	for i := range s.Topics {
		topic := &s.Topics[i]

		binary.BigEndian.PutUint16(buffer, uint16(len(*topic)))
		buffer = buffer[2:]

		copy(buffer, *topic)
		buffer = buffer[len(*topic):]
	}
	return nil
}

func (s *UnSubscribe) Unmarshal(data []byte) error {
	s.MessageID = binary.BigEndian.Uint16(data)
	data = data[2:]
	for len(data) > 0 {
		dataLength := binary.BigEndian.Uint16(data)
		data = data[2:]
		dataBuffer := data[:dataLength]
		data = data[dataLength:]
		s.Topics = append(s.Topics, *(*string)(unsafe.Pointer(&dataBuffer)))
	}
	return nil
}
