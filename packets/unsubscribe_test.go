package packets

import (
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"testing"
)

func TestUnSubscribe(t *testing.T) {
	unsub := UnSubscribe{
		Header: Header{
			MessageType:     UnsubscribeType,
			Dup:             false,
			Qos:             1,
			Retain:          false,
			RemainingLength: 8,
		},
		MessageID: 1024,               //2
		Topics:    []string{"a", "b"}, //2+4
	}

	pUnsub := &paho.UnsubscribePacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     UnsubscribeType,
			Dup:             false,
			Qos:             1,
			Retain:          false,
			RemainingLength: 8,
		},
		MessageID: 1024,               //2
		Topics:    []string{"a", "b"}, //2+4
	}

	data, _ := unsub.Marshal()
	if reflect.DeepEqual(data, marshalPacket(pUnsub)) == false {
		t.Fatal("marshal failed")
	}

	var unsub2 UnSubscribe

	unsub2.Header = unsub.Header

	unsub2.Unmarshal(data[unsub.Header.Size():])

	if reflect.DeepEqual(unsub, unsub2) == false {
		t.Fatal("Unmarshal failed")
	}

}
