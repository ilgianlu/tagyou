package message

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mockings "github.com/ilgianlu/tagyou/mockings"
	"github.com/julienschmidt/httprouter"
)

func TestPostMessage(t *testing.T) {

	var _ mqtt.Client = (*mockings.MockClient)(nil)

	c := mockings.MockClient{}
	mcPostBody := map[string]interface{}{
		"Topic":       "/a",
		"Qos":         0,
		"Retained":    false,
		"Payload":     "hello world",
		"PayloadType": 0,
	}
	body, _ := json.Marshal(mcPostBody)
	r := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))

	mc := New(&c)
	recorder := httptest.NewRecorder()
	params := httprouter.Params{}
	mc.PostMessage(recorder, r, params)
	if recorder.Code != 200 {
		t.Errorf("expecting code 200, received %d", recorder.Code)
	}

	resultMessage := recorder.Body.String()
	if !strings.EqualFold(resultMessage, "\"message published\"") {
		t.Errorf("expecting \"message published\", received %s", resultMessage)
	}
}
