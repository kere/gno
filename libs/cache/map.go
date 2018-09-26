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
}

// ExpiresVal class
type ExpiresVal struct {
	Value     interface{}
	Expires   time.Duration
	ExpiresAt time.Time
}

// Map class
type Map struct {
	target ICachedMap
	Lock   sync.RWMutex
	// Data        map[string]interface{}
	Data    map[string]*ExpiresVal
	expires time.Duration
}

// Init class
func (m *Map) Init(t ICachedMap, expires int) {
	m.Data = make(map[string]*ExpiresVal, 0)
	m.target = t
	m.expires = time.Duration(expires) * time.Second
}

// SetExpires func
func (m *Map) SetExpires(expires int) {
	m.expires = time.Duration(expires) * time.Second
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

// Get func
func (m *Map) Get(args ...interface{}) interface{} {
	key := m.buildKey(args...)
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

	obj, err := m.target.Build(args...)
	if err != nil {
		log.App.Error(err)
		m.Lock.Unlock()
		return nil
	}
	log.App.Debug("cache build", key, m.GetExpires())

	if obj == nil || !m.target.Validate(obj) {
		m.Lock.Unlock()
		return nil
	}

	ex := time.Now().Add(m.expires)
	v := &ExpiresVal{Value: obj, Expires: m.expires, ExpiresAt: ex}

	m.Data[key] = v

	m.Lock.Unlock()
	return v.Value
}

// ClearAll release all
func (m *Map) ClearAll() {
	m.Lock.Lock()
	for k, v := range m.Data {
		v.Value = nil
		m.Data[k] = nil
	}
	m.Data = make(map[string]*ExpiresVal, 0)
	m.Lock.Unlock()
}

// Release 释放缓存
func (m *Map) Release(args ...interface{}) {
	key := m.buildKey(args...)
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
	key := m.buildKey(args...)
	m.Lock.RLock()
	// key = market+code:dtype
	_, isok := m.Data[key]
	m.Lock.RUnlock()
	return isok
}

// Print 打印缓存的buildKey
func (m *Map) Print() {
	count := 0
	for k := range m.Data {
		count++
		fmt.Println(count, ":", k)
	}
}

// isExpired value is expired
func isExpired(e *ExpiresVal) bool {
	if e == nil || e.Expires == 0 {
		return true
	}
	return e.ExpiresAt.Before(time.Now())
}

// isNotExpired value is expired
func isNotExpired(e *ExpiresVal) bool {
	if e == nil || e.Expires == 0 {
		return false
	}
	return e.ExpiresAt.After(time.Now())
}
