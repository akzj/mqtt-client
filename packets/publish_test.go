package packets

import (
	"bytes"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	p := Publish{
		header: Header{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}

	pp := paho.PublishPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}

	var buffer = new(bytes.Buffer)
	pp.Write(buffer)

	data, _ := p.Marshal()

	if bytes.Compare(buffer.Bytes(), data) != 0 {
		t.Error("publish marshal failed")
		fmt.Println(buffer.Bytes())
		fmt.Println(data)
	}

	var np Publish
	np.header = p.header

	np.Unmarshal(data[np.header.Size():])

	data2, _ := np.Marshal()

	if bytes.Compare(buffer.Bytes(), data2) != 0 {
		t.Error("publish marshal failed")
		fmt.Println(buffer.Bytes())
		fmt.Println(data2)
	}

}

func BenchmarkPahoPublishWrite(b *testing.B) {
	pp := paho.PublishPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}
	for i := 0; i < b.N; i++ {
		var buffer = new(bytes.Buffer)
		pp.Write(buffer)
	}
}

func BenchmarkPublish_Marshal(b *testing.B) {
	p := Publish{
		header: Header{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}
	var buffer = make([]byte, p.Size())
	for i := 0; i < b.N; i++ {
		p.MarshalTo(buffer)
	}
}

func BenchmarkPublish_Unmarshal(b *testing.B) {
	b.ReportAllocs()
	p := Publish{
		header: Header{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}
	data, _ := p.Marshal()
	data = data[p.header.Size():]
	for i := 0; i < b.N; i++ {
		var nn Publish
		nn.Unmarshal(data)
	}
}

func BenchmarkPahoUnpack(b *testing.B) {
	b.ReportAllocs()
	pp := paho.PublishPacket{
		FixedHeader: paho.FixedHeader{
			MessageType:     PublishType,
			Dup:             true,
			Qos:             2,
			Retain:          true,
			RemainingLength: 0,
		},
		TopicName: "mqtt-client/test",
		MessageID: 1024,
		Payload:   []byte(strings.Repeat("hello", 1024)),
	}
	var buffer = new(bytes.Buffer)
	pp.Write(buffer)

	for i := 0; i < b.N; i++ {
		_,err := paho.ReadPacket(bytes.NewReader(buffer.Bytes()))
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}
