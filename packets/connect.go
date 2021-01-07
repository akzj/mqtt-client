package packets

import (
	"encoding/binary"
	"io"
	"unsafe"
)

type Connect struct {
	Header          Header
	ProtocolName    string
	ProtocolVersion byte
	CleanSession    bool
	WillFlag        bool
	WillQos         byte
	WillRetain      bool
	UsernameFlag    bool
	PasswordFlag    bool
	ReservedBit     byte
	Keepalive       uint16

	ClientIdentifier string
	WillTopic        string
	WillMessage      []byte
	Username         string
	Password         []byte
}

func (c *Connect) Size() int {
	var size int

	size += 2 + len(c.ProtocolName)     //ProtocolName
	size += 1                           //ProtocolVersion
	size += 1                           //flags
	size += 2                           //Keepalive
	size += 2 + len(c.ClientIdentifier) // ClientIdentifier
	if c.WillFlag {
		size += 2 + len(c.WillTopic)   //WillTopic
		size += 2 + len(c.WillMessage) //WillMessage
	}
	if c.UsernameFlag {
		size += 2 + len(c.Username) //Username
	}
	if c.PasswordFlag {
		size += 2 + len(c.Password) //Password
	}

	c.Header.RemainingLength = size
	return size + c.Header.Size()
}

func (c *Connect) Marshal() ([]byte, error) {
	data := make([]byte, c.Size())
	_ = c.MarshalTo(data)
	return data, nil
}

func (c *Connect) MarshalTo(buffer []byte) error {
	if len(buffer) < c.Size() {
		return io.ErrShortBuffer
	}

	c.Header.MarshalTo(buffer)
	buffer = buffer[c.Header.Size():]

	binary.BigEndian.PutUint16(buffer, uint16(len(c.ProtocolName))) //c.ProtocolName length
	buffer = buffer[2:]
	copy(buffer, c.ProtocolName) //ProtocolName
	buffer = buffer[len(c.ProtocolName):]

	buffer[0] = c.ProtocolVersion //ProtocolVersion
	buffer = buffer[1:]

	var flag byte

	if c.CleanSession {
		flag |= 1 << 1
	}
	if c.WillFlag {
		flag |= 1 << 2
	}
	flag |= c.WillQos << 3
	if c.WillRetain {
		flag |= 1 << 5
	}
	if c.PasswordFlag {
		flag |= 1 << 6
	}
	if c.UsernameFlag {
		flag |= 1 << 7
	}

	buffer[0] = flag
	buffer = buffer[1:] //flags

	binary.BigEndian.PutUint16(buffer, c.Keepalive) //Keepalive
	buffer = buffer[2:]

	binary.BigEndian.PutUint16(buffer, uint16(len(c.ClientIdentifier))) //c.ClientIdentifier length
	buffer = buffer[2:]
	copy(buffer, c.ClientIdentifier) //ClientIdentifier
	buffer = buffer[len(c.ClientIdentifier):]

	if c.WillFlag {
		binary.BigEndian.PutUint16(buffer, uint16(len(c.WillTopic))) //c.WillTopic length
		buffer = buffer[2:]
		copy(buffer, c.WillTopic) //WillTopic
		buffer = buffer[len(c.WillTopic):]

		binary.BigEndian.PutUint16(buffer, uint16(len(c.WillMessage))) //c.WillMessage length
		buffer = buffer[2:]
		copy(buffer, c.WillMessage) //WillMessage
		buffer = buffer[len(c.WillMessage):]
	}
	if c.UsernameFlag {
		binary.BigEndian.PutUint16(buffer, uint16(len(c.Username))) //c.Username length
		buffer = buffer[2:]
		copy(buffer, c.Username) //Username
		buffer = buffer[len(c.Username):]
	}
	if c.PasswordFlag {
		binary.BigEndian.PutUint16(buffer, uint16(len(c.Password))) //c.Password length
		buffer = buffer[2:]
		copy(buffer, c.Password) //Password
		buffer = buffer[len(c.Password):]
	}
	return nil
}

func (c *Connect) Unmarshal(data []byte) error {

	//ProtocolName
	{
		dataLength := binary.BigEndian.Uint16(data)
		data = data[2:]
		dataBuffer := data[:dataLength]
		data = data[dataLength:]
		c.ProtocolName = *(*string)(unsafe.Pointer(&dataBuffer))
	}

	//ProtocolVersion
	c.ProtocolVersion = data[0]
	data = data[1:]

	options := data[0]
	data = data[1:]

	c.ReservedBit = 1 & options
	c.CleanSession = 1&(options>>1) > 0
	c.WillFlag = 1&(options>>2) > 0
	c.WillQos = 3 & (options >> 3)
	c.WillRetain = 1&(options>>5) > 0
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0

	//Keepalive
	c.Keepalive = binary.BigEndian.Uint16(data)
	data = data[2:]

	//clientIdentifier
	{
		dataLength := binary.BigEndian.Uint16(data)
		data = data[2:]
		dataBuffer := data[:dataLength]
		data = data[dataLength:]
		c.ClientIdentifier = *(*string)(unsafe.Pointer(&dataBuffer))
	}

	if c.WillFlag {
		//WillTopic
		{
			dataLength := binary.BigEndian.Uint16(data)
			data = data[2:]
			dataBuffer := data[:dataLength]
			data = data[dataLength:]
			c.WillTopic = *(*string)(unsafe.Pointer(&dataBuffer))
		}

		//WillMessage
		{
			dataLength := binary.BigEndian.Uint16(data)
			data = data[2:]
			dataBuffer := data[:dataLength]
			data = data[dataLength:]
			c.WillMessage = dataBuffer
		}
	}

	if c.UsernameFlag {
		//Username
		dataLength := binary.BigEndian.Uint16(data)
		data = data[2:]
		dataBuffer := data[:dataLength]
		data = data[dataLength:]
		c.Username = *(*string)(unsafe.Pointer(&dataBuffer))
	}

	if c.PasswordFlag {
		//Password
		dataLength := binary.BigEndian.Uint16(data)
		data = data[2:]
		dataBuffer := data[:dataLength]
		data = data[dataLength:]
		c.Password = dataBuffer
	}
	return nil
}

//Validate performs validation of the fields of a Connect packet
func (c *Connect) Validate() byte {
	if c.PasswordFlag && !c.UsernameFlag {
		return ErrRefusedBadUsernameOrPassword
	}
	if c.ReservedBit != 0 {
		//Bad reserved bit
		return ErrProtocolViolation
	}
	if (c.ProtocolName == "MQIsdp" && c.ProtocolVersion != 3) ||
		(c.ProtocolName == "MQTT" && c.ProtocolVersion != 4) {
		//Mismatched or unsupported protocol version
		return ErrRefusedBadProtocolVersion
	}
	if c.ProtocolName != "MQIsdp" && c.ProtocolName != "MQTT" {
		//Bad protocol name
		return ErrProtocolViolation
	}
	if len(c.ClientIdentifier) > 65535 ||
		len(c.Username) > 65535 ||
		len(c.Password) > 65535 {
		//Bad size field
		return ErrProtocolViolation
	}
	if len(c.ClientIdentifier) == 0 && !c.CleanSession {
		//Bad client identifier
		return ErrRefusedIDRejected
	}
	return Accepted
}
