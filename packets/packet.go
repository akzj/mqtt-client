package packets

import "errors"

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