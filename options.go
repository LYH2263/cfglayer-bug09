package cfglayer

type Options struct {
    NodeID              string
    MaxLayers           int
    TreatEmptyAsDelete  bool
    AuditPath           string
}

func (o Options) withDefaults() Options {
    if o.NodeID == "" {
        o.NodeID = "local"
    }
    if o.MaxLayers <= 0 {
        o.MaxLayers = 64
    }
    return o
}
