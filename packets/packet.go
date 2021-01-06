package packets

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
	Unmarshal([]byte) error
}
