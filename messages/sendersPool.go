package messages

import (
	"sync"
)

var once sync.Once
var SendersPool *Pool

type Pool struct {
	mx sync.RWMutex
	m  map[string]bool
}

func NewPool() {
	SendersPool = &Pool{
		m: make(map[string]bool),
	}
}

func (p *Pool) Load(key string) bool {
	p.mx.RLock()
	defer p.mx.RUnlock()
	val, _ := p.m[key]
	return val
}

func (p *Pool) Store(key string, value bool) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.m[key] = value
}
