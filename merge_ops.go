package cfglayer

import (
	"context"

	"github.com/LYH2263/go-cfglayer/internal/merge"
)

func (m *Merger) MergeStack(ctx context.Context) (map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := m.checkOpen(); err != nil {
		return nil, err
	}
	layers := layersFromStore(m.store.List())
	res := merge.MergeLayers(layersToMerge(layers), merge.Config{TreatEmptyAsDelete: m.opts.TreatEmptyAsDelete})
	m.mu.Lock()
	m.merges++
	m.mu.Unlock()
	if m.audit != nil {
		m.audit.Printf("merge layers=%d keys=%d", len(layers), len(res.Merged))
	}
	return merge.CopyMap(res.Merged), nil
}

func (m *Merger) Resolve(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := m.checkOpen(); err != nil {
		return "", err
	}
	if key == "" {
		return "", ErrBadInput
	}
	layers := layersToMerge(layersFromStore(m.store.List()))
	v, ok := merge.ResolveKey(layers, key, m.opts.TreatEmptyAsDelete)
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (m *Merger) ExplainKey(key string) ([]KeyStep, error) {
	if err := m.checkOpen(); err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrBadInput
	}
	layers := layersToMerge(layersFromStore(m.store.List()))
	steps := merge.ExplainChain(layers, key)
	if len(steps) == 0 {
		return nil, ErrNotFound
	}
	out := make([]KeyStep, len(steps))
	for i, s := range steps {
		out[i] = stepFromMerge(s)
	}
	return out, nil
}
