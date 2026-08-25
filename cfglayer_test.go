package cfglayer_test

import (
	"context"
	"testing"

	cl "github.com/LYH2263/go-cfglayer"
)

func newMerger(t *testing.T) *cl.Merger {
	t.Helper()
	m, err := cl.New(cl.Options{MaxLayers: 32})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = m.Close() })
	return m
}

func push(t *testing.T, m *cl.Merger, id string, vals map[string]string) {
	t.Helper()
	if err := m.PushLayer(context.Background(), cl.Layer{ID: id, Values: vals}); err != nil {
		t.Fatal(err)
	}
}
