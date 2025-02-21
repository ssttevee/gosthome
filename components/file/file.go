package file

import (
	"context"
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/gosthome/gosthome/core/component"
	"github.com/gosthome/gosthome/core/component/cid"
)

type File struct {
	cid.CID
	watcher *fsnotify.Watcher
}

func New(ctx context.Context, cfg *Config) ([]component.Component, error) {
	return []component.Component{&File{
		CID: cid.NewID("file"),
	}}, nil
}

// Setup implements component.Component.
func (f *File) Setup() {
	var err error
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to initialize file component", "err", err)
	}
}

// Close implements component.Component.
func (f *File) Close() error {
	if f.watcher != nil {
		f.watcher.Close()
	}
	return nil
}

// InitializationPriority implements component.Component.
func (c *File) InitializationPriority() component.InitializationPriority {
	return component.InitializationPriorityHardware
}

var _ component.Component = (*File)(nil)
