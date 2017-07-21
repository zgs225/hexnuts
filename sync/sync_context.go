package sync

import (
	"bytes"
	"context"

	"git.youplus.cc/tiny/hexnuts/client"
)

type SyncContext struct {
	Ctx     context.Context
	Client  client.Client
	Symbols map[string]string
	Key     []byte
	Replace []byte
	Value   string
	Got     bool
}

func (sc *SyncContext) Do(data []byte) ([]byte, error) {
	if err := sc.FetchValue(); err != nil {
		return data, err
	}
	return bytes.Replace(data, sc.Replace, []byte(sc.Value), 1), nil
}

func (sc *SyncContext) FetchValue() error {
	if sc.Got {
		return nil
	}
	k := string(sc.Key)
	if v, ok := sc.Symbols[k]; ok {
		sc.Value = v
		sc.Got = true
		return nil
	}
	v, err := sc.Client.Get(k)
	if err != nil {
		return err
	}
	sc.Value = v
	sc.Symbols[k] = v
	sc.Got = true
	return nil
}
