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

func cacheGet(key string, mode int) (DataSet, MapRows, error) {
	reply, err := cacheIns.GetString(key)
	var dataset DataSet
	if err != nil {
		return dataset, nil, err
	}

	if mode == 1 { // DataRows
		rows := MapRows{}
		if err := json.Unmarshal([]byte(reply), &rows); err != nil {
			return dataset, nil, err
		}
		return dataset, rows, nil
	}

	json.Unmarshal([]byte(reply), &dataset)
	return dataset, nil, err
}

// func cacheGetX(key string, cls IVO) (VODataSet, error) {
// 	reply, err := cacheIns.GetString(key)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	sm := NewStructConvert(cls)
//
// 	styp := reflect.SliceOf(sm.GetType())
// 	val := reflect.New(styp)
// 	if err := json.Unmarshal([]byte(reply), val.Interface()); err != nil {
// 		return nil, err
// 	}
// 	val = val.Elem()
//
// 	l := val.Len()
// 	d := make([]IVO, l)
// 	for i := 0; i < l; i++ {
// 		d[i] = (val.Index(i).Interface()).(IVO)
// 	}
// 	return VODataSet(d), nil
// }
