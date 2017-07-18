package server

type Configer interface {
	// 设置配置
	// 参数是key, value
	Set(string, string) error

	// 获取配置
	Get(string) (string, error)

	// 删除配置
	Del(string) error
}
