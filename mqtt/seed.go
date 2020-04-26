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

	createSubscriptions(db)
	createRetains(db)
	createAuth(db)
}

func createSubscriptions(db *sql.DB) {
	sqlStmt := `
	create table subscriptions (
		topic text,
		clientid text,
		qos integer,
		retain_handling integer,
		retain_as_published integer,
		no_local integer,
		enabled integer,
		created_at integer
	);
	create index subscribed_topics_idx on subscriptions(topic);
	create index subscribers_idx on subscriptions(clientid);
	create unique index topic_client_sub_idx on subscriptions(topic, clientid);
	delete from subscriptions;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func createRetains(db *sql.DB) {
	sqlStmt := `
	create table retains (
		topic text,
		application_message blob,
		created_at integer
	);
	create unique index topic_retain_idx on retains(topic);
	delete from retains;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func createAuth(db *sql.DB) {
	sqlStmt := `
	create table auths (
		clientid text,
		username text,
		password blob,
		subscribe_acl text,
		publish_acl text,
		created_at integer
	);
	create unique index clientid_idx on auths(clientid);
	create index clientid_username_idx on auths(clientid, username);
	delete from auths;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
