package packets

import (
	"bytes"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"testing"
)

func TestPubRec(t *testing.T) {

	pr := PubRec{
		Header: Header{
			MessageType:     PubrecType,
			Dup:             false,
			Qos:             2,
			Retain:          false,
			RemainingLength: 2,
		},
		MessageID: 1024,
	}

	ppr := paho.PubrecPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     PubrecType,
			Dup:             false,
			Qos:             2,
			Retain:          false,
			RemainingLength: 2,
		},
		MessageID: 1024,
	}

	data, _ := pr.Marshal()

	var buffer = new(bytes.Buffer)
	ppr.Write(buffer)

	if bytes.Compare(data, buffer.Bytes()) != 0 {
		t.Fatal("pubRec marshal failed")
	}

	var pr2 = PubRec{Header: pr.Header}

	pr2.Unmarshal(data[pr.Header.Size():])

	if reflect.DeepEqual(pr, pr2) != true {
		t.Fatal("Unmarshal failed")
	}
}
