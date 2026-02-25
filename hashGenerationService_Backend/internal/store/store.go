package store

import "errors"

var ErrNotFound = errors.New("hash not found")

type Store interface {
	SaveIfNotExists(hash, input string) (bool, error)
	Get(hash string) (string, error)
}
