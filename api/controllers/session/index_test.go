package session

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetSessions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("sqlite3.test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("[API] [TEST] failed to connect database %s", err)
	}
	// db.LogMode(true)
	defer closeDb(db)
	r := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sc := New(db)
	recorder := httptest.NewRecorder()
	params := httprouter.Params{}
	sc.GetSessions(recorder, r, params)
	if recorder.Code != 200 {
		t.Errorf("expecting code 200, received %d", recorder.Code)
	}
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Println("could not close DB", err)
		return
	}
	sql.Close()
}
