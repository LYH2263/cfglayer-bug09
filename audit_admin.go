package cfglayer

func (m *Merger) RotateAudit() error {
	if m.audit == nil {
		return nil
	}
	return m.audit.Rotate()
}
