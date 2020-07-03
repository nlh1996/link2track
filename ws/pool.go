package ws

import (
	"sync"
)

// ConnPool .
type ConnPool struct {
	Pool map[string]*Connection
	sync.RWMutex
}

var instance *ConnPool

// GetConnPool .
func GetConnPool() *ConnPool {
	if instance == nil {
		instance = &ConnPool{}
		instance.Pool = make(map[string]*Connection, 10)
	}
	return instance
}

// GetConnByID .
func (p *ConnPool) GetConnByID(id string) *Connection {
	p.Lock()
	defer p.Unlock()
	if v, ok := p.Pool[id]; ok {
		return v
	}
	return nil
}

// Set .
func (p *ConnPool) Set(c *Connection) {
	p.Lock()
	defer p.Unlock()
	p.Pool[c.ID] = c
}

// DelByID .
func (p *ConnPool) DelByID(id string) {
	p.Lock()
	defer p.Unlock()
	delete(p.Pool, id)
}
