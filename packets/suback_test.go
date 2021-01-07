package packets

import (
	"bytes"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"testing"
)

func marshalPacket(packet paho.ControlPacket) []byte {
	var buffer = new(bytes.Buffer)
	packet.Write(buffer)
	return buffer.Bytes()
}

func TestSubAck(t *testing.T) {
	sa := SubAck{
		Header: Header{
			MessageType:     SubackType,
			Dup:             false,
			Qos:             0,
			Retain:          false,
			RemainingLength: 5 + 2,
		},
		MessageID:   1024,
		ReturnCodes: []byte{1, 2, 2, 2, 2},
	}

	psp := &paho.SubackPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     SubackType,
			Dup:             false,
			Qos:             0,
			Retain:          false,
			RemainingLength: 5 + 2,
		},
		MessageID:   1024,
		ReturnCodes: []byte{1, 2, 2, 2, 2},
	}

	data, _ := sa.Marshal()
	if bytes.Compare(data, marshalPacket(psp)) != 0 {
		t.Fatal("Marshal failed")
	}

	if _, err := paho.ReadPacket(bytes.NewReader(data)); err != nil {
		t.Fatal(err)
	}
}
