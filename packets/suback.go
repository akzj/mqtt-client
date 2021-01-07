package packets

import (
	"encoding/binary"
	"io"
)

type SubAck struct {
	Header      Header
	MessageID   uint16
	ReturnCodes []byte
}

func (p *SubAck) Size() int {
	p.Header.RemainingLength = 2 + len(p.ReturnCodes)
	return p.Header.Size() + 2 + len(p.ReturnCodes)
}

func (p *SubAck) Marshal() ([]byte, error) {
	data := make([]byte, p.Size())
	_ = p.MarshalTo(data)
	return data, nil
}

func (p *SubAck) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	p.Header.MarshalTo(buffer)
	buffer = buffer[p.Header.Size():]
	binary.BigEndian.PutUint16(buffer, p.MessageID)
	buffer = buffer[2:]
	copy(buffer, p.ReturnCodes)
	return nil
}

func (p *SubAck) Unmarshal(data []byte) error {
	p.MessageID = binary.BigEndian.Uint16(data)
	data = data[2:]
	p.ReturnCodes = data
	return nil
}
