package sync

import (
	"io"
)

type Pair struct {
	Dst io.Writer
	Src io.Reader
}
