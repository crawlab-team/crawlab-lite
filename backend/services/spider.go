package services

import (
	"crawlab-lite/constants"
	"crawlab-lite/database"
	"crawlab-lite/forms"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/xujiajun/nutsdb"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func QuerySpiderList(pageNum int, pageSize int) (total int, spiders []*models.Spider, err error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize

	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		// 查询区间内的所有爬虫
		if nodes, err := tx.ZRangeByRank(constants.SpiderListBucket, start, end); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var spider *models.Spider
				if err := json.Unmarshal(node.Value, &spider); err != nil {
					return err
				}
				spiders = append(spiders, spider)
			}
		}

		// 查询数据总数目
		if total, err = tx.ZCard(constants.SpiderListBucket); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, spiders, nil
}

func QuerySpiderByName(name string) (spider *models.Spider, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(constants.SpiderListBucket, []byte(name)); err != nil {
			if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
				return nil
			}
			return err
		} else if err := json.Unmarshal(node.Value, &spider); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return spider, nil
}

func SaveSpider(form forms.SpiderForm) (spider *models.Spider, err error) {
	spiderName := form.Name

	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider != nil {
		return nil, errors.New("spider already exists")
	}

	spider = &models.Spider{
		Name:     spiderName,
		CreateTs: utils.NowTimestamp(),
	}

	db := database.GetKvDB()
	defer db.Close()

	// 存储爬虫信息
	if err := db.Update(func(tx *nutsdb.Tx) error {
		score := float64(spider.CreateTs) / math.Pow10(0)
		value, _ := json.Marshal(&spider)
		return tx.ZAdd(constants.SpiderListBucket, []byte(spiderName), score, value)
	}); err != nil {
		return nil, err
	}

	return spider, nil
}

func RemoveSpider(spiderName string) (res interface{}, err error) {
	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider == nil {
		return nil, errors.New("spider not found")
	}

	// 删除版本文件
	spiderDir := viper.GetString("spider.path")
	dirs, _ := ioutil.ReadDir(spiderDir)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == spiderName {
			if err := os.RemoveAll(filepath.Join(spiderDir, spiderName)); err != nil {
				return nil, err
			}
		}
	}

	db := database.GetKvDB()
	defer db.Close()

	// 删除爬虫相关数据
	if err := db.Update(func(tx *nutsdb.Tx) error {
		// 删除爬虫
		if err := tx.ZRem(constants.SpiderListBucket, spiderName); err != nil {
			return err
		}
		// 删除爬虫版本
		verBucket := getVersionBucket(spiderName)
		if nodes, err := tx.ZMembers(verBucket); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				if err := tx.ZRem(verBucket, node.Key()); err != nil {
					return err
				}
			}
		}
		// 删除爬虫任务
		if nodes, err := tx.ZMembers(constants.TaskListBucket); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var task *models.Task
				if err := json.Unmarshal(node.Value, &task); err != nil {
					return err
				}
				if task.SpiderName == spiderName {
					if err := tx.ZRem(constants.TaskListBucket, task.Id); err != nil {
						return err
					}
					return nil
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func QuerySpiderVersionList(spiderName string) (versions []*models.SpiderVersion, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		// 查询爬虫下的所有版本信息
		if nodes, err := tx.ZMembers(getVersionBucket(spiderName)); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var version *models.SpiderVersion
				if err := json.Unmarshal(node.Value, &version); err != nil {
					return err
				}
				versions = append(versions, version)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return versions, nil
}

func GetSpiderVersionById(spiderName string, versionId string) (version *models.SpiderVersion, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(getVersionBucket(spiderName), []byte(versionId)); err != nil {
			if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
				return nil
			}
			return err
		} else if err := json.Unmarshal(node.Value, &version); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return version, nil
}

func SaveSpiderVersion(spiderName string, form forms.SpiderUploadForm) (version *models.SpiderVersion, err error) {
	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider == nil {
		return nil, errors.New("spider not found")
	}

	// 生成存储版本文件的路径
	fileName := strconv.FormatInt(utils.NowTimestamp(), 10) + ".zip"
	dir := filepath.Join(viper.GetString("spider.path"), spiderName)
	path := filepath.Join(dir, fileName)

	file, err := form.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 保存版本文件
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	if err := utils.SaveFile(file, path); err != nil {
		return nil, err
	}

	version = &models.SpiderVersion{
		Id:       utils.GetFileMD5(file),
		Path:     filepath.Join(spiderName, fileName),
		UploadTs: utils.NowTimestamp(),
	}

	db := database.GetKvDB()
	defer db.Close()

	// 根据 MD5 判断是否为重复的版本
	if version, err := GetSpiderVersionById(spiderName, version.Id); err != nil {
		return nil, err
	} else if version != nil {
		return nil, errors.New("version already exists")
	}

	// 存储版本信息
	if err := db.Update(func(tx *nutsdb.Tx) error {
		spiderBucket := getVersionBucket(spiderName)
		score := float64(version.UploadTs) / math.Pow10(0)
		value, _ := json.Marshal(&version)
		return tx.ZAdd(spiderBucket, []byte(version.Id), score, value)
	}); err != nil {
		return nil, err
	}

	return version, nil
}

func RemoveSpiderVersion(spiderName string, versionId string) (res interface{}, err error) {
	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider == nil {
		return nil, errors.New("spider not found")
	}

	// 查询版本信息
	version, err := GetSpiderVersionById(spiderName, versionId)
	if err != nil {
		return nil, err
	} else if version == nil {
		return nil, errors.New("version not found")
	}

	// 删除版本文件
	dirPath := viper.GetString("spider.path")
	if err := os.Remove(filepath.Join(dirPath, version.Path)); err != nil {
		return nil, err
	}

	db := database.GetKvDB()
	defer db.Close()

	// 删除爬虫与版本数据
	if err := db.Update(func(tx *nutsdb.Tx) error {
		if err := tx.ZRem(getVersionBucket(spiderName), versionId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func getVersionBucket(spiderName string) string {
	return constants.SpiderVersionBucket + spiderName
}
