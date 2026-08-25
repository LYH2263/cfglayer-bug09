package provenance

import (
    "fmt"
    "strings"

    "github.com/LYH2263/go-cfglayer"
)

type Entry struct {
    Key   string
    Steps []cfglayer.KeyStep
}

func FormatExplain(key string, steps []cfglayer.KeyStep) string {
    if len(steps) == 0 {
        return fmt.Sprintf("%s: (not set)", key)
    }
    var b strings.Builder
    b.WriteString(key)
    b.WriteString(":\n")
    for i, s := range steps {
        b.WriteString(fmt.Sprintf("  [%d] layer=%s prio=%d source=%s value=%q\n",
            i, s.LayerID, s.Priority, s.Source, s.Value))
    }
    if w := steps[len(steps)-1]; w.Value != "" {
        b.WriteString(fmt.Sprintf("  => winner: %q from %s\n", w.Value, w.LayerID))
    }
    return b.String()
}

func CollectConflicts(layers []cfglayer.Layer) []Entry {
    seen := make(map[string][]cfglayer.KeyStep)
    for _, layer := range layers {
        for k, v := range layer.Values {
            seen[k] = append(seen[k], cfglayer.KeyStep{
                LayerID:  layer.ID,
                Priority: layer.Priority,
                Value:    v,
                Source:   layer.Source,
            })
        }
    }
    var out []Entry
    for k, steps := range seen {
        if len(steps) > 1 {
            out = append(out, Entry{Key: k, Steps: steps})
        }
    }
    return out
}
