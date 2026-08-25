package merge

type Tracker struct {
    steps map[string][]KeyStep
}

func NewTracker() *Tracker {
    return &Tracker{steps: make(map[string][]KeyStep)}
}

func (t *Tracker) Record(key string, step KeyStep) {
    t.steps[key] = append(t.steps[key], step)
}

func (t *Tracker) Chain(key string) []KeyStep {
    src := t.steps[key]
    out := make([]KeyStep, len(src))
    copy(out, src)
    return out
}

func (t *Tracker) Winner(key string) (KeyStep, bool) {
    chain := t.steps[key]
    if len(chain) == 0 {
        return KeyStep{}, false
    }
    return chain[len(chain)-1], true
}

func (t *Tracker) Keys() []string {
    keys := make([]string, 0, len(t.steps))
    for k := range t.steps {
        keys = append(keys, k)
    }
    return keys
}

func (t *Tracker) Reset() {
    t.steps = make(map[string][]KeyStep)
}

func BuildTracker(layers []Layer) *Tracker {
    t := NewTracker()
    res := MergeLayers(layers, Config{})
    for k, step := range res.Provenance {
        t.steps[k] = ExplainChain(layers, k)
        if len(t.steps[k]) == 0 {
            t.steps[k] = []KeyStep{step}
        }
    }
    return t
}
