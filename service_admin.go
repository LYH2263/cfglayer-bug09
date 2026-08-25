package cfglayer

func (m *Merger) Health() Health {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return Health{OK: !m.closed, Closed: m.closed, NodeID: m.opts.NodeID}
}

func (m *Merger) Stats() Stats {
    return Stats{
        LayersTotal: m.store.Len(),
        KeysTotal:   m.store.KeyCount(),
        MergeCount:  m.merges,
    }
}
