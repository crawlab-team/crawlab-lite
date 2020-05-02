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

func QueryProjectList(pageNum int, pageSize int) (total int, projects []*models.Project, err error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize

	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		// 查询区间内的所有项目
		if nodes, err := tx.ZRangeByRank(constants.ProjectListBucket, start, end); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var project *models.Project
				if err := json.Unmarshal(node.Value, &project); err != nil {
					return err
				}
				projects = append(projects, project)
			}
		}

		// 查询数据总数目
		if total, err = tx.ZCard(constants.ProjectListBucket); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, projects, nil
}

func QueryProjectByName(name string) (project *models.Project, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(constants.ProjectListBucket, []byte(name)); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else if err := json.Unmarshal(node.Value, &project); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return project, nil
}

func SaveProject(form forms.ProjectForm) (project *models.Project, err error) {
	projectName := form.Name

	// 检查项目是否已存在
	if project, err := QueryProjectByName(projectName); err != nil {
		return nil, err
	} else if project != nil {
		return nil, errors.New("project already exists")
	}

	project = &models.Project{
		Name:     projectName,
		CreateTs: utils.NowTimestamp(),
	}

	db := database.GetKvDB()
	defer db.Close()

	// 存储项目信息
	if err := db.Update(func(tx *nutsdb.Tx) error {
		score := float64(project.CreateTs) / math.Pow10(0)
		value, _ := json.Marshal(&project)
		return tx.ZAdd(constants.ProjectListBucket, []byte(projectName), score, value)
	}); err != nil {
		return nil, err
	}

	return project, nil
}

func RemoveProject(projectName string) (res interface{}, err error) {
	// 检查项目是否已存在
	if project, err := QueryProjectByName(projectName); err != nil {
		return nil, err
	} else if project == nil {
		return nil, errors.New("project not found")
	}

	// 删除版本文件
	projectDir := viper.GetString("spider.path")
	dirs, _ := ioutil.ReadDir(projectDir)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == projectName {
			if err := os.RemoveAll(filepath.Join(projectDir, projectName)); err != nil {
				return nil, err
			}
		}
	}

	db := database.GetKvDB()
	defer db.Close()

	// 删除项目与版本数据
	if err := db.Update(func(tx *nutsdb.Tx) error {
		if err := tx.ZRem(constants.ProjectListBucket, projectName); err != nil {
			return err
		}
		verBucket := getVersionBucket(projectName)
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
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func QueryProjectVersionList(projectName string) (versions []*models.ProjectVersion, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		// 查询项目下的所有版本信息
		if nodes, err := tx.ZMembers(getVersionBucket(projectName)); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var version *models.ProjectVersion
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

func GetProjectVersionById(projectName string, versionId string) (version *models.ProjectVersion, err error) {
	db := database.GetKvDB()
	defer db.Close()

	if err := db.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(getVersionBucket(projectName), []byte(versionId)); err != nil {
			if err == nutsdb.ErrBucket {
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

func SaveProjectVersion(projectName string, form forms.ProjectUploadForm) (version *models.ProjectVersion, err error) {
	// 检查项目是否已存在
	if project, err := QueryProjectByName(projectName); err != nil {
		return nil, err
	} else if project == nil {
		return nil, errors.New("project not found")
	}

	// 生成存储版本文件的路径
	fileName := strconv.FormatInt(utils.NowTimestamp(), 10) + ".zip"
	dir := filepath.Join(viper.GetString("spider.path"), projectName)
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

	version = &models.ProjectVersion{
		Id:       utils.GetFileMD5(file),
		Path:     filepath.Join(projectName, fileName),
		UploadTs: utils.NowTimestamp(),
	}

	db := database.GetKvDB()
	defer db.Close()

	// 根据 MD5 判断是否为重复的版本
	if version, err := GetProjectVersionById(projectName, version.Id); err != nil {
		return nil, err
	} else if version != nil {
		return nil, errors.New("version already exists")
	}

	// 存储版本信息
	if err := db.Update(func(tx *nutsdb.Tx) error {
		projectBucket := getVersionBucket(projectName)
		score := float64(version.UploadTs) / math.Pow10(0)
		value, _ := json.Marshal(&version)
		return tx.ZAdd(projectBucket, []byte(version.Id), score, value)
	}); err != nil {
		return nil, err
	}

	return version, nil
}

func RemoveProjectVersion(projectName string, versionId string) (res interface{}, err error) {
	// 检查项目是否已存在
	if project, err := QueryProjectByName(projectName); err != nil {
		return nil, err
	} else if project == nil {
		return nil, errors.New("project not found")
	}

	// 查询版本信息
	version, err := GetProjectVersionById(projectName, versionId)
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

	// 删除项目与版本数据
	if err := db.Update(func(tx *nutsdb.Tx) error {
		if err := tx.ZRem(getVersionBucket(projectName), versionId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func getVersionBucket(projectName string) string {
	return constants.ProjectVersionBucket + projectName
}
