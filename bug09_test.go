package cfglayer_test

import (
	"testing"
)

func TestBug09_ExportPreservesLayerKeys(t *testing.T) {
	m := newMerger(t)
	push(t, m, "e1", map[string]string{"k": "v"})
	a, err := m.ExportMerged()
	if err != nil {
		t.Fatal(err)
	}
	b, err := m.ExportMerged()
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Layers[0].Values) != len(b.Layers[0].Values) {
		t.Fatalf("keys %d vs %d", len(a.Layers[0].Values), len(b.Layers[0].Values))
	}
}
