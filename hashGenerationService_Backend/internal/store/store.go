package store

type Store interface {
	Save(hash, input string) error
	Exists(hash string) bool
}
