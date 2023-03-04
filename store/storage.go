package store

type Storage interface {
	Save(path, name string, data []byte) error
	MkdirAll(path string) error
	Exist(file string) bool
}
