package payment

import (
	"database/sql"
	"log"
	"time"
)

const (
	createTransactionsIfNotExistSQL = `CREATE TABLE IF NOT EXISTS transactions
    (
        id SERIAL PRIMARY KEY,
        track_id VARCHAR(127) UNIQUE NOT NULL,
        provider VARCHAR(255) NOT NULL,
        provider_transaction_key VARCHAR(255) NOT NULL,
        project_id INTEGER NOT NULL,
        amount VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE,
        updated_at TIMESTAMP WITH TIME ZONE
    )`
	createReceiversIfNotExistSQL = `CREATE TABLE IF NOT EXISTS receivers
    (
        id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES transactions(id),
		fairlance_id INTEGER NOT NULL,
		email VARCHAR(255) NOT NULL,
		amount VARCHAR(255) NOT NULL
    )`
	insertReceiverSQL    = `INSERT INTO receivers (transaction_id, fairlance_id, email, amount) VALUES ($1,$2,$3,$4)`
	insertTransactionSQL = `INSERT INTO transactions (track_id, provider, provider_transaction_key, project_id, amount, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	// selectTransactionSQL = `SELECT id, track_id, provider, provider_transaction_key, project_id, amount, status, created_at, updated_at FROM transactions WHERE track_id = $1`
	// selectReceiverSQL    = `SELECT id, fairlance_id, email, amount FROM receivers WHERE transaction_id = $1`
	updateTransactionPaymentTransactionKeySQL = `UPDATE transactions SET status=$1,provider_transaction_key=$2, updated_at=$3 WHERE track_id=$4`
	updateTransactionSatatusSQL               = `UPDATE transactions SET status=$1, updated_at=$2 WHERE track_id=$3`
)

type db interface {
	init()
	// get(trackID string) (transaction, error)
	insert(t *transaction) error
	updatePaymentKeyAndStatusByTrackID(trackID, paymentTransactionKey, status string) error
	updateStatusByTractID(trackID, status string) error
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
	if _, err := db.storage.Exec(createTransactionsIfNotExistSQL); err != nil {
		log.Fatalf("could not create  transactions table: %v", err)
	}
	if _, err := db.storage.Exec(createReceiversIfNotExistSQL); err != nil {
		log.Fatalf("could not create receivers table: %v", err)
	}
}

func (db *sqlDB) insert(t *transaction) error {
	txn, err := db.storage.Begin()
	if err != nil {
		return err
	}
	now := time.Now()
	var transactionID uint
	err = db.storage.QueryRow(insertTransactionSQL, t.trackID, t.provider, t.paymentKey, t.projectID, t.amount, t.status, &now, &now).Scan(&transactionID)
	if err != nil {
		return err
	}
	for _, receiver := range t.receivers {
		_, err = txn.Exec(insertReceiverSQL, transactionID, receiver.fairlanceID, receiver.email, receiver.amount)
		if err != nil {
			return err
		}
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *sqlDB) updatePaymentKeyAndStatusByTrackID(trackID, paymentTransactionKey, status string) error {
	_, err := db.storage.Exec(updateTransactionPaymentTransactionKeySQL,
		status, paymentTransactionKey, time.Now(), trackID,
	)
	return err
}

func (db *sqlDB) updateStatusByTractID(trackID, status string) error {
	_, err := db.storage.Exec(updateTransactionSatatusSQL,
		status, time.Now(), trackID,
	)
	return err
}

// func (db *sqlDB) get(trackID string) (transaction, error) {
// 	var t transaction
// 	if err := db.storage.QueryRow(selectTransactionSQL, trackID).Scan(
// 		&t.id, &t.trackID, &t.provider, &t.paymentKey, &t.projectID, &t.amount, &t.status, &t.createdAt, &t.updatedAt,
// 	); err != nil {
// 		return transaction{}, err
// 	}
// 	rows, err := db.storage.Query(selectReceiverSQL, t.id)
// 	if err != nil {
// 		log.Printf("%v", t)
// 		return transaction{}, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var receiver paymentReceiver
// 		err := rows.Scan(&receiver.id, &receiver.fairlanceID, &receiver.email, &receiver.amount)
// 		if err != nil {
// 			return transaction{}, err
// 		}
// 		t.receivers = append(t.receivers, receiver)
// 	}
// 	return t, nil
// }
