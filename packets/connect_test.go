package packets

import "testing"

func TestConnect(t *testing.T) {
	conn := Connect{
		Header: Header{
			MessageType:     ConnectType,
			Dup:             false,
			Qos:             1,
			Retain:          false,
			RemainingLength: 0,
		},
		ProtocolName:     "",
		ProtocolVersion:  0,
		CleanSession:     false,
		WillFlag:         false,
		WillQos:          0,
		WillRetain:       false,
		UsernameFlag:     false,
		PasswordFlag:     false,
		ReservedBit:      0,
		Keepalive:        0,
		ClientIdentifier: "",
		WillTopic:        "",
		WillMessage:      nil,
		Username:         "",
		Password:         nil,
	}
}
