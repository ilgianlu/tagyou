package mqtt

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const RETRY_INSERT = `insert into retries(
	clientid, application_message, packet_identifier,
	qos, retries, ack_status, created_at)
	values(?, ?, ?, ?, ?, ?, ?)`
const RETRY_DELETE = "delete from retries where clientid = ? and packet_identifier = ?"
const RETRY_SELECT_BY_CLIENTID_PACKETIDENT = "select * from retries where clientid = ? and packet_identifier = ?"

type SqliteRetries struct {
	db *sql.DB
}

func (is SqliteRetries) addRetry(retry Retry) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare(RETRY_INSERT)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		retry.clientId, retry.applicationMessage, retry.packetIdentifier,
		retry.qos, retry.retries, retry.ackStatus, retry.createdAt.Unix())
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetries) remRetry(clientId string, packetIdentifier int) error {
	tx, err := is.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare(RETRY_DELETE)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientId, packetIdentifier)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

func (is SqliteRetries) findRetriesByClientId(clientId string, packetIdentifier int) []Retry {
	retries := []Retry{}
	rows, err := is.db.Query(RETRY_SELECT_BY_CLIENTID_PACKETIDENT, clientId, packetIdentifier)
	if err != nil {
		log.Println(err)
		return retries
	}
	defer rows.Close()
	for rows.Next() {
		var r Retry
		err = rows.Scan(
			&r.clientId, &r.applicationMessage, &r.packetIdentifier,
			&r.qos, &r.retries, &r.ackStatus)
		if err != nil {
			log.Println(err)
			return retries
		}
		retries = append(retries, r)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return retries
	}

	return retries
}
