package store

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestSaveIfNotExists(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*InMemoryStore)
		hash      string
		input     string
		wantSaved bool
		wantErr   error
	}{
		{
			name:      "saves new entry",
			hash:      "abc123",
			input:     "hello",
			wantSaved: true,
		},
		{
			name: "returns false on collision",
			setup: func(s *InMemoryStore) {
				s.data["abc123"] = "hello"
			},
			hash:      "abc123",
			input:     "world",
			wantSaved: false,
		},
		{
			name: "returns ErrStoreFull when at capacity",
			setup: func(s *InMemoryStore) {
				for i := range maxEntries {
					s.data[fmt.Sprintf("key%d", i)] = "x"
				}
			},
			hash:    "newhash",
			input:   "newval",
			wantErr: ErrStoreFull,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewInMemoryStore()
			if tc.setup != nil {
				tc.setup(s)
			}
			saved, err := s.SaveIfNotExists(tc.hash, tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if saved != tc.wantSaved {
				t.Errorf("saved = %v, want %v", saved, tc.wantSaved)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*InMemoryStore)
		hash      string
		wantInput string
		wantErr   error
	}{
		{
			name: "returns input for existing hash",
			setup: func(s *InMemoryStore) {
				s.data["abc123"] = "hello"
			},
			hash:      "abc123",
			wantInput: "hello",
		},
		{
			name:    "returns ErrNotFound for unknown hash",
			hash:    "unknown",
			wantErr: ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewInMemoryStore()
			if tc.setup != nil {
				tc.setup(s)
			}
			input, err := s.Get(tc.hash)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if input != tc.wantInput {
				t.Errorf("input = %q, want %q", input, tc.wantInput)
			}
		})
	}
}

func TestSaveIfNotExistsConcurrent(t *testing.T) {
	s := NewInMemoryStore()
	const goroutines = 100
	saved := make([]bool, goroutines)
	var wg sync.WaitGroup

	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ok, err := s.SaveIfNotExists("samehash", "input")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			saved[i] = ok
		}(i)
	}
	wg.Wait()

	count := 0
	for _, ok := range saved {
		if ok {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 successful save, got %d", count)
	}
}
