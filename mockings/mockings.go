package mockings

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MockToken struct {
	ready bool
	err   error
}

func (mt *MockToken) Wait() bool {
	return mt.ready
}

func (mt *MockToken) WaitTimeout(time.Duration) bool {
	return mt.ready
}

func (mt *MockToken) Error() error {
	return mt.err
}

type MockClient struct {
	isConnected bool
	token       MQTT.Token
	Published   int
}

func (mc MockClient) IsConnected() bool {
	return mc.isConnected
}

func (mc MockClient) IsConnectionOpen() bool {
	return mc.isConnected
}

func (mc MockClient) Connect() MQTT.Token {
	return mc.token
}

func (mc MockClient) Disconnect(quiesce uint) {
}

func (mc MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) MQTT.Token {
	mc.Published = mc.Published + 1
	return mc.token
}

func (mc MockClient) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) MQTT.Token {
	return mc.token
}

func (mc MockClient) SubscribeMultiple(filters map[string]byte, callback MQTT.MessageHandler) MQTT.Token {
	return mc.token
}

func (mc MockClient) Unsubscribe(topics ...string) MQTT.Token {
	return mc.token
}

func (mc MockClient) AddRoute(topic string, callback MQTT.MessageHandler) {
}

func (mc MockClient) OptionsReader() MQTT.ClientOptionsReader {
	return MQTT.ClientOptionsReader{}
}

func (mc MockClient) GetPublished() int {
	return mc.Published
}
