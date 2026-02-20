package husocket

import "sync"

type Context struct {
	m       sync.RWMutex
	content map[string]interface{}
}

func NewContext() *Context {
	return &Context{
		content: make(map[string]interface{}),
		m:       sync.RWMutex{},
	}
}

func (c *Context) Set(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()
	c.content[key] = value
}

func (c *Context) Get(key string) interface{} {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.content[key]
}

func (c *Context) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.content, key)
}

func (c *Context) Exists(key string) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	_, ok := c.content[key]
	return ok
}
