package payment

import (
	"database/sql"
	"log"
	"time"
)

const (
	createIfNotExistSQL = ` CREATE TABLE IF NOT EXISTS transactions
    (
        id SERIAL PRIMARY KEY,
        track_id varchar(127) UNIQUE NOT NULL,
        provider varchar(255) NOT NULL,
        provider_key varchar(255) NOT NULL,
        project_id integer NOT NULL,
        amount varchar(255) NOT NULL,
        status varchar(255) NOT NULL,
        receivers JSONB NOT NULL,
        created_at date,
        updated_at date
    )`
	insertSQL = `INSERT INTO transactions (track_id, provider, provider_key, project_id, amount, status, receivers, created_at, updated_at)  
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	selectSQL = `SELECT track_id, provider, provider_key, project_id, amount, status, receivers, created_at, updated_at FROM transactions WHERE track_id = $1`
	updateSQL = `UPDATE transactions SET track_id=$1, provider=$2, provider_key=$3, project_id=$4, amount=$5, status=$6, receivers=$7, created_at=$8, updated_at=$9 WHERE track_id=$10`
)

type db interface {
	init()
	get(trackID string) (transaction, error)
	insert(t transaction) error
	update(t transaction) error
}

func newDB(db *sql.DB) db {
	return &sqlDB{db}
}

type sqlDB struct {
	storage *sql.DB
}

func (db *sqlDB) init() {
	if err := db.storage.Ping(); err != nil {
		log.Fatalf("could not ping db: %v", err)
	}
	_, err := db.storage.Exec(createIfNotExistSQL)
	if err != nil {
		log.Fatalf("could not create table: %v", err)
	}
}

func (db *sqlDB) insert(t transaction) error {
	now := time.Now()
	_, err := db.storage.Exec(insertSQL,
		t.trackID,
		t.provider,
		t.providerKey,
		t.projectID,
		t.amount,
		t.status,
		t.receivers,
		&now, // createdAt
		&now, // updatedAt
	)
	return err
}

func (db *sqlDB) get(trackID string) (transaction, error) {
	var t transaction
	if err := db.storage.QueryRow(selectSQL, trackID).Scan(
		&t.trackID,
		&t.provider,
		&t.providerKey,
		&t.projectID,
		&t.amount,
		&t.status,
		&t.receivers,
		&t.createdAt,
		&t.updatedAt,
	); err != nil {
		return transaction{}, err
	}
	return t, nil
}

func (db *sqlDB) update(t transaction) error {
	_, err := db.storage.Exec(updateSQL,
		t.trackID,
		t.provider,
		t.providerKey,
		t.projectID,
		t.amount,
		t.status,
		t.receivers,
		t.createdAt,
		time.Now(),
		t.trackID,
	)
	return err
}
