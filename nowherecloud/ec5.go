package nowherecloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/nowherecloud/kura"
	"github.com/ilgianlu/tagyou/packet"
	kgo "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type metric struct {
	ValueType string      `json:"valueType"`
	Value     interface{} `json:"value"`
}

type metricPayload struct {
	Metrics map[string]metric `json:"metrics"`
}
type messageChannel struct {
	SemanticParts []string `json:"semanticParts"`
}

type kuraJson struct {
	ScopeId  string         `json:"scopeId"`
	DeviceId string         `json:"deviceId"`
	ClientId string         `json:"clientId"`
	Channel  messageChannel `json:"channel"`
	Payload  metricPayload  `json:"payload"`
}

type NowhereConnector struct {
	kWriter *kgo.Writer
}

type NcMessage struct {
	Topic string
	P     *packet.Packet
}

func (nc *NowhereConnector) Init() (chan NcMessage, error) {
	Loader()
	ncMessages := make(chan NcMessage)
	kWriter, err := StartKafka(KAFKA_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("[NOWHERE-CLOUD] failed to connect to kafka")
	}
	nc.kWriter = kWriter
	log.Info().Msg("[NOWHERE-CLOUD] kafka connected")
	go nc.rangeNcMessages(ncMessages)
	return ncMessages, nil
}

func StartKafka(url string) (*kgo.Writer, error) {
	if !KAFKA_ON {
		return nil, nil
	}
	hosts := strings.Split(url, ",")
	w := &kgo.Writer{
		Addr:     kgo.TCP(hosts...),
		Balancer: &kgo.LeastBytes{},
		Async:    true,
	}
	return w, nil
}

func StopKafka(writer *kgo.Writer) {
	if !KAFKA_ON {
		return
	}
	if err := writer.Close(); err != nil {
		log.Fatal().Err(err).Msg("[NOWHERE-CLOUD] failed to close writer")
	}
}

func (nc *NowhereConnector) rangeNcMessages(ncMessage chan NcMessage) {
	for ncMessage := range ncMessage {
		respected, found := respectFilter(ncMessage.Topic)
		if !found {
			return
		}
		prepared, _ := preparePacket(ncMessage.Topic, ncMessage.P)
		log.Debug().Msg(fmt.Sprintf("[NOWHERE-CLOUD] Publishing to %s", respected))
		err := nc.kWriter.WriteMessages(context.Background(), kgo.Message{Topic: respected, Value: prepared})
		if err != nil {
			log.Fatal().Err(err).Msg("[NOWHERE-CLOUD] failed to write messages")
		}
	}
}

func semanticParts(parts []string) []string {
	return parts[2:]
}

func topicParts(topic string) []string {
	return strings.Split(topic, "/")
}

func preparePacket(topic string, p *packet.Packet) ([]byte, error) {
	decoded := kura.KuraPayload{}
	err := proto.Unmarshal(p.Payload(), &decoded)
	if err != nil {
		return []byte{}, err
	}

	parts := topicParts(topic)

	kp := kuraJson{
		ScopeId:  parts[0],
		ClientId: parts[1],
		Channel: messageChannel{
			SemanticParts: semanticParts(parts),
		},
	}
	kp.Payload.Metrics = make(map[string]metric)
	for _, m := range decoded.Metric {
		newMetric := metric{
			ValueType: kura.KuraPayload_KuraMetric_ValueType_name[int32(m.GetType())],
		}
		switch newMetric.ValueType {
		case "STRING":
			{
				newMetric.Value = m.GetStringValue()
			}
		case "DOUBLE":
			{
				newMetric.Value = m.GetDoubleValue()
			}
		case "FLOAT":
			{
				newMetric.Value = m.GetFloatValue()
			}
		case "INT64":
			{
				newMetric.Value = m.GetIntValue()
			}
		case "INT32":
			{
				newMetric.Value = m.GetIntValue()
			}
		case "BOOL":
			{
				newMetric.Value = m.GetBoolValue()
			}
		case "BYTES":
			{
				newMetric.Value = m.GetBytesValue()
			}
		}
		kp.Payload.Metrics[m.GetName()] = newMetric
	}

	toSend, err := json.Marshal(kp)
	if err != nil {
		return []byte{}, err
	}
	log.Debug().Msg(string(toSend))
	return toSend, nil
}
