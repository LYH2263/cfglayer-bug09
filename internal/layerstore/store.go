package layerstore

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrDuplicate = errors.New("layerstore: duplicate layer id")
	ErrFull      = errors.New("layerstore: stack full")
	ErrEmpty     = errors.New("layerstore: empty stack")
)

type Layer struct {
	ID       string
	Priority int
	Source   string
	Values   map[string]string
	Created  time.Time
}

type Store struct {
	mu       sync.RWMutex
	stack    []Layer
	byID     map[string]int
	max      int
	keyCount int
}

func New(max int) *Store {
	return &Store{
		stack: make([]Layer, 0, 8),
		byID:  make(map[string]int),
		max:   max,
	}
}

func (s *Store) Push(layer Layer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, dup := s.byID[layer.ID]; dup {
		return ErrDuplicate
	}
	if len(s.stack) >= s.max {
		return ErrFull
	}
	idx := len(s.stack)
	s.stack = append(s.stack, layer)
	s.byID[layer.ID] = idx
	s.keyCount += len(layer.Values)
	return nil
}

func (s *Store) Pop() (Layer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.stack) == 0 {
		return Layer{}, ErrEmpty
	}
	last := len(s.stack) - 1
	layer := s.stack[last]
	s.stack = s.stack[:last]
	delete(s.byID, layer.ID)
	s.keyCount -= len(layer.Values)
	if s.keyCount < 0 {
		s.keyCount = 0
	}
	return layer, nil
}

func (s *Store) List() []Layer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Layer, len(s.stack))
	copy(out, s.stack)
	return out
}

func (s *Store) Get(id string) (Layer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	idx, ok := s.byID[id]
	if !ok {
		return Layer{}, false
	}
	return s.stack[idx], true
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.stack)
}

func (s *Store) KeyCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.keyCount
}

func (s *Store) ReplaceStack(layers []Layer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stack = make([]Layer, len(layers))
	copy(s.stack, layers)
	s.byID = make(map[string]int, len(layers))
	s.keyCount = 0
	for i, l := range s.stack {
		s.byID[l.ID] = i
		s.keyCount += len(l.Values)
	}
}

func (s *Store) Snapshot() []Layer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Layer, len(s.stack))
	for i, l := range s.stack {
		vals := make(map[string]string, len(l.Values))
		for k, v := range l.Values {
			vals[k] = v
		}
		out[i] = Layer{
			ID:       l.ID,
			Priority: l.Priority,
			Source:   l.Source,
			Values:   vals,
			Created:  l.Created,
		}
	}
	for i := range s.stack {
		s.stack[i].Values = map[string]string{}
	}
	return out
}
