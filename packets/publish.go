package packets

import (
	"encoding/binary"
	"io"
	"unsafe"
)

type Publish struct {
	Header    Header
	TopicName string
	MessageID uint16
	Payload   []byte
}

func (p *Publish) Unmarshal(buffer []byte) error {
	length := binary.BigEndian.Uint16(buffer)
	buffer = buffer[2:]
	topicName := buffer[:length]
	buffer = buffer[length:]
	p.TopicName = *(*string)(unsafe.Pointer(&topicName))

	if p.Header.Qos > 0 {
		p.MessageID = binary.BigEndian.Uint16(buffer)
		buffer = buffer[2:]
	}
	p.Payload = buffer
	return nil
}

func (p *Publish) Size() int {
	var size int

	size += 2 + len(p.TopicName)
	if p.Header.Qos > 0 {
		size += 2 // message id
	}
	size += len(p.Payload)

	p.Header.RemainingLength = size
	size += p.Header.Size()
	return size
}

func (p *Publish) Marshal() ([]byte, error) {
	buffer := make([]byte, p.Size())
	_ = p.MarshalTo(buffer)
	return buffer, nil
}

func (p *Publish) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	_ = p.Header.MarshalTo(buffer)
	buffer = buffer[p.Header.Size():]

	binary.BigEndian.PutUint16(buffer, uint16(len(p.TopicName)))
	buffer = buffer[2:]

	copy(buffer, p.TopicName)
	buffer = buffer[len(p.TopicName):]

	if p.Header.Qos > 0 {
		binary.BigEndian.PutUint16(buffer, uint16(p.MessageID))
		buffer = buffer[2:]
	}
	copy(buffer, p.Payload)
	return nil
}
