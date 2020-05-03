package database

import (
	"github.com/spf13/viper"
	"github.com/xujiajun/nutsdb"
)

var KvDB *nutsdb.DB

func InitKvDB() error {
	opt := nutsdb.DefaultOptions
	opt.SegmentSize = 64 * 1000 * 1000
	opt.Dir = viper.GetString("kvdb.path")
	db, err := nutsdb.Open(opt)
	if err != nil {
		return err
	}
	KvDB = db
	return nil
}
