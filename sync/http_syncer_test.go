package sync

import (
	"bytes"
	"context"
	"testing"

	"github.com/zgs225/hexnuts/client"
)

func TestSync(t *testing.T) {
	ctx := context.Background()
	text := `debug = "{% hello.world %}"`
	r := bytes.NewBufferString(text)
	w := new(bytes.Buffer)
	c := &client.HTTPClient{Addr: "http://localhost:5678"}
	s := &HTTPSyncer{Client: c, Symbols: make(map[string]string)}

	if err := s.Sync(ctx, r, w); err != nil {
		t.Error("Sync error: ", err)
	}

	if w.String() != `debug = "1"` {
		t.Errorf("Sync error:\n\texpect: %s\n\t   get: %s", `debug = "1"`, w.String())
	}
}
