package db

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	// bucketName is the name of the bucket where all keys/values will be stored within boltDB
	bucketName string
	database   *bolt.DB
}

func NewBoltDB(path string, mode os.FileMode, buketName string) (*BoltDB, error) {
	db, err := bolt.Open(path, mode, nil)
	if err != nil {
		return nil, err
	}

	// Make sure the bucket we are going to use exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(buketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BoltDB{
		bucketName: buketName,
		database:   db,
	}, nil
}

func (d *BoltDB) Insert(shortURL, longURL string) error {
	d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucketName))
		err := b.Put([]byte(shortURL), []byte(longURL))
		return err
	})
	return nil
}

func (d *BoltDB) GetFullURL(shortURL string) (longURL string, err error) {
	d.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucketName))
		longURL = string(b.Get([]byte(shortURL)))
		return nil
	})
	return
}

func (d *BoltDB) Delete(shortURL string) error {
	err := d.database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(d.bucketName))
		return b.Delete([]byte(shortURL))
	})
	return err
}

func (d *BoltDB) Close() error {
	return d.database.Close()
}
