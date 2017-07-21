package monitor

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

type TCPServer struct {
	Addr      string
	TLS       bool
	Cert      string
	Key       string
	Audiences map[string]*Audience
	Ch        chan *Event
	Logger    *logrus.Logger

	mu sync.Mutex
}

func (s *TCPServer) ServeLoop() error {
	var lis net.Listener

	if s.TLS {
		cer, err := tls.LoadX509KeyPair(s.Cert, s.Key)
		if err != nil {
			return err
		}
		conf := &tls.Config{Certificates: []tls.Certificate{cer}}
		l, err := tls.Listen("tcp", s.Addr, conf)
		if err != nil {
			return err
		}
		lis = l
	} else {
		l, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}
		lis = l
	}

	go s.Heartbeat()

	go s.EventLoop()

	for {
		conn, err := lis.Accept()
		if err != nil {
			s.Logger.Error(err)
			continue
		}
		go s.handle(conn)
	}

	return nil
}

func (s *TCPServer) handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				s.Logger.Errorf("Read from conn error: %v", err)
				continue
			} else {
				return
			}
		}
		s.Logger.Debug(msg)
		if len(msg) < 3 {
			s.Logger.Errorf("Unrecognized message: %q", msg)
			continue
		}
		cmd := string(bytes.ToUpper(msg[:3]))
		switch cmd {
		case "REG":
			name := string(bytes.Trim(msg[3:], " \t\n\r\f"))
			if err := s.Register(name, conn); err != nil {
				s.Logger.Errorf("Register %s error: %v", name, err)
				s.Deregister(name)
				return
			}
		case "LIV":
			name := string(bytes.Trim(msg[3:], " \t\n\r\f"))
			if err := s.Live(name); err != nil {
				s.Logger.Errorf("Live %s error: %v", name, err)
				s.Deregister(name)
				return
			}
		case "DER":
			name := string(bytes.Trim(msg[3:], " \t\n\r\f"))
			s.Deregister(name)
		default:
			s.Logger.Errorf("Unrecognized message: %q", msg)
			continue
		}
	}
}

func (s *TCPServer) Register(name string, conn net.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ctx := context.Background()
	ctx, canc := context.WithCancel(ctx)
	au := &Audience{
		Ctx:     ctx,
		Cancel:  canc,
		Alive:   true,
		Ch:      make(chan *Event),
		Conn:    conn,
		Last:    time.Now(),
		Name:    name,
		Timeout: 30 * time.Second,
		Logger:  s.Logger.WithField("audience", name),
	}
	s.Audiences[name] = au

	go func(s *TCPServer) {
		au.Logger.Info("Event channel starting...")
		for {
			select {
			case <-au.Ctx.Done():
				return
			case e := <-au.Ch:
				err := au.Notify(e)
				if err != nil {
					au.Logger.WithError(err).Errorf("Event %s error", e.String)
				}
			}
		}
	}(s)

	s.Logger.Infof("REG %s from %s", name, conn.RemoteAddr().String())

	return nil
}

func (s *TCPServer) Live(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	au, ok := s.Audiences[name]
	if !ok {
		return fmt.Errorf("Audience[%s] does not exists.", name)
	}
	au.Alive = true
	au.Last = time.Now()
	return nil
}

func (s *TCPServer) Deregister(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	au, ok := s.Audiences[name]
	if ok {
		delete(s.Audiences, name)
		au.Alive = false
		au.Cancel()
	}
	s.Logger.Error("DER ", name)
}

func (s *TCPServer) Heartbeat() {
	tick := time.Tick(time.Second)
	for t := range tick {
		for n, au := range s.Audiences {
			delta := t.Sub(au.Last)
			if delta > time.Second*30 {
				s.Deregister(n)
			} else if delta > time.Second*2 {
				s.Logger.Warnf("LIV %s INACTIVE", au.String())
				au.Alive = false
			}
		}
	}
}

func (s *TCPServer) Notify(e *Event) {
	s.Logger.Infof("Notify event[%d] key=%s value=%s", e.T, e.K, e.V)
	for _, au := range s.Audiences {
		if au.Alive {
			au.Ch <- e
		}
	}
}

func (s *TCPServer) EventLoop() {
	for e := range s.Ch {
		s.Notify(e)
	}
}
