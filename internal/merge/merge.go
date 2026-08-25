package merge

import (
    "sort"
    "time"
)

type Layer struct {
    ID       string
    Priority int
    Source   string
    Values   map[string]string
    Created  time.Time
}

type KeyStep struct {
    LayerID  string
    Priority int
    Value    string
    Source   string
}

type Result struct {
    Merged     map[string]string
    Provenance map[string]KeyStep
}

type Config struct {
    TreatEmptyAsDelete bool
}

// MergeLayers applies bottom-to-top overlay; higher priority wins on conflict.
func MergeLayers(layers []Layer, cfg Config) Result {
    if len(layers) == 0 {
        return Result{
            Merged:     map[string]string{},
            Provenance: map[string]KeyStep{},
        }
    }
    ordered := make([]Layer, len(layers))
    copy(ordered, layers)
    sort.SliceStable(ordered, func(i, j int) bool {
        if ordered[i].Priority == ordered[j].Priority {
            return i < j
        }
        return ordered[i].Priority < ordered[j].Priority
    })

    merged := make(map[string]string)
    prov := make(map[string]KeyStep)

    for _, layer := range ordered {
        for k, v := range layer.Values {
            step := KeyStep{
                LayerID:  layer.ID,
                Priority: layer.Priority,
                Value:    v,
                Source:   layer.Source,
            }
            if cfg.TreatEmptyAsDelete && v == "" {
                delete(merged, k)
                prov[k] = step
                continue
            }
            merged[k] = v
            prov[k] = step
        }
    }
    return Result{Merged: merged, Provenance: prov}
}

func ExplainChain(layers []Layer, key string) []KeyStep {
    ordered := make([]Layer, len(layers))
    copy(ordered, layers)
    sort.SliceStable(ordered, func(i, j int) bool {
        if ordered[i].Priority == ordered[j].Priority {
            return i < j
        }
        return ordered[i].Priority < ordered[j].Priority
    })
    var steps []KeyStep
    for _, layer := range ordered {
        v, ok := layer.Values[key]
        if !ok {
            continue
        }
        steps = append(steps, KeyStep{
            LayerID:  layer.ID,
            Priority: layer.Priority,
            Value:    v,
            Source:   layer.Source,
        })
    }
    return steps
}

func ResolveKey(layers []Layer, key string, treatEmptyAsDelete bool) (string, bool) {
    res := MergeLayers(layers, Config{TreatEmptyAsDelete: treatEmptyAsDelete})
    v, ok := res.Merged[key]
    return v, ok
}

func CopyMap(src map[string]string) map[string]string {
    if src == nil {
        return map[string]string{}
    }
    out := make(map[string]string, len(src))
    for k, v := range src {
        out[k] = v
    }
    return out
}
