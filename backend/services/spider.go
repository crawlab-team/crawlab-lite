package services

import (
	"crawlab-lite/model"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetSpiderList() ([]*model.Spider, error) {
	var spiderList []*model.Spider

	spiderDir := viper.GetString("spider.path")
	dirs, _ := ioutil.ReadDir(spiderDir)
	for _, dir := range dirs {
		if dir.IsDir() {
			data, rErr := ioutil.ReadFile(filepath.Join(spiderDir, dir.Name(), "spider.json"))
			if rErr != nil {
				continue
			}
			var spider *model.Spider
			jErr := json.Unmarshal(data, &spider)
			if jErr != nil {
				continue
			}
			spiderList = append(spiderList, spider)
		}
	}

	return spiderList, nil
}

func GeySpiderByName(name string) (*model.Spider, error) {
	spiderDir := viper.GetString("spider.path")
	dirs, _ := ioutil.ReadDir(spiderDir)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == name {
			data, rErr := ioutil.ReadFile(filepath.Join(spiderDir, dir.Name(), "spider.json"))
			if rErr != nil {
				return nil, errors.New("read spider error")
			}
			var spider *model.Spider
			jErr := json.Unmarshal(data, &spider)
			if jErr != nil {
				return nil, errors.New("read spider error")
			}
			return spider, nil
		}
	}
	return nil, nil
}

func DeleteSpiderByName(name string) (interface{}, error) {
	spiderDir := viper.GetString("spider.path")
	dirs, _ := ioutil.ReadDir(spiderDir)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == name {
			return nil, os.RemoveAll(filepath.Join(spiderDir, dir.Name()))
		}
	}
	return nil, nil
}
