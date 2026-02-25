package service

import (
	"errors"
	"fmt"
	"hashGenerationService/internal/model"
	"hashGenerationService/internal/store"
	"hashGenerationService/internal/utils"
	"regexp"
)

const (
	maxRetries  = 5
	maxInputLen = 256
)

var (
	ErrInvalidInput       = errors.New("input must be alphanumeric")
	ErrMaxRetriesExceeded = errors.New("failed to generate a unique hash after 5 retries")
	ErrHashNotFound       = errors.New("hash not found")
	ErrStoreFull          = fmt.Errorf("service is at capacity: %w", store.ErrStoreFull)

	alphanumericRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

type Service struct {
	store store.Store
}

func NewService(s store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) GenerateHash(input string) (*model.HashResponse, error) {
	if len(input) > maxInputLen {
		return nil, ErrInvalidInput
	}
	if !alphanumericRe.MatchString(input) {
		return nil, ErrInvalidInput
	}

	for range maxRetries {
		hash, err := utils.GenerateHash(input)
		if err != nil {
			return nil, err
		}

		saved, err := s.store.SaveIfNotExists(hash, input)
		if err != nil {
			if errors.Is(err, store.ErrStoreFull) {
				return nil, ErrStoreFull
			}
			return nil, err
		}
		if saved {
			return &model.HashResponse{Input: input, Hash: hash}, nil
		}
	}

	return nil, ErrMaxRetriesExceeded
}

func (s *Service) GetHash(hash string) (*model.HashResponse, error) {
	input, err := s.store.Get(hash)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrHashNotFound
		}
		return nil, err
	}
	return &model.HashResponse{Input: input, Hash: hash}, nil
}
