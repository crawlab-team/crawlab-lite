package model

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Spider struct {
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	File string `json:"file"`
}

func (s *Spider) GetSpiderDirPath() string {
	return filepath.Join(viper.GetString("spider.path"), s.Name)
}

func (s *Spider) GetSpiderFilePath() string {
	return filepath.Join(s.GetSpiderDirPath(), s.File)
}

func (s *Spider) SaveJson() error {
	path := filepath.Join(s.GetSpiderDirPath(), "spider.json")
	data, _ := json.Marshal(s)
	return ioutil.WriteFile(path, data, os.ModePerm)
}
