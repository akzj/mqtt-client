package packets

import "io"

type ConnAck struct {
	Header         Header
	SessionPresent bool
	ReturnCode     byte
}

func (c *ConnAck) Size() int {
	c.Header.RemainingLength = 2
	return c.Header.Size() + 2
}

func (c *ConnAck) Marshal() ([]byte, error) {
	data := make([]byte, c.Size())
	_ = c.MarshalTo(data)
	return data, nil
}

func (c *ConnAck) MarshalTo(buffer []byte) error {
	if len(buffer) < c.Size() {
		return io.ErrShortBuffer
	}
	if c.SessionPresent {
		buffer[0] = 1
	}
	buffer[1] = c.ReturnCode
	return nil
}

func (c *ConnAck) Unmarshal(buffer []byte) error {
	if buffer[0] > 0 {
		c.SessionPresent = true
	}
	c.ReturnCode = buffer[1]
	return nil
}
