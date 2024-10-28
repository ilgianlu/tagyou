package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilgianlu/tagyou/persistence"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetSessions(t *testing.T) {
	p := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	p.Init(false, false, []byte{})

	r := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sc := NewController()
	recorder := httptest.NewRecorder()
	sc.GetSessions(recorder, r)
	if recorder.Code != 200 {
		t.Errorf("expecting code 200, received %d", recorder.Code)
	}
}
