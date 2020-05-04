package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	"github.com/google/uuid"
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
			if err := json.Unmarshal(node.Value, &spider); err != nil {
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
		return 0, err
	}

	return total, nil
}

// 根据名称查询爬虫
func (t *Tx) SelectSpiderWhereName(name string) (spider *models.Spider, err error) {
	if node, err := t.tx.ZGetByKey(constants.SpiderListBucket, []byte(name)); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err := json.Unmarshal(node.Value, &spider); err != nil {
		return nil, err
	}
	return spider, nil
}

// 插入新爬虫
func (t *Tx) InsertSpider(spider *models.Spider) (err error) {
	spider.Id = uuid.New().String()
	spider.CreateTs = time.Now()
	spider.UpdateTs = time.Now()

	score := utils.ConvertTimestamp(spider.UpdateTs)
	value, _ := json.Marshal(&spider)
	if err = t.tx.ZAdd(constants.SpiderListBucket, []byte(spider.Name), score, value); err != nil {
		return err
	}
	return nil
}

// 更新爬虫
func (t *Tx) UpdateSpider(spider *models.Spider) (err error) {
	spider.UpdateTs = time.Now()
	score := utils.ConvertTimestamp(spider.UpdateTs)
	value, _ := json.Marshal(&spider)
	if err = t.tx.ZAdd(constants.SpiderListBucket, []byte(spider.Name), score, value); err != nil {
		return err
	}
	return nil
}

// 通过名称删除爬虫
func (t *Tx) DeleteSpiderFromName(name string) (err error) {
	if err := t.tx.ZRem(constants.SpiderListBucket, name); err != nil {
		return err
	}
	return nil
}

// 查询所有爬虫版本
func (t *Tx) SelectAllSpiderVersionsWhereSpiderName(spiderName string) (versions []*models.SpiderVersion, err error) {
	if nodes, err := t.tx.ZMembers(joinVersionBucket(spiderName)); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var version *models.SpiderVersion
			if err := json.Unmarshal(node.Value, &version); err != nil {
				return nil, err
			}
			versions = append(versions, version)
		}
	}

	return versions, nil
}

// 根据 ID 查询爬虫版本
func (t *Tx) SelectSpiderVersionWhereSpiderNameAndId(spiderName string, id string) (version *models.SpiderVersion, err error) {
	if node, err := t.tx.ZGetByKey(joinVersionBucket(spiderName), []byte(id)); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err := json.Unmarshal(node.Value, &version); err != nil {
		return nil, err
	}
	return version, nil
}

// 查询爬虫下最新的爬虫版本
func (t *Tx) SelectLatestSpiderVersionWhereSpiderName(spiderName string) (version *models.SpiderVersion, err error) {
	if node, err := t.tx.ZPeekMax(joinVersionBucket(spiderName)); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if node != nil {
		if err := json.Unmarshal(node.Value, &version); err != nil {
			return nil, err
		}
	}
	return version, nil
}

// 插入新爬虫版本
func (t *Tx) InsertSpiderVersion(version *models.SpiderVersion) (err error) {
	if version.Id == "" {
		version.Id = uuid.New().String()
	}
	version.CreateTs = time.Now()
	version.UpdateTs = time.Now()

	score := utils.ConvertTimestamp(version.UpdateTs)
	value, _ := json.Marshal(&version)
	if err = t.tx.ZAdd(joinVersionBucket(version.SpiderName), []byte(version.Id), score, value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除爬虫版本
func (t *Tx) DeleteSpiderVersionFromSpiderNameAndId(spiderName string, id string) (err error) {
	if err := t.tx.ZRem(joinVersionBucket(spiderName), id); err != nil {
		return err
	}
	return nil
}

// 根据爬虫名称删除所有爬虫版本
func (t *Tx) DeleteAllSpiderVersionsFromSpiderName(spiderName string) (err error) {
	verBucket := joinVersionBucket(spiderName)
	if nodes, err := t.tx.ZMembers(verBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	} else {
		for _, node := range nodes {
			if err := t.tx.ZRem(verBucket, node.Key()); err != nil {
				return err
			}
		}
	}
	return nil
}

func joinVersionBucket(spiderName string) string {
	return constants.SpiderVersionBucket + spiderName
}
