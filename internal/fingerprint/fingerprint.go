package fingerprint

import (
    "crypto/sha256"
    "encoding/hex"
    "sort"
)

func LayerHash(values map[string]string) string {
    if len(values) == 0 {
        return hex.EncodeToString(sha256.New().Sum(nil))
    }
    keys := make([]string, 0, len(values))
    for k := range values {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    h := sha256.New()
    for _, k := range keys {
        h.Write([]byte(k))
        h.Write([]byte{0})
        h.Write([]byte(values[k]))
        h.Write([]byte{0})
    }
    return hex.EncodeToString(h.Sum(nil))
}

func StackHash(layers []struct {
    ID     string
    Values map[string]string
}) string {
    h := sha256.New()
    for _, l := range layers {
        h.Write([]byte(l.ID))
        h.Write([]byte(LayerHash(l.Values)))
    }
    return hex.EncodeToString(h.Sum(nil))
}
