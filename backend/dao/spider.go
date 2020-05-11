package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/xujiajun/nutsdb"
	"time"
)

// 查询区间内的所有爬虫
func (t *Tx) SelectAllSpidersLimit(start int, end int) (spiders []*models.Spider, err error) {
	if nodes, err := t.tx.ZRangeByRank(constants.SpiderListBucket, start, end); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var spider *models.Spider
			if err = json.Unmarshal(node.Value, &spider); err != nil {
				return nil, err
			}
			spiders = append(spiders, spider)
		}
	}

	return spiders, nil
}

// 所有爬虫的总数目
func (t *Tx) CountSpiders() (total int, err error) {
	if total, err = t.tx.ZCard(constants.SpiderListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}

// 根据 ID 查询爬虫
func (t *Tx) SelectSpider(id uuid.UUID) (spider *models.Spider, err error) {
	if node, err := t.tx.ZGetByKey(constants.SpiderListBucket, []byte(id.String())); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err = json.Unmarshal(node.Value, &spider); err != nil {
		return nil, err
	}
	return spider, nil
}

// 根据名称查询爬虫
func (t *Tx) SelectSpiderWhereName(spiderName string) (spider *models.Spider, err error) {
	if nodes, err := t.tx.ZMembers(constants.SpiderListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var spider *models.Spider
			if err = json.Unmarshal(node.Value, &spider); err != nil {
				return nil, err
			}
			if spider.Name == spiderName {
				return spider, nil
			}
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

	score := utils.ConvertTimestamp(spider.UpdateTs)
	value, _ := json.Marshal(&spider)
	if err = t.tx.ZAdd(constants.SpiderListBucket, []byte(spider.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 更新爬虫
func (t *Tx) UpdateSpider(spider *models.Spider) (err error) {
	spider.UpdateTs = time.Now()
	score := utils.ConvertTimestamp(spider.UpdateTs)
	value, _ := json.Marshal(&spider)
	if err = t.tx.ZAdd(constants.SpiderListBucket, []byte(spider.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除爬虫
func (t *Tx) DeleteSpider(id uuid.UUID) (err error) {
	if err = t.tx.ZRem(constants.SpiderListBucket, id.String()); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	}
	return nil
}

// 查询所有爬虫版本
func (t *Tx) SelectAllSpiderVersions(spiderId uuid.UUID) (versions []*models.SpiderVersion, err error) {
	if nodes, err := t.tx.ZMembers(joinVersionBucket(spiderId)); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var version *models.SpiderVersion
			if err = json.Unmarshal(node.Value, &version); err != nil {
				return nil, err
			}
			versions = append(versions, version)
		}
	}

	return versions, nil
}

// 根据 ID 查询爬虫版本
func (t *Tx) SelectSpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (version *models.SpiderVersion, err error) {
	if node, err := t.tx.ZGetByKey(joinVersionBucket(spiderId), []byte(versionId.String())); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err = json.Unmarshal(node.Value, &version); err != nil {
		return nil, err
	}
	return version, nil
}

// 根据 MD5 查询爬虫版本
func (t *Tx) SelectSpiderVersionWhereFileHash(spiderId uuid.UUID, fileHash string) (version *models.SpiderVersion, err error) {
	if nodes, err := t.tx.ZMembers(joinVersionBucket(spiderId)); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var version *models.SpiderVersion
			if err = json.Unmarshal(node.Value, &version); err != nil {
				return nil, err
			}
			if version.FileHash == fileHash {
				return version, nil
			}
		}
	}
	return nil, nil
}

// 查询爬虫下最新的爬虫版本
func (t *Tx) SelectLatestSpiderVersion(spiderId uuid.UUID) (version *models.SpiderVersion, err error) {
	if node, err := t.tx.ZPeekMax(joinVersionBucket(spiderId)); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else if node != nil {
		if err = json.Unmarshal(node.Value, &version); err != nil {
			return nil, err
		}
	}
	return version, nil
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

	score := utils.ConvertTimestamp(version.UpdateTs)
	value, _ := json.Marshal(&version)
	if err = t.tx.ZAdd(joinVersionBucket(version.SpiderId), []byte(version.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除爬虫版本
func (t *Tx) DeleteSpiderVersion(spiderId uuid.UUID, versionId uuid.UUID) (err error) {
	if err := t.tx.ZRem(joinVersionBucket(spiderId), versionId.String()); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	}
	return nil
}

// 根据 ID 删除所有爬虫版本
func (t *Tx) DeleteAllSpiderVersions(spiderId uuid.UUID) (err error) {
	verBucket := joinVersionBucket(spiderId)
	if nodes, err := t.tx.ZMembers(verBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	} else {
		for _, node := range nodes {
			print(node.Key())

			if err = t.tx.ZRem(verBucket, node.Key()); err != nil {
				return err
			}
		}
	}
	return nil
}

func joinVersionBucket(spiderId uuid.UUID) string {
	return constants.SpiderVersionBucket + spiderId.String()
}
