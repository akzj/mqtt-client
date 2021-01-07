package packets

import "io"

type PingResp struct {
	Header Header
}

func (p *PingResp) Size() int {
	p.Header.RemainingLength = 0
	return p.Header.Size()
}

func (p *PingResp) Marshal() ([]byte, error) {
	data := make([]byte, p.Size())
	_ = p.MarshalTo(data)
	return data, nil
}

func (p *PingResp) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	p.Header.MarshalTo(buffer)
	return nil
}

func (p *PingResp) Unmarshal(_ []byte) error {
	return nil
}
