package cfglayer

import (
	"time"

	"github.com/LYH2263/go-cfglayer/internal/layerstore"
	"github.com/LYH2263/go-cfglayer/internal/merge"
)

func (m *Merger) ExportMerged() (Snapshot, error) {
	if err := m.checkOpen(); err != nil {
		return Snapshot{}, err
	}
	layers := layersFromStore(m.store.Snapshot())
	res := merge.MergeLayers(layersToMerge(layers), merge.Config{TreatEmptyAsDelete: m.opts.TreatEmptyAsDelete})
	snap := Snapshot{
		NodeID: m.opts.NodeID,
		Layers: layers,
		Merged: merge.CopyMap(res.Merged),
	}
	if m.audit != nil {
		m.audit.Printf("export merged keys=%d layers=%d", len(snap.Merged), len(layers))
	}
	return snap, nil
}

func (m *Merger) ImportSnapshot(snap Snapshot) error {
	if err := m.checkOpen(); err != nil {
		return err
	}
	layers := make([]layerstore.Layer, len(snap.Layers))
	for i, l := range snap.Layers {
		layers[i] = layerToStore(l)
		if layers[i].Created.IsZero() {
			layers[i].Created = time.Now()
		}
	}
	m.store.ReplaceStack(layers)
	if m.audit != nil {
		m.audit.Printf("import snapshot layers=%d", len(layers))
	}
	return nil
}
