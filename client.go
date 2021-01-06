package mqttclient

import (
	"bufio"
	"bytes"
	"github.com/akzj/block-queue"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"log"
	"net"
	"reflect"
	"sync"
)

type Client struct {
	err           error
	messageID     uint16
	conn          net.Conn
	pubQueue      blockqueue.QueueWithContext
	callbackQueue blockqueue.QueueWithContext

	pubRequestsMtx sync.Mutex
	pubRequests    map[uint16]*pubRequest
}

type Callback func(err error)

type connectRequest struct {
	packets.ConnectPacket
	callback Callback
}

type callback struct {
	err    error
	handle Callback
}

type pubRequest struct {
	packets.PublishPacket
	callback Callback
}

func (client *Client) readLoop() {
	var reader = bufio.NewReaderSize(client.conn, 64*1024)
	for {
		packet, err := packets.ReadPacket(reader)
		if err != nil {
			client.err = err
			client.Close()
			return
		}
		switch packet := packet.(type) {
		case *packets.PubackPacket:
			client.handlePubackPacket(packet)
		}
	}
}

func (client *Client) handlePubackPacket(packet *packets.PubackPacket) {
	client.pubRequestsMtx.Lock()
	request, ok := client.pubRequests[packet.MessageID]
	if ok == false {
		client.pubRequestsMtx.Unlock()
		log.Printf("no find messageID with %d \n", packet.MessageID)
		return
	}
	delete(client.pubRequests, packet.MessageID)
	client.pubRequestsMtx.Unlock()

	client.callbackQueue.Push(callback{
		err:    nil,
		handle: request.callback,
	})
}

func (client *Client) writeLoop() {
	for {
		buffer := bytes.NewBuffer(make([]byte, 0, 64*1024))
		items, err := client.pubQueue.PopAll(nil)
		if err != nil {
			return
		}
		for _, item := range items {
			switch obj := item.(type) {
			case *pubRequest:
				messageID := client.getMessageID(obj.Qos)
				if messageID != 0 {
					client.pubRequests[messageID] = obj
				}
				if err := obj.Write(buffer); err != nil {
					panic(err) //todo remove panic
				}

			case *connectRequest:
				if err := obj.ConnectPacket.Write(buffer); err != nil {
					panic(err) //todo remove panic
				}
			default:
				panic("unknown type " + reflect.TypeOf(obj).String()) ////todo remove panic
			}
		}
		if _, err := client.conn.Write(buffer.Bytes()); err != nil {
			client.err = err
			client.Close()
			return
		}
	}
}

func (client *Client) Pub(topic string, Qos byte, retain bool, payload []byte, callback Callback) {
	if client.err != nil {
		go callback(client.err)
		return
	}
	_ = client.pubQueue.Push(&pubRequest{
		PublishPacket: packets.PublishPacket{
			FixedHeader: packets.FixedHeader{
				MessageType:     packets.Publish,
				Dup:             false,
				Qos:             Qos,
				Retain:          retain,
				RemainingLength: 0,
			},
			TopicName: topic,
			MessageID: 0,
			Payload:   payload,
		},
		callback: callback,
	})
}

func (client *Client) getMessageID(qos byte) uint16 {
	if qos == 0 {
		return 0
	}
	client.messageID++
	return client.messageID
}

func (client *Client) Close() {

}
