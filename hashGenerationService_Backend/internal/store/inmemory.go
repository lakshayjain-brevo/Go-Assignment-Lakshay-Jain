package store

import (
	"errors"
	"sync"
)

const maxEntries = 10_000

var ErrStoreFull = errors.New("store is at capacity")

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

func (s *InMemoryStore) SaveIfNotExists(hash, input string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[hash]; exists {
		return false, nil
	}
	if len(s.data) >= maxEntries {
		return false, ErrStoreFull
	}
	s.data[hash] = input
	return true, nil
}

func (s *InMemoryStore) Get(hash string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	input, ok := s.data[hash]
	if !ok {
		return "", ErrNotFound
	}
	return input, nil
}
