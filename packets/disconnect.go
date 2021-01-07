package packets

import "io"

type Disconnect struct {
	Header Header
}

func (p *Disconnect) Size() int {
	p.Header.RemainingLength = 0
	return p.Header.Size()
}

func (p *Disconnect) Marshal() ([]byte, error) {
	data := make([]byte, p.Size())
	_ = p.MarshalTo(data)
	return data, nil
}

func (p *Disconnect) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	p.Header.MarshalTo(buffer)
	return nil
}

func (p *Disconnect) Unmarshal(_ []byte) error {
	return nil
}
