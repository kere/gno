package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/kere/gno/libs/log"
)

// ICachedMap interface
type ICachedMap interface {
	Build(args ...interface{}) (interface{}, int, error)
	Init(ICachedMap)
	CheckValue(v interface{}) bool
}

// ExpiresVal class
type ExpiresVal struct {
	Value     interface{}
	Expires   int
	ExpiresAt time.Time
}

// Map class
type Map struct {
	target ICachedMap
	Lock   *sync.RWMutex
	// Data        map[string]interface{}
	Data map[string]*ExpiresVal
}

// Init class
func (m *Map) Init(t ICachedMap) {
	m.Lock = new(sync.RWMutex)
	m.Data = make(map[string]*ExpiresVal, 0)
	m.target = t
}

// Key func
func (m *Map) Key(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// CheckValue 检查缓存值是否正确，如果正确才保存
func (m *Map) CheckValue(v interface{}) bool {
	return true
}

// ClearAll release all
func (m *Map) ClearAll() {
	m.Lock.Lock()
	m.Data = make(map[string]*ExpiresVal, 0)
	m.Lock.Unlock()
}

// Get func
func (m *Map) Get(args ...interface{}) interface{} {
	key := m.Key(args...)
	m.Lock.RLock()
	// key = market+code:dtype
	if v, isok := m.Data[key]; isok && isNotExpired(v) {
		m.Lock.RUnlock()
		return v.Value
	}
	m.Lock.RUnlock()

	// not found ---------------------
	m.Lock.Lock()
	if v, isok := m.Data[key]; isok && isNotExpired(v) {
		m.Lock.RUnlock()
		return v.Value
	}

	obj, expires, err := m.target.Build(args...)
	if err != nil {
		log.App.Notice(err)
		m.Lock.Unlock()
		return nil
	}

	if obj == nil || !m.target.CheckValue(obj) {
		m.Lock.Unlock()
		return nil
	}

	ex := time.Now().Add(time.Duration(expires) * time.Second)
	v := &ExpiresVal{Value: obj, Expires: expires, ExpiresAt: ex}

	m.Data[key] = v

	m.Lock.Unlock()
	return v.Value
}

// Release 释放缓存
func (m *Map) Release(args ...interface{}) {
	key := m.Key(args...)
	m.Lock.RLock()
	// key = market+code:dtype
	if _, isok := m.Data[key]; !isok {
		m.Lock.RUnlock()
		return
	}
	m.Lock.RUnlock()
	// ----------------------
	m.Lock.Lock()
	if _, isok := m.Data[key]; !isok {
		m.Lock.Unlock()
		return
	}

	delete(m.Data, key)
	m.Lock.Unlock()
}

// IsCached bool
func (m *Map) IsCached(args ...interface{}) bool {
	key := m.Key(args...)
	m.Lock.RLock()
	// key = market+code:dtype
	_, isok := m.Data[key]
	m.Lock.RUnlock()
	return isok
}

// Print 打印缓存的Key
func (m *Map) Print() {
	count := 0
	for k := range m.Data {
		count++
		fmt.Println(count, ":", k)
	}
}

// isExpired value is expired
func isExpired(e *ExpiresVal) bool {
	return e != nil && e.Expires > 0 && e.ExpiresAt.Before(time.Now())
}

// isNotExpired value is expired
func isNotExpired(e *ExpiresVal) bool {
	return e != nil && e.Expires > 0 && e.ExpiresAt.After(time.Now())
}
