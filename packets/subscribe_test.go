package packets

import (
	"bytes"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"testing"
)

func TestSubscribe(t *testing.T) {

	sub := Subscribe{
		Header: Header{
			MessageType:     SubscribeType,
			Dup:             false,
			Qos:             2,
			Retain:          false,
			RemainingLength: 0,
		},
		MessageID: 1024,
		Topics:    []string{"a", "b"},
		Qoss:      []byte{1, 2},
	}

	data, _ := sub.Marshal()

	psub := &paho.SubscribePacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     SubscribeType,
			Dup:             false,
			Qos:             2,
			Retain:          false,
			RemainingLength: 10,
		},
		MessageID: 1024,               //2
		Topics:    []string{"a", "b"}, //2 + 4
		Qoss:      []byte{1, 2},       //2
	}
	if psub.FixedHeader.RemainingLength != sub.Header.RemainingLength || sub.Header.RemainingLength != 10 {
		t.Fatal("sub.Header.RemainingLength error", sub.Header.RemainingLength, psub.FixedHeader.RemainingLength)
	}
	if reflect.DeepEqual(data, marshalPacket(psub)) != true {
		t.Fatal("marshal failed")
	}
	_, err := paho.ReadPacket(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	var sub2 = Subscribe{
		Header: sub.Header,
	}
	sub2.Unmarshal(data[sub.Header.Size():])
	if reflect.DeepEqual(sub, sub2) == false {
		t.Fatal("Unmarshal")
	}
}
