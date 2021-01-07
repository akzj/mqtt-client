package packets

import (
	"encoding/binary"
	"io"
	"unsafe"
)

type Subscribe struct {
	Header    Header
	MessageID uint16
	Topics    []string
	Qoss      []byte
}

func (s *Subscribe) Size() int {
	var size int
	size += 2 // message ID
	for _, topic := range s.Topics {
		size += 2 //topic length
		size += len(topic)
		size += 1 // qos
	}
	s.Header.RemainingLength = size
	return s.Header.Size() + size
}

func (s *Subscribe) Marshal() ([]byte, error) {
	buffer := make([]byte, s.Size())
	_ = s.MarshalTo(buffer)
	return buffer, nil
}

func (s *Subscribe) MarshalTo(buffer []byte) error {
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

		buffer[0] = s.Qoss[i]
		buffer = buffer[1:]
	}
	return nil
}

func (s *Subscribe) Unmarshal(data []byte) error {
	s.MessageID = binary.BigEndian.Uint16(data)
	data = data[2:]
	for len(data) > 0 {
		slen := binary.BigEndian.Uint16(data)
		data = data[2:]
		sbuff := data[:slen]
		s.Topics = append(s.Topics, *(*string)(unsafe.Pointer(&sbuff)))
		s.Qoss = append(s.Qoss, data[0])
		data = data[1:]
	}
	return nil
}
