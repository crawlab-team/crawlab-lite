package dao

import (
	bolt "go.etcd.io/bbolt"
)

type Tx struct {
	tx *bolt.Tx
}

// 读事务，同一时间允许多个读事务执行
func ReadTx(db *bolt.DB, fn func(tx Tx) error) (err error) {
	if err = db.View(func(tx *bolt.Tx) error {
		var t Tx
		t.tx = tx
		return fn(t)
	}); err != nil {
		return err
	}
	return nil
}

// 写事务，同一时间只允许一个写事务执行
func WriteTx(db *bolt.DB, fn func(tx Tx) error) (err error) {
	if err = db.Update(func(tx *bolt.Tx) error {
		var t Tx
		t.tx = tx
		return fn(t)
	}); err != nil {
		return err
	}
	return nil
}
