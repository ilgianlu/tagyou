package api

import (
	"log"
	"net/http"

	AuthController "github.com/ilgianlu/tagyou/api/controllers/auth"
	SessionController "github.com/ilgianlu/tagyou/api/controllers/session"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func StartApi(httpPort string) {
	db, err := gorm.Open("sqlite3", "sqlite.db3")
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// db.LogMode(true)
	defer db.Close()

	r := httprouter.New()
	uc := AuthController.New(db)
	uc.RegisterRoutes(r)
	sc := SessionController.New(db)
	sc.RegisterRoutes(r)

	log.Printf("http listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Panic(err)
	}
}
