package server

import (
	"errors"
	"strings"
)

var (
	ErrExists      = errors.New("配置已经存在")
	ErrNotExists   = errors.New("配置不存在")
	ErrDelNotEmpty = errors.New("配置中包含其他配置，不能删除")
)

type Configurer struct {
	Items map[string]interface{}
}

func (c *Configurer) Set(k, v string) error {
	return c.set(c.Items, k, v)
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

func (c *Configurer) Get(k string) (string, error) {
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
	return c.del(c.Items, k)
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
