package server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

func init() {
	gob.Register(new(Nodes))
	gob.Register(make(map[string]interface{}))
}

var (
	ErrExists      = errors.New("配置已经存在")
	ErrNotExists   = errors.New("配置不存在")
	ErrDelNotEmpty = errors.New("配置中包含其他配置，不能删除")
)

type Nodes map[string]interface{}

type Configurer struct {
	mu    sync.RWMutex
	dirty bool
	Items Nodes
}

func (c *Configurer) Set(k, v string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.set(c.Items, k, v); err != nil {
		return err
	}
	c.dirty = true
	return nil
}

func (c *Configurer) set(items map[string]interface{}, k, v string) error {
	i := strings.Index(k, ".")

	if i < 0 {
		if _, ok := items[k]; ok {
			return ErrExists
		}
		items[k] = v
		return nil
	}

	if items2, ok := items[k[:i]]; ok {
		if items3, ok := items2.(map[string]interface{}); !ok {
			return ErrExists
		} else {
			return c.set(items3, k[i+1:], v)
		}
	} else {
		items3 := make(map[string]interface{})
		items[k[:i]] = items3
		return c.set(items3, k[i+1:], v)
	}
}

func (c *Configurer) Update(k, v string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.update(c.Items, k, v); err != nil {
		return err
	}
	c.dirty = true
	return nil
}

func (c *Configurer) update(items map[string]interface{}, k, v string) error {
	i := strings.Index(k, ".")

	if i < 0 {
		if _, ok := items[k]; ok {
			items[k] = v
			return nil
		}
		return ErrNotExists
	}

	if items2, ok := items[k[:i]]; ok {
		if items3, ok := items2.(map[string]interface{}); !ok {
			return ErrNotExists
		} else {
			return c.update(items3, k[i+1:], v)
		}
	} else {
		return ErrNotExists
	}
}

func (c *Configurer) Get(k string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(c.Items, k)
}

func (c *Configurer) get(items map[string]interface{}, k string) (string, error) {
	i := strings.Index(k, ".")

	if i < 0 {
		if v, ok := items[k]; ok {
			if v2, ok := v.(string); ok {
				return v2, nil
			} else {
				return "", ErrNotExists
			}
		} else {
			return "", ErrNotExists
		}
	}

	if items2, ok := items[k[:i]]; ok {
		if items3, ok := items2.(map[string]interface{}); ok {
			return c.get(items3, k[i+1:])
		} else {
			return "", ErrNotExists
		}
	} else {
		return "", ErrNotExists
	}
}

func (c *Configurer) Del(k string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.del(c.Items, k); err != nil {
		return err
	}
	c.dirty = true
	return nil
}

func (c *Configurer) del(items map[string]interface{}, k string) error {
	i := strings.Index(k, ".")
	if i < 0 {
		if v, ok := items[k]; ok {
			if _, ok = v.(map[string]interface{}); !ok {
				delete(items, k)
				return nil
			}
			return ErrDelNotEmpty
		}
		return nil
	}

	if items2, ok := items[k[:i]]; ok {
		if items3, ok := items2.(map[string]interface{}); ok {
			return c.del(items3, k[i+1:])
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (c *Configurer) Dumps(w io.Writer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if err := gob.NewEncoder(w).Encode(c.Items); err != nil {
		return err
	}
	c.dirty = false
	return nil
}

func (c *Configurer) Loads(r io.Reader) (PersistentConfiger, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	items := Nodes(make(map[string]interface{}))
	if err := gob.NewDecoder(r).Decode(&items); err != nil {
		return nil, err
	}
	return &Configurer{Items: items}, nil
}

func (c *Configurer) Dirty() bool {
	return c.dirty
}

func (c *Configurer) All() [][2]string {
	rv := make([][2]string, 0)
	rv = c.all("", c.Items, rv)
	return rv
}

func (c *Configurer) all(parent string, items Nodes, rv [][2]string) [][2]string {
	for k, v := range items {
		var fk string
		if len(parent) > 0 {
			fk = fmt.Sprintf("%s.%s", parent, k)
		} else {
			fk = k
		}

		if v2, ok := v.(string); ok {
			i := [2]string{fk, v2}
			rv = append(rv, i)
		}

		if v2, ok := v.(map[string]interface{}); ok {
			rv = c.all(fk, v2, rv)
		}
	}

	return rv
}
