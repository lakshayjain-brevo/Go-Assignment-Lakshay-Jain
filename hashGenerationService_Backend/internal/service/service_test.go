package service

import (
	"errors"
	"hashGenerationService/internal/model"
	"hashGenerationService/internal/store"
	"testing"
)

// mockStore implements store.Store for testing without real storage.
type mockStore struct {
	data      map[string]string
	storeFull bool
}

func newMockStore() *mockStore {
	return &mockStore{data: make(map[string]string)}
}

func (m *mockStore) Save(hash, input string) error {
	m.data[hash] = input
	return nil
}

func (m *mockStore) SaveIfNotExists(hash, input string) (bool, error) {
	if m.storeFull {
		return false, store.ErrStoreFull
	}
	if _, exists := m.data[hash]; exists {
		return false, nil
	}
	m.data[hash] = input
	return true, nil
}

func (m *mockStore) Get(hash string) (string, error) {
	input, ok := m.data[hash]
	if !ok {
		return "", store.ErrNotFound
	}
	return input, nil
}

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:  "valid alphanumeric input",
			input: "hello123",
		},
		{
			name:    "empty input rejected",
			input:   "",
			wantErr: ErrInvalidInput,
		},
		{
			name:    "input with spaces rejected",
			input:   "hello world",
			wantErr: ErrInvalidInput,
		},
		{
			name:    "input with special chars rejected",
			input:   "hello!",
			wantErr: ErrInvalidInput,
		},
		{
			name:    "input exceeding max length rejected",
			input:   string(make([]byte, maxInputLen+1)),
			wantErr: ErrInvalidInput,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(newMockStore())
			resp, err := svc.GenerateHash(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil {
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Input != tc.input {
					t.Errorf("resp.Input = %q, want %q", resp.Input, tc.input)
				}
				if len(resp.Hash) != 10 {
					t.Errorf("hash length = %d, want 10", len(resp.Hash))
				}
			}
		})
	}
}

func TestGetHash(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockStore)
		hash     string
		wantResp *model.HashResponse
		wantErr  error
	}{
		{
			name: "returns response for existing hash",
			setup: func(m *mockStore) {
				m.data["abc123def0"] = "hello"
			},
			hash:     "abc123def0",
			wantResp: &model.HashResponse{Input: "hello", Hash: "abc123def0"},
		},
		{
			name:    "returns ErrHashNotFound for unknown hash",
			hash:    "doesnotexist",
			wantErr: ErrHashNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := newMockStore()
			if tc.setup != nil {
				tc.setup(ms)
			}
			svc := NewService(ms)
			resp, err := svc.GetHash(tc.hash)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if tc.wantResp != nil {
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Input != tc.wantResp.Input || resp.Hash != tc.wantResp.Hash {
					t.Errorf("resp = %+v, want %+v", resp, tc.wantResp)
				}
			}
		})
	}
}
