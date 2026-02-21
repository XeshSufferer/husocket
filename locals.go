package husocket

import "sync"

type Locals struct {
	m       sync.RWMutex
	content map[string]interface{}
}

func NewContext() *Locals {
	return &Locals{
		content: make(map[string]interface{}),
		m:       sync.RWMutex{},
	}
}

func (c *Locals) Set(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()
	c.content[key] = value
}

func (c *Locals) Get(key string) interface{} {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.content[key]
}

func (c *Locals) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.content, key)
}

func (c *Locals) Exists(key string) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	_, ok := c.content[key]
	return ok
}
