package monitor

import (
	"net"
)

type Server interface {
	Register(string, net.Conn) error

	Live(string) error

	Deregister(string)

	Notify(string, string, string) error
}
