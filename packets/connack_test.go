package packets

import (
	"bytes"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"testing"
)

func TestConnack(t *testing.T) {
	ca := ConnAck{
		Header: Header{
			MessageType:     ConnackType,
			Dup:             true,
			Qos:             0,
			Retain:          true,
			RemainingLength: 0,
		},
		SessionPresent: true,
		ReturnCode:     Accepted,
	}

	pca := paho.ConnackPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     ConnackType,
			Dup:             true,
			Qos:             0,
			Retain:          true,
			RemainingLength: 0,
		},
		SessionPresent: true,
		ReturnCode:     Accepted,
	}

	data, _ := ca.Marshal()
	if reflect.DeepEqual(data, marshalPacket(&pca)) == false {
		t.Fatal("marshal failed")
	}

	_, _, err := UnmarshalPacket(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
}
