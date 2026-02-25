package service

import (
	"errors"
	"hashGenerationService/internal/model"
	"hashGenerationService/internal/store"
	"hashGenerationService/internal/utils"
	"regexp"
)

const maxRetries = 5

var (
	ErrInvalidInput       = errors.New("input must be alphanumeric")
	ErrMaxRetriesExceeded = errors.New("failed to generate a unique hash after 5 retries")
	ErrHashNotFound       = errors.New("hash not found")

	alphanumericRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

type Service struct {
	store store.Store
}

func NewService(s store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) GenerateHash(input string) (*model.HashResponse, error) {
	if !alphanumericRe.MatchString(input) {
		return nil, ErrInvalidInput
	}

	for range maxRetries {
		hash, err := utils.GenerateHash(input)
		if err != nil {
			return nil, err
		}

		if !s.store.Exists(hash) {
			if err := s.store.Save(hash, input); err != nil {
				return nil, err
			}
			return &model.HashResponse{Input: input, Hash: hash}, nil
		}
	}

	return nil, ErrMaxRetriesExceeded
}
