package packets

import (
	"bytes"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	conn := Connect{
		Header: Header{
			MessageType:     ConnectType,
			Dup:             false,
			Qos:             1,
			Retain:          false,
			RemainingLength: 0,
		},
		ProtocolName:     "MQTT",
		ProtocolVersion:  4,
		CleanSession:     true,
		WillFlag:         true,
		WillQos:          1,
		WillRetain:       true,
		UsernameFlag:     true,
		PasswordFlag:     true,
		ReservedBit:      0,
		Keepalive:        32,
		ClientIdentifier: "ClientIdentifier",
		WillTopic:        "WillTopic",
		WillMessage:      []byte("WillMessage"),
		Username:         "Username",
		Password:         []byte("Password"),
	}

	pconn := &paho.ConnectPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     ConnectType,
			Dup:             false,
			Qos:             1,
			Retain:          false,
			RemainingLength: 0,
		},
		ProtocolName:     "MQTT",
		ProtocolVersion:  4,
		CleanSession:     true,
		WillFlag:         true,
		WillQos:          1,
		WillRetain:       true,
		UsernameFlag:     true,
		PasswordFlag:     true,
		ReservedBit:      0,
		Keepalive:        32,
		ClientIdentifier: "ClientIdentifier",
		WillTopic:        "WillTopic",
		WillMessage:      []byte("WillMessage"),
		Username:         "Username",
		Password:         []byte("Password"),
	}

	if conn.Validate() != Accepted {
		t.Fatal("Validate failed")
	}
	if pconn.Validate() != Accepted {
		t.Fatal("Validate failed")
	}

	data, _ := conn.Marshal()
	if reflect.DeepEqual(data, marshalPacket(pconn)) == false {
		t.Fatal("marshal failed")
	}

	var conn2 = Connect{
		Header: conn.Header,
	}

	conn2.Unmarshal(data[conn.Header.Size():])

	if reflect.DeepEqual(conn2, conn) == false {
		t.Fatal("Unmarshal failed")
	}

	p, d, err := UnmarshalPacket(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(d, data) == false {
		fmt.Println(data)
		fmt.Println(d)
		t.Fatal("UnmarshalPacket failed",len(data),len(data))
	}
	if reflect.DeepEqual(p, &conn) == false {
		fmt.Println(p)
		fmt.Println(conn)
		t.Fatal("UnmarshalPacket failed")
	}
}
