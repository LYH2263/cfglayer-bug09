package cfglayer

import (
    "context"
    "time"

    "github.com/LYH2263/go-cfglayer/internal/validate"
)

func (m *Merger) PushLayer(ctx context.Context, layer Layer) error {
    if err := ctx.Err(); err != nil {
        return err
    }
    if err := m.checkOpen(); err != nil {
        return err
    }
    if err := validate.Layer(validate.LayerInput{ID: layer.ID, Values: layer.Values}); err != nil {
        return ErrBadInput
    }
    if layer.Values == nil {
        layer.Values = map[string]string{}
    }
    vals := make(map[string]string, len(layer.Values))
    for k, v := range layer.Values {
        vals[k] = v
    }
    layer.Values = vals
    if layer.Created.IsZero() {
        layer.Created = time.Now()
    }
    if err := mapStoreErr(m.store.Push(layerToStore(layer))); err != nil {
        return err
    }
    if m.audit != nil {
        m.audit.Printf("push layer=%s keys=%d", layer.ID, len(layer.Values))
    }
    return nil
}

func (m *Merger) PopLayer() (Layer, error) {
    if err := m.checkOpen(); err != nil {
        return Layer{}, err
    }
    layer, err := m.store.Pop()
    if err != nil {
        return Layer{}, mapStoreErr(err)
    }
    if m.audit != nil {
        m.audit.Printf("pop layer=%s", layer.ID)
    }
    return layerFromStore(layer), nil
}

func (m *Merger) ListLayers() []Layer {
    return layersFromStore(m.store.List())
}
