package mqtt

import (
	"database/sql"
	"log"
	"os"
)

func Seed() {
	os.Remove(os.Getenv("DB_FILE"))

	db, err := sql.Open("sqlite3", os.Getenv("DB_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table subscriptions (topic text, clientid text);
	create index subscribed_topics on subscriptions(topic);
	create index subscribers on subscriptions(clientid);
	delete from subscriptions;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
