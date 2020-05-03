package services

import (
	"crawlab-lite/constants"
	"crawlab-lite/database"
	"crawlab-lite/forms"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/xujiajun/nutsdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func QuerySpiderList(pageNum int, pageSize int) (total int, spiders []*models.Spider, err error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize

	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
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
	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
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

func AddSpider(form forms.SpiderForm) (spider *models.Spider, err error) {
	spiderName := form.Name

	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider != nil {
		return nil, errors.New("spider already exists")
	}

	spider = &models.Spider{
		Name:     spiderName,
		CreateTs: time.Now(),
	}

	// 存储爬虫信息
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		score := utils.ConvertTimestamp(spider.CreateTs)
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

	// 删除爬虫相关数据
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		// 删除爬虫
		if err := tx.ZRem(constants.SpiderListBucket, spiderName); err != nil {
			return err
		}
		// 删除爬虫版本
		verBucket := joinVersionBucket(spiderName)
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
	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
		// 查询爬虫下的所有版本信息
		if nodes, err := tx.ZMembers(joinVersionBucket(spiderName)); err != nil {
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

func QuerySpiderVersionById(spiderName string, versionId string) (version *models.SpiderVersion, err error) {
	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(joinVersionBucket(spiderName), []byte(versionId)); err != nil {
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

func AddSpiderVersion(spiderName string, form forms.SpiderUploadForm) (version *models.SpiderVersion, err error) {
	// 检查爬虫是否已存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider == nil {
		return nil, errors.New("spider not found")
	}

	dirName := uuid.New().String()
	tmpPath := filepath.Join(viper.GetString("other.tmppath"), spiderName)
	zipPath := filepath.Join(tmpPath, dirName+".zip")

	formFile, err := form.File.Open()
	if err != nil {
		return nil, err
	}
	defer formFile.Close()

	// 检查创建临时目录
	if utils.PathExist(tmpPath) == false {
		if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 保存压缩文件到临时目录
	if err := utils.SaveFile(formFile, zipPath); err != nil {
		return nil, err
	}

	// 打开压缩文件
	zipFile, err := os.OpenFile(zipPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer func() {
		// 关闭压缩文件
		_ = zipFile.Close()
		// 删除压缩文件
		_ = os.Remove(zipPath)
	}()

	version = &models.SpiderVersion{
		Id:       utils.GetFileMD5(zipFile),
		Path:     filepath.Join(spiderName, dirName),
		CreateTs: time.Now(),
	}

	// 通过文件 MD5 作为 ID，可以根据 ID 判断是否为重复的版本
	if version, err := QuerySpiderVersionById(spiderName, version.Id); err != nil {
		return nil, err
	} else if version != nil {
		return nil, errors.New("version already exists")
	}

	// 解压文件
	unzipPath := filepath.Join(viper.GetString("spider.path"), version.Path)
	if err := utils.Unzip(zipFile, unzipPath); err != nil {
		_ = os.RemoveAll(unzipPath)
		return nil, err
	}

	// 修改权限
	if err := os.Chmod(unzipPath, os.ModePerm); err != nil {
		return nil, err
	}

	fileList, err := ioutil.ReadDir(unzipPath)
	if err == nil && len(fileList) == 1 && fileList[0].IsDir() {
		version.Path = filepath.Join(version.Path, fileList[0].Name())
	}

	// 存储版本信息
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		spiderBucket := joinVersionBucket(spiderName)
		score := utils.ConvertTimestamp(version.CreateTs)
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
	version, err := QuerySpiderVersionById(spiderName, versionId)
	if err != nil {
		return nil, err
	} else if version == nil {
		return nil, errors.New("version not found")
	}

	// 删除爬虫与版本数据
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		if err := tx.ZRem(joinVersionBucket(spiderName), versionId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// 删除版本文件
	dirPath := viper.GetString("spider.path")
	if err := os.RemoveAll(filepath.Join(dirPath, version.Path)); err != nil {
		return nil, err
	}

	return nil, nil
}

func joinVersionBucket(spiderName string) string {
	return constants.SpiderVersionBucket + spiderName
}
