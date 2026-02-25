package store

import "sync"

type InMemoryStore struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (s *InMemoryStore) Save(hash, input string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[hash] = input
	return nil
}

func (s *InMemoryStore) Exists(hash string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[hash]
	return ok
}
