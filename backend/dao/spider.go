package dao

import (
	"bytes"
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

// 查询所有爬虫
func (t *Tx) SelectAllSpiders() (spiders []*models.Spider, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderBucket))
	if b == nil {
		return spiders, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var spider *models.Spider
		if err = json.Unmarshal(v, &spider); err != nil {
			return nil, err
		}
		spiders = append(spiders, spider)
	}
	return spiders, nil
}

// 根据 ID 查询爬虫
func (t *Tx) SelectSpider(id uuid.UUID) (spider *models.Spider, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id.String()))
	if value == nil {
		return nil, nil
	}
	if err = json.Unmarshal(value, &spider); err != nil {
		return nil, err
	}
	return spider, nil
}

// 根据名称查询爬虫
func (t *Tx) SelectSpiderWhereName(spiderName string) (spider *models.Spider, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderBucket))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if err = json.Unmarshal(v, &spider); err != nil {
			return nil, err
		}
		if spider.Name == spiderName {
			return spider, nil
		}
	}
	return nil, nil
}

// 插入新爬虫
func (t *Tx) InsertSpider(spider *models.Spider) (err error) {
	if spider.Id == uuid.Nil {
		spider.Id = uuid.NewV4()
	}
	if spider.CreateTs.IsZero() {
		spider.CreateTs = time.Now()
	}
	if spider.UpdateTs.IsZero() {
		spider.UpdateTs = time.Now()
	}

	value, err := json.Marshal(&spider)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.SpiderBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(spider.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 更新爬虫
func (t *Tx) UpdateSpider(spider *models.Spider) (err error) {
	b := t.tx.Bucket([]byte(constants.SpiderBucket))
	if b == nil {
		return nil
	}
	spider.UpdateTs = time.Now()
	value, _ := json.Marshal(&spider)
	if err = b.Put([]byte(spider.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除爬虫
func (t *Tx) DeleteSpider(id uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.SpiderBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id.String())); err != nil {
		return err
	}
	return nil
}

// 查询所有爬虫版本
func (t *Tx) SelectAllSpiderVersions(spiderId uuid.UUID) (versions []*models.SpiderVersion, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderVersionBucket))
	if b == nil {
		return versions, nil
	}
	c := b.Cursor()
	prefix := []byte(spiderId.String())
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var version *models.SpiderVersion
		if err = json.Unmarshal(v, &version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// 根据 ID 查询爬虫版本
func (t *Tx) SelectSpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (version *models.SpiderVersion, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderVersionBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(joinVersionKey(spiderId, versionId)))
	if value == nil {
		return nil, nil
	}
	if err = json.Unmarshal(value, &version); err != nil {
		return nil, err
	}
	return version, nil
}

// 根据 MD5 查询爬虫版本
func (t *Tx) SelectSpiderVersionWhereFileHash(spiderId uuid.UUID, fileHash string) (version *models.SpiderVersion, err error) {
	b := t.tx.Bucket([]byte(constants.SpiderVersionBucket))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	prefix := []byte(spiderId.String())
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		if err = json.Unmarshal(v, &version); err != nil {
			return nil, err
		}
		if version.FileHash == fileHash {
			return version, nil
		}
	}
	return nil, nil
}

// 插入新爬虫版本
func (t *Tx) InsertSpiderVersion(version *models.SpiderVersion) (err error) {
	if version.Id == uuid.Nil {
		version.Id = uuid.NewV4()
	}
	if version.CreateTs.IsZero() {
		version.CreateTs = time.Now()
	}
	if version.UpdateTs.IsZero() {
		version.UpdateTs = time.Now()
	}

	value, err := json.Marshal(&version)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.SpiderVersionBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(joinVersionKey(version.SpiderId, version.Id)), value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除爬虫版本
func (t *Tx) DeleteSpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.SpiderVersionBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(joinVersionKey(spiderId, versionId))); err != nil {
		return err
	}
	return nil
}

// 根据 ID 删除所有爬虫版本
func (t *Tx) DeleteAllSpiderVersions(spiderId uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.SpiderVersionBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	prefix := []byte(spiderId.String())
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var version models.SpiderVersion
		if err = json.Unmarshal(v, &version); err != nil {
			return err
		}
		if err = b.Delete(k); err != nil {
			return err
		}
	}
	return nil
}

func joinVersionKey(spiderId uuid.UUID, versionId uuid.UUID) string {
	return fmt.Sprintf("%s:%s", spiderId.String(), versionId.String())
}
