package store

import (
	"fmt"
	"os"
	"path/filepath"
)

type fsStorage struct {
	dir string
}

func NewFsStorage(dir string) Storage {
	return &fsStorage{
		dir: dir,
	}
}

func (f *fsStorage) Save(path, name string, data []byte) error {
	var cnt int
	fullpath := filepath.Join(path, name)
	for {
		st, err := os.Stat(fullpath)
		if err != nil { // os.PathError.Err == syscall.ENOSPC
			return os.WriteFile(fullpath, data, 0644)
		}
		if st.IsDir() {
			return fmt.Errorf("%s is dir", fullpath)
		}
		// conflict
		fullpath = filepath.Join(path, fmt.Sprintf("%d%s", cnt, name))
		cnt++
	}
}

func (f *fsStorage) MkdirAll(path string) error {
	return os.MkdirAll(path, 0644)
}

func (f *fsStorage) Exist(name string) bool {
	_, err := os.Stat(name)
	return err != nil
}
