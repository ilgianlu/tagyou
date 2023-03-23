package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetSessions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("sqlite3.test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Msgf("[API] [TEST] failed to connect database %s", err)
	}
	p := persistence.SqlPersistence{}
	p.InnerInit(db, false, false)
	// db.LogMode(true)
	defer closeDb(db)
	r := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sc := New()
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
		log.Err(err).Msg("could not close DB")
		return
	}
	sql.Close()
}
