package services

import (
	"crawlab-lite/dao"
	"crawlab-lite/forms"
	"crawlab-lite/models"
	"crawlab-lite/results"
	"crawlab-lite/utils"
	"errors"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func QuerySpiderPage(page forms.PageForm) (total int, resultList []*results.Spider, err error) {
	start, end := page.Range()

	if err := dao.ReadTx(func(tx dao.Tx) error {
		spiders, err := tx.SelectAllSpidersLimit(start, end)
		if err != nil {
			return err
		}
		if total, err = tx.CountSpiders(); err != nil {
			return err
		}
		cache := map[uuid.UUID]*models.Task{}
		for _, spider := range spiders {
			var result results.Spider
			if err := copier.Copy(&result, spider); err != nil {
				return err
			}
			var task *models.Task
			task, exists := cache[spider.Id]
			if !exists {
				if task, err = tx.SelectLatestTaskWhereSpiderId(spider.Id); err != nil {
					return err
				}
			}
			if task != nil {
				result.LastRunTs = task.StartTs
				result.LastStatus = task.Status
			}
			resultList = append(resultList, &result)
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, resultList, nil
}

func QuerySpider(id uuid.UUID) (result *results.Spider, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		spider, err := tx.SelectSpider(id)
		if err != nil {
			return err
		}
		if err := copier.Copy(&result, spider); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func AddSpider(form forms.SpiderForm) (result *results.Spider, err error) {
	spiderName := form.Name

	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查爬虫是否已存在
		spider, err := tx.SelectSpiderWhereName(spiderName)
		if err != nil {
			return err
		}
		if spider != nil {
			return errors.New("spider name already exists")
		}

		// 存储爬虫信息
		spider = &models.Spider{
			Name: spiderName,
		}
		if err = tx.InsertSpider(spider); err != nil {
			return err
		}

		if err := copier.Copy(&result, spider); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveSpider(id uuid.UUID) (res interface{}, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查爬虫是否存在
		if spider, err := tx.SelectSpider(id); err != nil {
			return err
		} else if spider == nil {
			return errors.New("spider not found")
		}

		// 删除爬虫
		if err = tx.DeleteSpider(id); err != nil {
			return err
		}

		// 删除爬虫的所有版本
		if err = tx.DeleteAllSpiderVersions(id); err != nil {
			return err
		}

		// 删除爬虫的所有任务
		if err = tx.DeleteAllTasksWhereSpiderId(id); err != nil {
			return err
		}

		// 删除版本文件
		spiderDir := viper.GetString("spider.path")
		dirs, _ := ioutil.ReadDir(spiderDir)
		for _, dir := range dirs {
			if dir.IsDir() && dir.Name() == id.String() {
				if err := os.RemoveAll(filepath.Join(spiderDir, id.String())); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func QuerySpiderVersionList(spiderId uuid.UUID) (resultList []*results.SpiderVersion, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		versions, err := tx.SelectAllSpiderVersions(spiderId)
		if err != nil {
			return err
		}
		for _, version := range versions {
			var result results.SpiderVersion
			if err := copier.Copy(&result, version); err != nil {
				return err
			}
			resultList = append(resultList, &result)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return resultList, nil
}

func QuerySpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (result *results.SpiderVersion, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		version, err := tx.SelectSpiderVersion(spiderId, versionId)
		if err != nil {
			return err
		}
		if err := copier.Copy(&result, version); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func QueryLatestSpiderVersion(spiderId uuid.UUID) (result *results.SpiderVersion, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		version, err := tx.SelectLatestSpiderVersion(spiderId)
		if err != nil {
			return err
		}
		if err := copier.Copy(&result, version); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func AddSpiderVersion(spiderId uuid.UUID, form forms.SpiderUploadForm) (result *results.SpiderVersion, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查爬虫是否存在
		if spider, err := tx.SelectSpider(spiderId); err != nil {
			return err
		} else if spider == nil {
			return errors.New("spider not found")
		}

		uid := uuid.NewV4()
		tmpPath := filepath.Join(viper.GetString("other.tmppath"), spiderId.String())
		zipPath := filepath.Join(tmpPath, uid.String()+".zip")

		formFile, err := form.File.Open()
		if err != nil {
			return err
		}
		defer formFile.Close()

		// 检查创建临时目录
		if utils.PathExist(tmpPath) == false {
			if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
				return err
			}
		}

		// 保存压缩文件到临时目录
		if err := utils.SaveFile(formFile, zipPath); err != nil {
			return err
		}

		// 打开压缩文件
		zipFile, err := os.OpenFile(zipPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		defer func() {
			// 关闭压缩文件
			_ = zipFile.Close()
			// 删除解压临时文件夹
			_ = os.RemoveAll(tmpPath)
		}()

		version := &models.SpiderVersion{
			Id:       uid,
			FileHash: utils.GetFileMD5(zipFile),
			SpiderId: spiderId,
			Path:     filepath.Join(spiderId.String(), uid.String()),
		}

		// 通过文件 MD5 判断是否为重复的版本
		if _version, err := tx.SelectSpiderVersionWhereFileHash(spiderId, version.FileHash); err != nil {
			return err
		} else if _version != nil {
			return errors.New("spider version already exists")
		}

		// 解压文件
		unzipPath := filepath.Join(viper.GetString("spider.path"), version.Path)
		if err := utils.Unzip(zipFile, unzipPath); err != nil {
			_ = os.RemoveAll(unzipPath)
			return err
		}

		// 修改权限
		if err := os.Chmod(unzipPath, os.ModePerm); err != nil {
			return err
		}

		fileList, err := ioutil.ReadDir(unzipPath)
		if err == nil && len(fileList) == 1 && fileList[0].IsDir() {
			version.Path = filepath.Join(version.Path, fileList[0].Name())
		}

		// 存储版本信息
		if err := tx.InsertSpiderVersion(version); err != nil {
			return err
		}

		if err := copier.Copy(&result, version); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveSpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (res interface{}, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查爬虫是否存在
		if spider, err := tx.SelectSpider(spiderId); err != nil {
			return err
		} else if spider == nil {
			return errors.New("spider not found")
		}

		// 查询版本信息
		version, err := tx.SelectSpiderVersion(spiderId, versionId)
		if err != nil {
			return err
		} else if version == nil {
			return errors.New("spider version not found")
		}

		// 删除爬虫版本信息
		if err = tx.DeleteSpiderVersion(spiderId, versionId); err != nil {
			return err
		}

		// 查询版本信息
		version2, err := tx.SelectSpiderVersion(spiderId, versionId)
		if err != nil {
			return err
		} else if version2 == nil {
			return errors.New("spider version not found")
		}

		// 删除版本文件
		dirPath := viper.GetString("spider.path")
		if err := os.RemoveAll(filepath.Join(dirPath, version.Path)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
