package mqtt

import (
	"database/sql"
	"log"
	"os"
)

func Seed(filename string) {
	os.Remove(filename)

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// SUBSCRIPTIONS
	sqlStmt := `
	create table subscriptions (
		topic text,
		clientid text,
		qos integer,
		retain_handling integer,
		retain_as_published integer
		no_local integer,
		enabled integer
	);
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
