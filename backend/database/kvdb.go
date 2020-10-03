package database

import (
	"crawlab-lite/utils"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

var MainDB, LogDB *bolt.DB

func InitKvDB() error {
	path := viper.GetString("kvdb.path")
	if utils.PathExist(path) == false {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	mainDB, err := bolt.Open(filepath.Join(path, "main.db"), 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}
	MainDB = mainDB

	logDB, err := bolt.Open(filepath.Join(path, "log.db"), 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}
	LogDB = logDB

	return nil
}
