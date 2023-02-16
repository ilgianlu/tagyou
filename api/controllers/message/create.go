package message

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

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

	// mc.mqttClient.Publish(mess.Topic, mess.Qos, mess.Retained, payloadFromPayloadType(mess.Payload, mess.PayloadType))

	if res, err := json.Marshal("message published"); err != nil {
		log.Printf("error marshaling response message: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		numBytes, err := w.Write(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Wrote %d bytes json result\n", numBytes)
	}
}

// func payloadFromPayloadType(payload string, payloadType byte) []byte {
// 	return []byte(payload)
// }
