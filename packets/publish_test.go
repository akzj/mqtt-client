package packets

import (
	"bytes"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang/packets"
	"reflect"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	p := Publish{
		Header: Header{
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
	np.Header = p.Header

	np.Unmarshal(data[np.Header.Size():])

	data2, _ := np.Marshal()

	if bytes.Compare(buffer.Bytes(), data2) != 0 {
		t.Error("publish marshal failed")
		fmt.Println(buffer.Bytes())
		fmt.Println(data2)
	}

	p2, d, err := UnmarshalPacket(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(d, data) == false {
		fmt.Println(data)
		fmt.Println(d)
		t.Fatal("UnmarshalPacket failed", len(data), len(data))
	}
	if reflect.DeepEqual(&p, p2) == false {
		fmt.Println(p)
		fmt.Println(p2)
		t.Fatal("UnmarshalPacket failed")
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
		Header: Header{
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
		Header: Header{
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
	reader := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		_, _, err := UnmarshalPacket(reader)
		if err != nil {
			b.Fatal("UnmarshalPacket failed")
		}
		reader.Reset(data)
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

	reader := bytes.NewReader(buffer.Bytes())
	for i := 0; i < b.N; i++ {
		_, err := paho.ReadPacket(reader)
		if err != nil {
			b.Fatalf(err.Error())
		}
		reader.Reset(buffer.Bytes())
	}
}
