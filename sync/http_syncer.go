package sync

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/zgs225/hexnuts/client"
)

const (
	Symbols_START = "{%"
	Symbols_END   = "%}"
)

type HTTPSyncer struct {
	Client  client.Client
	Symbols map[string]string
}

func (hs *HTTPSyncer) SyncFile(ctx context.Context, src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()
	return hs.Sync(ctx, r, w)
}

func (hs *HTTPSyncer) Sync(ctx context.Context, r io.Reader, w io.Writer) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return hs.sync(r, w)
	}
}

func (hs *HTTPSyncer) sync(r io.Reader, w io.Writer) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	todos := hs.getSyncContexts(data)
	if todos != nil {
		for _, sc := range todos {
			if data, err = sc.Do(data); err != nil {
				return err
			}
		}
	}
	_, err = w.Write(data)
	return err
}

func (hs *HTTPSyncer) getSyncContexts(data []byte) []*SyncContext {
	ctx := context.Background()
	pattern, _ := regexp.Compile(fmt.Sprintf("%s\\s([_a-zA-Z][\\w\\.]*)\\s*%s", Symbols_START, Symbols_END))
	matches := pattern.FindAllSubmatch(data, -1)
	if len(matches) == 0 {
		return nil
	}
	rv := make([]*SyncContext, 0, len(matches))
	for _, group := range matches {
		if len(group) != 2 {
			continue
		}
		sc := &SyncContext{
			Client:  hs.Client,
			Ctx:     ctx,
			Key:     group[1],
			Replace: group[0],
			Symbols: hs.Symbols,
		}
		rv = append(rv, sc)
	}
	return rv
}

func (hs *HTTPSyncer) DelSymbol(key string) {
	delete(hs.Symbols, key)
}
