package monitor

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	stdsync "github.com/zgs225/hexnuts/sync"
)

type Client struct {
	Ctx        context.Context
	RemoteAddr net.Addr
	TLS        bool
	Name       string
	Conn       net.Conn
	Dialed     bool
	Syncer     stdsync.FileSyncer
	Pairs      map[string]*stdsync.Pair

	w  *bufio.Writer
	r  *bufio.Reader
	mu sync.Mutex
}

func (c *Client) Dial() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.Dialed {
		dialer := &net.Dialer{Timeout: 30 * time.Second}
		if c.TLS {
			conf := &tls.Config{InsecureSkipVerify: true}
			conn, err := tls.DialWithDialer(dialer, c.RemoteAddr.Network(), c.RemoteAddr.String(), conf)
			if err != nil {
				return err
			}
			c.Conn = conn
		} else {
			conn, err := dialer.DialContext(c.Ctx, c.RemoteAddr.Network(), c.RemoteAddr.String())
			if err != nil {
				return err
			}
			c.Conn = conn
		}
		c.Dialed = true

		c.w = bufio.NewWriter(c.Conn)
		c.r = bufio.NewReader(c.Conn)
	}

	return nil
}

func (c *Client) writer() *bufio.Writer {
	return c.w
}

func (c *Client) reader() *bufio.Reader {
	return c.r
}

func (c *Client) Register() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := c.writer().WriteString(fmt.Sprintf("REG %s\n", c.Name)); err != nil {
		return err
	}
	return c.writer().Flush()
}

func (c *Client) Live() error {
	c.mu.Lock()
	c.mu.Unlock()

	if _, err := c.writer().WriteString(fmt.Sprintf("LIV %s\n", c.Name)); err != nil {
		return err
	}
	return c.writer().Flush()
}

func (c *Client) Deregister() error {
	c.mu.Lock()
	c.mu.Unlock()

	if _, err := c.writer().WriteString(fmt.Sprintf("DER %s\n", c.Name)); err != nil {
		return err
	}
	return c.writer().Flush()
}

func (c *Client) ReadEvent() error {
	c.mu.Lock()
	c.mu.Unlock()

	data, err := c.reader().ReadBytes('\n')
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(string(data[0]), 10, 8)
	if err != nil {
		return err
	}

	kv := bytes.Split(bytes.Trim(data[1:], " \r\n"), []byte(" "))
	if len(kv) == 0 {
		return fmt.Errorf("不支持的消息: %q", data)
	}
	c.Syncer.DelSymbol(string(kv[0]))

	t := int8(i)
	switch t {
	case Events_ADD:
		return c.SyncPairs()
	case Events_Update:
		return c.SyncPairs()
	case Events_DEL:
		return c.SyncPairs()
	default:
		return fmt.Errorf("不支持的消息: %q", data)
	}
}

func (c *Client) SyncPairs() error {
	log.Println("Sync files...")
	ch := make(chan error)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	for _, p := range c.Pairs {
		wg.Add(1)
		go func(ch chan error, p *stdsync.Pair) {
			if err := c.Syncer.SyncFile(c.Ctx, p.Src, p.Dst); err != nil {
				ch <- err
			}
			wg.Done()
		}(ch, p)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case err := <-ch:
		return err
	case <-done:
		return nil
	}
}
