package api

import (
	"log"
	"net/http"

	"github.com/ilgianlu/tagyou/api/controllers"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func StartApi(httpPort string) {
	db, err := gorm.Open("sqlite3", "sqlite.db3")
	if err != nil {
		log.Fatal("failed to connect database")
	}
	defer db.Close()

	r := httprouter.New()
	uc := controllers.NewAuthController(db)

	r.GET("/auths", uc.GetAuths)
	r.GET("/auths/:id", uc.GetAuth)
	r.POST("/auths", uc.CreateAuth)
	r.DELETE("/auths/:id", uc.RemoveAuth)

	log.Printf("http listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Panic(err)
	}
}
