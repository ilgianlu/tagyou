package session

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func TestGetSessions(t *testing.T) {
	dbConn, err := sql.Open("sqlite3", "sqlite3.test.db3")
	if err != nil {
		t.Errorf("[API] [TEST] failed to connect database %s", err)
	}
	db := dbaccess.New(dbConn)

	p := persistence.SqlPersistence{}
	p.InnerInit(db, false, false, []byte(""))

	// db.LogMode(true)
	defer closeDb(dbConn)
	r := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sc := New()
	recorder := httptest.NewRecorder()
	params := httprouter.Params{}
	sc.GetSessions(recorder, r, params)
	if recorder.Code != 200 {
		t.Errorf("expecting code 200, received %d", recorder.Code)
	}
}

func closeDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		slog.Info("could not close DB")
	}
}
