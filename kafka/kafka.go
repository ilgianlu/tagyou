package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/kafka/kura"
	"github.com/ilgianlu/tagyou/packet"
	kgo "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type metric struct {
	Name      string      `json:"name"`
	ValueType string      `json:"valueType"`
	Value     interface{} `json:"value"`
}

type metricPayload struct {
	Metrics []metric `json:"metrics"`
}

type kuraJson struct {
	Channel string        `json:"channel"`
	Payload metricPayload `json:"payload"`
}

func StartKafka(host string, topic string, partition int) (*kgo.Conn, error) {
	if !conf.KAFKA_ON {
		return nil, nil
	}
	return kgo.DialLeader(context.Background(), "tcp", host, topic, partition)
}

func Publish(conn *kgo.Conn, p *packet.Packet) {

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	prepared, _ := preparePacket(p)
	_, err := conn.WriteMessages(
		kgo.Message{Value: prepared},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("[KAFKA] failed to write messages")
	}

	if err := conn.Close(); err != nil {
		log.Fatal().Err(err).Msg("[KAFKA] failed to close writer")
	}
}

func preparePacket(p *packet.Packet) ([]byte, error) {
	decoded := kura.KuraPayload{}
	err := proto.Unmarshal(p.Payload(), &decoded)
	if err != nil {
		return []byte{}, err
	}

	kp := kuraJson{
		Channel: p.Topic,
	}
	for _, m := range decoded.Metric {
		newMetric := metric{
			Name:      m.GetName(),
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
		kp.Payload.Metrics = append(kp.Payload.Metrics, newMetric)
	}

	toSend, err := json.Marshal(kp)
	if err != nil {
		return []byte{}, err
	}
	log.Debug().Msg(string(toSend))
	return toSend, nil
}
