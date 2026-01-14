package db

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

// DB represents a BadgerDB database
type DB struct {
	db *badger.DB
}

// Config holds database configuration
type Config struct {
	Path string // Database directory path
}

// NewDB creates a new database connection
func NewDB(cfg Config) (*DB, error) {
	opts := badger.DefaultOptions(cfg.Path)

	// Disable default logger for cleaner output
	opts.Logger = nil

	// Tune for Raspberry Pi 3 (1GB RAM)
	// Use default options for v4 which are already optimized
	opts.InMemory = false

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &DB{db: db}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// Set stores a value by key
func (d *DB) Set(key, value []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Get retrieves a value by key
func (d *DB) Get(key []byte) ([]byte, error) {
	var val []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete removes a key
func (d *DB) Delete(key []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// SetJSON stores a JSON-encoded value
func (d *DB) SetJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return d.Set([]byte(key), data)
}

// GetJSON retrieves and decodes a JSON value
func (d *DB) GetJSON(key string, dest interface{}) error {
	data, err := d.Get([]byte(key))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// ListPrefix returns all keys with a given prefix
func (d *DB) ListPrefix(prefix string) ([][]byte, error) {
	var keys [][]byte
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := make([]byte, len(item.Key()))
			copy(key, item.Key())
			keys = append(keys, key)
		}
		return nil
	})
	return keys, err
}

// ListPrefixJSON returns all JSON values with a given prefix
func (d *DB) ListPrefixJSON(prefix string) ([]json.RawMessage, error) {
	var values []json.RawMessage
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			values = append(values, json.RawMessage(val))
		}
		return nil
	})
	return values, err
}

// RunGC manually triggers garbage collection
// Call this periodically to reclaim disk space
func (d *DB) RunGC() error {
	return d.db.RunValueLogGC(0.5)
}
