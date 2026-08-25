package cfglayer

import "time"

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

type Snapshot struct {
    NodeID string
    Layers []Layer
    Merged map[string]string
}

type Stats struct {
    LayersTotal int
    KeysTotal   int
    MergeCount  int64
}

type Health struct {
    OK     bool
    Closed bool
    NodeID string
}
