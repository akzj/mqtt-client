package packets

import "io"

type PingReq struct {
	Header Header
}

func (p *PingReq) Size() int {
	p.Header.RemainingLength = 0
	return p.Header.Size()
}

func (p *PingReq) Marshal() ([]byte, error) {
	data := make([]byte, p.Size())
	_ = p.MarshalTo(data)
	return data, nil
}

func (p *PingReq) MarshalTo(buffer []byte) error {
	if len(buffer) < p.Size() {
		return io.ErrShortBuffer
	}
	p.Header.MarshalTo(buffer)
	return nil
}

func (p *PingReq) Unmarshal(_ []byte) error {
	return nil
}
