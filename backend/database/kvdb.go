package database

import (
	"github.com/spf13/viper"
	"github.com/xujiajun/nutsdb"
)

func InitKvDB() error {
	opt := nutsdb.DefaultOptions
	opt.Dir = viper.GetString("kvdb.path")
	db, err := nutsdb.Open(opt)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func GetKvDB() *nutsdb.DB {
	opt := nutsdb.DefaultOptions
	opt.Dir = viper.GetString("kvdb.path")
	db, _ := nutsdb.Open(opt)
	return db
}
