package ec5

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/ec5/kura"
	"github.com/ilgianlu/tagyou/packet"
	kgo "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

var EC5_TOPIC_FILTER string = "dev-auto"

type metric struct {
	ValueType string      `json:"valueType"`
	Value     interface{} `json:"value"`
}

type metricPayload struct {
	Metrics map[string]metric `json:"metrics"`
}

type kuraJson struct {
	Channel string        `json:"channel"`
	Payload metricPayload `json:"payload"`
}

func StartKafka(host string, topic string, partition int) (*kgo.Writer, error) {
	if !conf.KAFKA_ON {
		return nil, nil
	}
	w := &kgo.Writer{
		Addr:     kgo.TCP(host),
		Topic:    topic,
		Balancer: &kgo.LeastBytes{},
	}
	return w, nil
}

func StopKafka(writer *kgo.Writer) {
	if !conf.KAFKA_ON {
		return
	}
	if err := writer.Close(); err != nil {
		log.Fatal().Err(err).Msg("[KAFKA] failed to close writer")
	}
}

func Publish(writer *kgo.Writer, topic string, p *packet.Packet) {
	if len(topic) <= len(EC5_TOPIC_FILTER) {
		return
	}
	if topic[:len(EC5_TOPIC_FILTER)] != EC5_TOPIC_FILTER {
		return
	}
	prepared, _ := preparePacket(topic, p)
	err := writer.WriteMessages(context.Background(), kgo.Message{Value: prepared})
	if err != nil {
		log.Fatal().Err(err).Msg("[KAFKA] failed to write messages")
	}
}

func preparePacket(topic string, p *packet.Packet) ([]byte, error) {
	decoded := kura.KuraPayload{}
	err := proto.Unmarshal(p.Payload(), &decoded)
	if err != nil {
		return []byte{}, err
	}

	kp := kuraJson{
		Channel: topic,
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
