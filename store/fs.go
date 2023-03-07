package store

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

type fsStorage struct {
	basePath string
}

func NewFsStorage(dir string) Storage {
	y, m, d := time.Now().Date()
	prefix := path.Join(dir,
		fmt.Sprintf("%d-%d-%d", y, m, d),
	)
	return &fsStorage{
		basePath: prefix,
	}
}

func (f *fsStorage) Save(pc *PostContent, pi *PostImage) error {
	var cnt int
	fullpath := filepath.Join(f.basePath, pc.Title, pi.Name)
	for {
		st, err := os.Stat(fullpath)
		if err != nil { // os.PathError.Err == syscall.ENOSPC
			return os.WriteFile(fullpath, pi.Data, 0644)
		}
		if st.IsDir() {
			return fmt.Errorf("%s is dir", fullpath)
		}
		// conflict
		fullpath = filepath.Join(f.basePath, pc.Title, fmt.Sprintf("%d%s", cnt, pi.Name))
		cnt++
	}
}

func (f *fsStorage) SaveFailed(pc *PostContent, items string) error {

	return nil
}

func (f *fsStorage) MkdirAll(pc *PostContent) error {
	return os.MkdirAll(filepath.Join(f.basePath, pc.Title), 0644)
}
