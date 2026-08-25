package cfglayer

import (
	"sync"

	"github.com/LYH2263/go-cfglayer/internal/audit"
	"github.com/LYH2263/go-cfglayer/internal/layerstore"
)

type Merger struct {
	opts       Options
	store      *layerstore.Store
	audit      *audit.Logger
	mu         sync.RWMutex
	closed     bool
	merges     int64
}

func New(opts Options) (*Merger, error) {
	o := opts.withDefaults()
	m := &Merger{
		opts:  o,
		store: layerstore.New(o.MaxLayers),
	}
	if o.AuditPath != "" {
		lg, err := audit.Open(o.AuditPath)
		if err != nil {
			return nil, err
		}
		m.audit = lg
	}
	return m, nil
}

func (m *Merger) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	if m.audit != nil {
		return m.audit.Close()
	}
	return nil
}

func (m *Merger) checkOpen() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.closed {
		return ErrClosed
	}
	return nil
}
