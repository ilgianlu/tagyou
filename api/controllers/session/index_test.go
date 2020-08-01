package session

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func TestGetSessions(t *testing.T) {
	db, err := gorm.Open("sqlite3", "./testdata/sqlite.test")
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// db.LogMode(true)
	defer db.Close()
	r := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sc := New(db)
	recorder := httptest.NewRecorder()
	params := httprouter.Params{}
	sc.GetSessions(recorder, r, params)
	if recorder.Code != 200 {
		t.Errorf("expecting code 200, received %d", recorder.Code)
	}
}
