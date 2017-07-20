package monitor

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Audience struct {
	Ctx     context.Context
	Name    string
	Conn    net.Conn
	Alive   bool
	Last    time.Time
	Ch      chan *Event
	Timeout time.Duration
	Cancel  context.CancelFunc
}

func (au *Audience) Notify(e *Event) error {
	return nil
}

func (au *Audience) String() string {
	return fmt.Sprintf("%s[%s]", au.Name, au.Conn.RemoteAddr().String())
}
