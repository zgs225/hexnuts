package sync

import (
	"context"
	"io"
)

type Syncer interface {
	Sync(context.Context, io.Reader, io.Writer) error
}

type FileSyncer interface {
	Syncer

	SyncFile(context.Context, string, string) error
}
