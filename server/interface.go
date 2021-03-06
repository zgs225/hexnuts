package server

import (
	"io"
)

type Configer interface {
	// 设置配置
	// 参数是key, value
	Set(string, string) error

	// 获取配置
	Get(string) (string, error)

	// 删除配置
	Del(string) error

	// 更新配置
	Update(string, string) error

	// 获取所有配置
	All() [][2]string
}

type PersistentConfiger interface {
	Configer

	Dumps(w io.Writer) error

	Loads(r io.Reader) (PersistentConfiger, error)

	Dirty() bool
}
