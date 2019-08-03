package db

import (
	"encoding/json"

	"github.com/kere/gno/libs/cache"
)

var cacheIns cache.ICache

// SetCache f
func SetCache(c cache.ICache) {
	cacheIns = c
}

func cacheDel(key string) error {
	return cacheIns.Delete(key)
}

func cacheSet(key string, value interface{}, expire int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = cacheIns.Set(key, string(v), expire)
	return err
}

func cacheGetDataSet(key string) (DataSet, error) {
	reply, err := cacheIns.GetString(key)
	var dataset DataSet
	if err != nil {
		return dataset, err
	}

	json.Unmarshal([]byte(reply), &dataset)
	return dataset, err
}

func cacheGetRows(key string) (MapRows, error) {
	reply, err := cacheIns.GetString(key)
	if err != nil {
		return nil, err
	}

	rows := MapRows{}
	if err := json.Unmarshal([]byte(reply), &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
