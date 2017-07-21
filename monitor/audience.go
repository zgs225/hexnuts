package monitor

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
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
	Logger  *logrus.Entry
}

func (au *Audience) Notify(e *Event) error {
	au.Logger.Debugf("Event[%d] key=%s value=%s", e.T, e.K, e.V)
	w := bufio.NewWriter(au.Conn)
	_, err := w.WriteString(e.String())
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (au *Audience) String() string {
	return fmt.Sprintf("%s[%s]", au.Name, au.Conn.RemoteAddr().String())
}
