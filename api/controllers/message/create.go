package message

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Message struct {
	Topic       string
	Qos         byte
	Retained    bool
	Payload     string
	PayloadType byte
}

func (mc MessageController) PostMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	mess := Message{}
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		log.Printf("error decoding json input: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mc.mqttClient.Publish(mess.Topic, mess.Qos, mess.Retained, payloadFromPayloadType(mess.Payload, mess.PayloadType))
}

func payloadFromPayloadType(payload string, payloadType byte) []byte {
	return []byte(payload)
}
