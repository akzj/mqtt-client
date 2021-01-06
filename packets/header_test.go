package packets

import (
	"bytes"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"testing"
)

func TestHeader(t *testing.T) {

	EPHeader := paho.FixedHeader{
		MessageType:     PublishType,
		Dup:             true,
		Qos:             2,
		Retain:          true,
		RemainingLength: 10240,
	}

	var buffer bytes.Buffer
	packet := paho.DisconnectPacket{FixedHeader: EPHeader}
	_ = packet.Write(&buffer)

	var header = Header{
		MessageType:     PublishType,
		Dup:             true,
		Qos:             2,
		Retain:          true,
		RemainingLength: 10240,
	}

	if buffer.Len() != header.Size() {
		t.Fatalf("header size error")
	}
	data, _ := header.Marshal()
	if bytes.Compare(buffer.Bytes(), data) != 0 {
		t.Fatalf("marshal header failed")
	}

}
