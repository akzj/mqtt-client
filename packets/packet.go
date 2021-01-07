package packets

import (
	"errors"
	"fmt"
	"io"
)

//Below are the constants assigned to each of the MQTT packet types
const (
	ConnectType     = 1
	ConnackType     = 2
	PublishType     = 3
	PubackType      = 4
	PubrecType      = 5
	PubrelType      = 6
	PubcompType     = 7
	SubscribeType   = 8
	SubackType      = 9
	UnsubscribeType = 10
	UnsubackType    = 11
	PingreqType     = 12
	PingrespType    = 13
	DisconnectType  = 14
)

type Packet interface {
	Size() int
	Marshal() ([]byte, error)
	MarshalTo(buffer []byte) error
	Unmarshal(data []byte) error
}

func UnmarshalPacket(reader io.Reader) (Packet, []byte, error) {
	var header Header
	var headerBuffer [5]byte
	var pos = 0
	_, err := io.ReadFull(reader, headerBuffer[pos:pos+1])
	if err != nil {
		return nil, nil, err
	}
	pos++
	var typeAndFlags = headerBuffer[0]
	header.MessageType = typeAndFlags >> 4
	header.Dup = (typeAndFlags>>3)&0x01 > 0
	header.Qos = (typeAndFlags >> 1) & 0x03
	header.Retain = typeAndFlags&0x01 > 0

	var rLength uint32
	var multiplier uint32
	for {
		if _, err := io.ReadFull(reader, headerBuffer[pos:pos+1]); err != nil {
			return nil, nil, err
		}
		digit := headerBuffer[pos]
		pos++
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
		if multiplier > 28 {
			return nil, nil, fmt.Errorf("malformed remaining length")
		}
	}
	header.RemainingLength = int(rLength)

	data := make([]byte, int(rLength)+pos)
	copy(data, headerBuffer[:pos])
	if _, err := io.ReadFull(reader, data[pos:]); err != nil {
		return nil, nil, err
	}

	switch header.MessageType {
	case ConnackType:
		var packet ConnAck
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case ConnectType:
		var packet Connect
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case DisconnectType:
		var packet Disconnect
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PingreqType:
		var packet PingReq
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PingrespType:
		var packet PingResp
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PubackType:
		var packet PubAck
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PubcompType:
		var packet PubComp
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PublishType:
		var packet Publish
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PubrecType:
		var packet PubRec
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case PubrelType:
		var packet PubRel
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case SubackType:
		var packet SubAck
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case SubscribeType:
		var packet Subscribe
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case UnsubackType:
		var packet UnSubAck
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	case UnsubscribeType:
		var packet UnSubscribe
		packet.Header = header
		_ = packet.Unmarshal(data[pos:])
		return &packet, data, nil
	default:
		return nil, data, fmt.Errorf("unknowm packet type %d", header.MessageType)
	}
}

//Below are the const definitions for error codes returned by
//Connect()
const (
	Accepted                        = 0x00
	ErrRefusedBadProtocolVersion    = 0x01
	ErrRefusedIDRejected            = 0x02
	ErrRefusedServerUnavailable     = 0x03
	ErrRefusedBadUsernameOrPassword = 0x04
	ErrRefusedNotAuthorised         = 0x05
	ErrNetworkError                 = 0xFE
	ErrProtocolViolation            = 0xFF
)

//ConnackReturnCodes is a map of the error codes constants for Connect()
//to a string representation of the error
var ConnackReturnCodes = map[uint8]string{
	0:   "Connection Accepted",
	1:   "Connection Refused: Bad Protocol Version",
	2:   "Connection Refused: Client Identifier Rejected",
	3:   "Connection Refused: Server Unavailable",
	4:   "Connection Refused: Username or Password in unknown format",
	5:   "Connection Refused: Not Authorised",
	254: "Connection Error",
	255: "Connection Refused: Protocol Violation",
}

//ConnErrors is a map of the errors codes constants for Connect()
//to a Go error
var ConnErrors = map[byte]error{
	Accepted:                        nil,
	ErrRefusedBadProtocolVersion:    errors.New("Unnacceptable protocol version"),
	ErrRefusedIDRejected:            errors.New("Identifier rejected"),
	ErrRefusedServerUnavailable:     errors.New("Server Unavailable"),
	ErrRefusedBadUsernameOrPassword: errors.New("Bad user name or password"),
	ErrRefusedNotAuthorised:         errors.New("Not Authorized"),
	ErrNetworkError:                 errors.New("Network Error"),
	ErrProtocolViolation:            errors.New("Protocol Violation"),
}
