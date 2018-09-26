package websock

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// GetClient get client
func GetClient(path string, id int) (Client, error) {
	m := connMap[path]
	if m == nil {
		return Client{}, errors.New("not found")
	}

	return Client{}, errors.New("not found")
}

// GetManager get connections
func GetManager(path string) *Manager {
	return connMap[path]
}

// make(map[path] map[clientID]*websocket.Conn, 0)
var connMap = make(map[string]*Manager, 0)

// Manager class
type Manager struct {
	counter int
	lock    sync.Mutex
	list    map[int]Client
}

// NewManager new
func NewManager() *Manager {
	return &Manager{list: make(map[int]Client, 0)}
}

// AllClients client
func (m *Manager) AllClients(f func(c Client)) {
	for _, c := range m.list {
		f(c)
	}
}

// ClientsCountByIP ip地址下的连接数量
func (m *Manager) ClientsCountByIP(ip string) int {
	count := 0
	for _, c := range m.list {
		if ip == c.IP {
			count++
		}
	}
	return count
}

// ClientsByIP ip地址下的连接
func (m *Manager) ClientsByIP(ip string) []Client {
	arr := make([]Client, 0)
	for _, c := range m.list {
		if ip == c.IP {
			arr = append(arr, c)
		}
	}
	return arr
}

// ClientCount client
func (m *Manager) ClientCount() int {
	return len(m.list)
}

// Close client
func (m *Manager) Close(id int) {
	c, isok := m.list[id]
	if !isok {
		return
	}
	c.Close()
	delete(m.list, id)
}

// SetClient by id
func (m *Manager) SetClient(id int, conn *websocket.Conn, req *http.Request) *Client {
	addr := req.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = req.Header.Get("X-Real-IP")
	}
	c := Client{ID: id, Conn: conn, Cookies: req.Cookies(), Form: req.Form, Header: req.Header, IP: addr}
	m.list[id] = c
	return &c
}

// GetClient by id
func (m *Manager) GetClient(id int) (Client, error) {
	c, isok := m.list[id]
	if !isok {
		return Client{}, errors.New("not found client")
	}
	return c, nil
}

// GenrateClientID genrate id
func (m *Manager) GenrateClientID() int {
	m.lock.Lock()
	m.counter++
	m.lock.Unlock()
	return m.counter
}
