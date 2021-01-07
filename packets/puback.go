package packets

import (
	"encoding/binary"
	"io"
)

type PubAck struct {
	Header    Header
	MessageID uint16
}

func (p *PubAck) Size() int {
	p.Header.RemainingLength = 2
	return p.Header.Size() + 2
}

func (p *PubAck) Marshal() ([]byte, error) {
	data := make([]byte, p.Size())
	_ = p.MarshalTo(data)
	return data, nil
}

func (p *PubAck) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	p.Header.MarshalTo(buffer)
	buffer = buffer[p.Header.Size():]
	binary.BigEndian.PutUint16(buffer, p.MessageID)
	return nil
}

func (p *PubAck) Unmarshal(bytes []byte) error {
	p.MessageID = binary.BigEndian.Uint16(bytes)
	return nil
}
