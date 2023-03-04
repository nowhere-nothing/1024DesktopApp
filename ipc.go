package main

import (
	"fmt"
	"github.com/gookit/goutil/fmtutil"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"path"
	"strings"
	"sync"
	"time"
	"webview_demo/config"
	"webview_demo/downloader"
	"webview_demo/store"
)

type Wrapper struct {
	D      *downloader.Downloader
	P      *downloader.Progress
	C      *config.Config
	DH     *store.DownloadHistory
	S      store.Storage
	Pool   *ants.Pool
	Logger *logrus.Logger
}

func DownloadFunc(app *App) func(title, url string, images []string) {
	return func(title, url string, images []string) {
		if app.C.SavePath == "" {

		}
		w := &Wrapper{
			D:      app.D,
			P:      app.P,
			C:      app.C,
			S:      app.S,
			Pool:   app.Pool,
			Logger: app.logger,
			DH:     app.DH,
		}

		title = strings.Trim(title, "!@#$%^&*()_+-=`~,.<>/?{}\\| ")

		if err := app.Pool.Submit(func() {
			batchDownload(title, url, images, w)
		}); err != nil {
			app.logger.WithError(err).WithFields(map[string]interface{}{
				"title": title, "url": url,
				"images": images,
			}).Errorf("add images download task")
		}
	}
}

func fetchAndSave(w *Wrapper, url, path string) error {
	start := time.Now()
	w.Logger.Println("start fetch", url)
	fd, err := w.D.Fetch(url)
	if err != nil {
		return err
	}

	w.Logger.Printf("fetch  done %s size %s time %.2fs",
		url, fmtutil.DataSize(uint64(len(fd.Data))), time.Since(start).Seconds())

	err = w.S.Save(path, fd.Name(), fd.Data)
	if err != nil {
		return err
	}
	return nil
}

type SyncList[T any] struct {
	mu   sync.Mutex
	list []T
}

func NewSyncList[T any]() *SyncList[T] {
	return &SyncList[T]{
		list: make([]T, 0),
	}
}

func (sl *SyncList[T]) Add(v T) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list = append(sl.list, v)
}

func (sl *SyncList[T]) List() []T {
	return sl.list
}

func batchDownload(title, url string, images []string, w *Wrapper) {
	w.DH.Save(title, url, images)

	y, m, day := time.Now().Date()
	prefix := path.Join(w.C.SavePath,
		fmt.Sprintf("%d-%d-%d", y, m, day),
		title,
	)
	if err := w.S.MkdirAll(prefix); err != nil {
		w.Logger.WithError(err).Errorf("create directory [%s]", prefix)
		return
	}

	wg := sync.WaitGroup{}
	failed := NewSyncList[string]()
	var added int32

	for i, target := range images {
		idx := i
		subTarget := target

		logger := w.Logger.WithFields(map[string]interface{}{
			"index": idx, "title": title, "url": url,
			"subUrl": subTarget,
		})

		if err := w.Pool.Submit(func() {
			if err := fetchAndSave(w, subTarget, prefix); err != nil {
				logger.WithError(err).Errorf("fetch or save")
				failed.Add(subTarget)
			}
			wg.Done()
			w.P.AddCur(1)
		}); err != nil {
			logger.WithError(err).Errorf("add sub image download task")
			failed.Add(subTarget)
		} else {
			wg.Add(1)
			added++
		}
	}
	w.P.AddMax(added)
	wg.Wait()
	w.DH.Failed(title, url, failed.List())
}

func testFunc(app *App) func() {
	return func() {
		app.RunJS(`console.log("Hello");`)
		app.RunJS(`alert("Hello");`)
	}
}

func FolderPicker(title, start string) (string, error) {
	return dialog.Directory().
		Title(title).
		SetStartDir(start).
		Browse()
}
