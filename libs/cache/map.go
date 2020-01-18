package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/kere/gno/libs/log"
)

// ICachedMap interface
type ICachedMap interface {
	Build(args ...interface{}) (interface{}, error)
	Init(c ICachedMap, expires int)
	Validate(v interface{}) bool
	SetExpires(v int)
	GetExpires() int
	Get(args ...interface{}) interface{}
	Set(obj interface{}, args ...interface{})
	Release(args ...interface{})
	IsCached(args ...interface{}) bool
	ClearAll()
}

// ExpiresVal class
type ExpiresVal struct {
	Value     interface{}
	Expires   time.Duration
	ExpiresAt time.Time
}

// Map class
type Map struct {
	target  ICachedMap
	Lock    sync.RWMutex
	Data    map[string]ExpiresVal
	expires time.Duration
}

// Init class
func (m *Map) Init(t ICachedMap, expires int) {
	m.target = t
	m.expires = time.Duration(expires) * time.Second
	m.Data = make(map[string]ExpiresVal, 0)
	// m.Data = sync.Map{}
}

// SetExpires func
func (m *Map) SetExpires(second int) {
	m.expires = time.Duration(second) * time.Second
}

// GetExpires func
func (m *Map) GetExpires() int {
	return int(m.expires.Seconds())
}

// buildKey func
func (m *Map) buildKey(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Validate 检查缓存值是否正确，如果正确才保存
func (m *Map) Validate(v interface{}) bool {
	return true
}

// Set func
func (m *Map) Set(obj interface{}, args ...interface{}) {
	if obj == nil || !m.target.Validate(obj) {
		return
	}
	key := m.buildKey(args...)
	ex := time.Now().Add(m.expires)
	m.Lock.Lock()
	m.Data[key] = ExpiresVal{Value: obj, Expires: m.expires, ExpiresAt: ex}
	m.Lock.Unlock()
}

// Get func
func (m *Map) Get(args ...interface{}) interface{} {
	// 读取数据
	key := m.buildKey(args...)
	m.Lock.RLock()
	// v, isok := m.Data.Load(key)
	v, isok := m.Data[key]
	if isok && v.isNotExpired() {
		m.Lock.RUnlock()
		return v.Value
	}
	m.Lock.RUnlock()

	// 同步读取数据
	m.Lock.Lock()
	// v, isok := m.Data.Load(key)
	v, isok = m.Data[key]
	if isok && v.isNotExpired() {
		m.Lock.Unlock()
		return v.Value
	}
	// debug.PrintStack()

	obj, err := m.target.Build(args...)
	if err != nil {
		log.App.Error(err)
		m.Lock.Unlock()
		return nil
	}

	log.App.Debug("[map cache] build", key, "expires:", m.GetExpires())

	if obj == nil || !m.target.Validate(obj) {
		m.Lock.Unlock()
		return nil
	}

	ex := time.Now().Add(m.expires)
	v = ExpiresVal{Value: obj, Expires: m.expires, ExpiresAt: ex}
	// m.Data.Store(key, val)
	m.Data[key] = v
	m.Lock.Unlock()
	return obj
}

// ClearAll release all
func (m *Map) ClearAll() {
	m.Lock.Lock()
	m.Data = make(map[string]ExpiresVal, 0)
	// m.Data = sync.Map{}
	m.Lock.Unlock()
}

// Release 释放缓存
func (m *Map) Release(args ...interface{}) {
	m.Lock.Lock()
	key := m.buildKey(args...)
	// m.Data.Delete(key)
	delete(m.Data, key)
	m.Lock.Unlock()
}

// IsCached bool
func (m *Map) IsCached(args ...interface{}) bool {
	m.Lock.Lock()
	key := m.buildKey(args...)
	_, isok := m.Data[key]
	m.Lock.Unlock()
	return isok
}

// Print 打印缓存的buildKey
func (m *Map) Print() {
	count := 0
	// m.Data.Range(func(k interface{}, value interface{}) bool {
	// 	count++
	// 	fmt.Println(count, ":", k, ":", value)
	// 	return true
	// })
	for k, v := range m.Data {
		count++
		fmt.Println(count, ":", k, ":", v)
	}
}

// isExpired value is expired
func (e ExpiresVal) isExpired() bool {
	if e.Expires == 0 {
		return false
	}
	return e.ExpiresAt.Before(time.Now())
}

// isNotExpired value is expired
func (e ExpiresVal) isNotExpired() bool {
	if e.Expires == 0 {
		return true
	}
	return e.ExpiresAt.After(time.Now())
}
