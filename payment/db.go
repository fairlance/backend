package payment

import (
	"database/sql"
	"log"
	"time"
)

const (
	createEventsIfNotExistSQL = `CREATE TABLE IF NOT EXISTS events
    (
        id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES transactions(id),
        type VARCHAR(255) NOT NULL,
        provider_status VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE
    )`
	createTransactionsIfNotExistSQL = `CREATE TABLE IF NOT EXISTS transactions
    (
        id SERIAL PRIMARY KEY,
        project_id INTEGER NOT NULL,
        track_id VARCHAR(127) UNIQUE NOT NULL,
        provider VARCHAR(255) NOT NULL DEFAULT '',
        provider_transaction_key VARCHAR(255) NOT NULL,
        provider_status VARCHAR(255) NOT NULL,
        amount VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
		error_msg VARCHAR(255) NOT NULL DEFAULT '',
        created_at TIMESTAMP WITH TIME ZONE,
        updated_at TIMESTAMP WITH TIME ZONE
    )`
	createReceiversIfNotExistSQL = `CREATE TABLE IF NOT EXISTS receivers
    (
        id SERIAL PRIMARY KEY,
		transaction_id INTEGER NOT NULL REFERENCES transactions(id),
		fairlance_id INTEGER NOT NULL,
        provider_status VARCHAR(255) NOT NULL,
        provider_transaction_key VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		amount VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE,
        updated_at TIMESTAMP WITH TIME ZONE
	)`

	insertEventSQL                               = `INSERT INTO events (transaction_id, type, provider_status, status, created_at) VALUES ($1,$2,$3,$4,$5)`
	insertTransactionSQL                         = `INSERT INTO transactions (track_id, provider, provider_transaction_key, provider_status, project_id, amount, status, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	insertReceiverSQL                            = `INSERT INTO receivers (transaction_id, provider_transaction_key, provider_status, fairlance_id, email, amount, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	updateTransactionSQL                         = `UPDATE transactions SET provider=$1,provider_transaction_key=$2, provider_status=$3, status=$4, error_msg=$5, updated_at=$6 WHERE id=$7`
	updateReceiverSQL                            = `UPDATE receivers SET provider_transaction_key=$1, provider_status=$2, updated_at=$3 WHERE id=$4`
	selectTransactionByProjectIDSQL              = `SELECT id, track_id, provider, provider_transaction_key, provider_status, project_id, amount, status, error_msg, created_at, updated_at FROM transactions WHERE project_id = $1`
	selectTransactionByProviderTransactionKeySQL = `SELECT id, track_id, provider, provider_transaction_key, provider_status, project_id, amount, status, error_msg, created_at, updated_at FROM transactions WHERE provider_transaction_key = $1`
	selectReceiverSQL                            = `SELECT id, fairlance_id, email, amount FROM receivers WHERE transaction_id = $1`

	eventTypeTransaction = "transaction"
	eventTypeReceiver    = "receiver"
)

type DB interface {
	Init()
	Insert(t *Transaction) error
	UpdateTransaction(t *Transaction) error
	UpdateReceiver(r *TransactionReceiver) error
	GetByProjectID(projectID uint) (*Transaction, error)
	GetByProviderTransactionKey(providerTransactionKey string) (*Transaction, error)
}

func NewDB(db *sql.DB) DB {
	return &sqlDB{db}
}

type sqlDB struct {
	storage *sql.DB
}

func (db *sqlDB) Init() {
	if err := db.storage.Ping(); err != nil {
		log.Fatalf("could not ping db: %v", err)
	}
	if _, err := db.storage.Exec(createTransactionsIfNotExistSQL); err != nil {
		log.Fatalf("could not create  transactions table: %v", err)
	}
	if _, err := db.storage.Exec(createReceiversIfNotExistSQL); err != nil {
		log.Fatalf("could not create receivers table: %v", err)
	}
	if _, err := db.storage.Exec(createEventsIfNotExistSQL); err != nil {
		log.Fatalf("could not create events table: %v", err)
	}
}

func (db *sqlDB) Insert(t *Transaction) error {
	txn, err := db.storage.Begin()
	if err != nil {
		return err
	}
	now := time.Now()
	var transactionID uint
	if err = txn.QueryRow(insertTransactionSQL, t.TrackID, t.Provider, t.ProviderTransactionKey, t.ProviderStatus, t.ProjectID, t.Amount, t.Status, now, now).Scan(&transactionID); err != nil {
		txn.Rollback()
		return err
	}
	for _, receiver := range t.Receivers {
		if _, err := txn.Exec(insertReceiverSQL, transactionID, "", "", receiver.FairlanceID, receiver.ProviderIdentifier, receiver.Amount, now, now); err != nil {
			txn.Rollback()
			return err
		}
	}
	if err = txn.Commit(); err != nil {
		return err
	}
	if _, err := db.storage.Exec(insertEventSQL, transactionID, eventTypeTransaction, t.ProviderStatus, t.Status, now); err != nil {
		log.Printf("could not create event for transaction %d: %v", t.ID, err)
	}
	return nil
}

func (db *sqlDB) UpdateTransaction(t *Transaction) error {
	now := time.Now()
	if _, err := db.storage.Exec(updateTransactionSQL,
		t.Provider, t.ProviderTransactionKey, t.ProviderStatus, t.Status, t.ErrorMsg, now, t.ID,
	); err != nil {
		return err
	}
	if _, err := db.storage.Exec(insertEventSQL, t.ID, eventTypeTransaction, t.ProviderStatus, t.Status, now); err != nil {
		log.Printf("could not create event for transaction %d: %v", t.ID, err)
	}
	return nil
}

func (db *sqlDB) UpdateReceiver(r *TransactionReceiver) error {
	now := time.Now()
	if _, err := db.storage.Exec(updateReceiverSQL,
		r.ProviderTransactionKey, r.ProviderStatus, now, r.ID,
	); err != nil {
		return err
	}
	if _, err := db.storage.Exec(insertEventSQL, r.ID, eventTypeReceiver, r.ProviderStatus, "", now); err != nil {
		log.Printf("could not create event for receiver %d: %v", r.ID, err)
	}
	return nil
}

func (db *sqlDB) GetByProjectID(projectID uint) (*Transaction, error) {
	var t Transaction
	if err := db.storage.QueryRow(selectTransactionByProjectIDSQL, projectID).Scan(
		&t.ID, &t.TrackID, &t.Provider, &t.ProviderTransactionKey, &t.ProviderStatus, &t.ProjectID, &t.Amount, &t.Status, &t.ErrorMsg, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	rows, err := db.storage.Query(selectReceiverSQL, t.ID)
	if err != nil {
		log.Printf("could not get receivers from db: %v", t)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var receiver TransactionReceiver
		if err := rows.Scan(&receiver.ID, &receiver.FairlanceID, &receiver.ProviderIdentifier, &receiver.Amount); err != nil {
			return nil, err
		}
		t.Receivers = append(t.Receivers, receiver)
	}
	return &t, nil
}

func (db *sqlDB) GetByProviderTransactionKey(providerTransactionKey string) (*Transaction, error) {
	var t Transaction
	if err := db.storage.QueryRow(selectTransactionByProviderTransactionKeySQL, providerTransactionKey).Scan(
		&t.ID, &t.TrackID, &t.Provider, &t.ProviderTransactionKey, &t.ProviderStatus, &t.ProjectID, &t.Amount, &t.Status, &t.ErrorMsg, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, err
	}
	rows, err := db.storage.Query(selectReceiverSQL, t.ID)
	if err != nil {
		log.Printf("could not get receivers from db: %v", t)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var receiver TransactionReceiver
		if err := rows.Scan(&receiver.ID, &receiver.FairlanceID, &receiver.ProviderIdentifier, &receiver.Amount); err != nil {
			return nil, err
		}
		t.Receivers = append(t.Receivers, receiver)
	}
	return &t, nil
}
