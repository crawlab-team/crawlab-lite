package database

import (
	"crawlab-lite/utils"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

var KvDB *bolt.DB

func InitKvDB() error {
	path := viper.GetString("kvdb.path")
	if utils.PathExist(path) == false {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	db, err := bolt.Open(filepath.Join(path, "kv.db"), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	KvDB = db
	return nil
}
