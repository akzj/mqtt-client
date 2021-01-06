package packets

type Header struct {
	MessageType     byte
	Dup             bool
	Qos             byte
	Retain          bool
	RemainingLength int
}

func (h Header) Size() int {
	var size = 1
	length := h.RemainingLength
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		size++
		if length == 0 {
			break
		}
	}
	return size
}

func (h Header) Marshal() ([]byte, error) {
	buffer := make([]byte, h.Size())
	_ = h.MarshalTo(buffer)
	return buffer, nil
}

func (h Header) MarshalTo(buffer []byte) error {
	var pos = 1
	var flag byte
	flag = h.MessageType << 4
	if h.Dup {
		flag |= 1 << 3
	}
	flag |= h.Qos << 1
	if h.Retain {
		flag |= 1
	}
	buffer[0] = flag
	length := h.RemainingLength
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		buffer[pos] = digit
		pos++
		if length == 0 {
			break
		}
	}
	return nil
}
