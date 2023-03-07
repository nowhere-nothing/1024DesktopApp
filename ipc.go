package main

import (
	"github.com/gookit/goutil/fmtutil"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"strings"
	"time"
	"webview_demo/config"
	"webview_demo/downloader"
	"webview_demo/store"
)

type Fetcher struct {
	fetchHandler *downloader.Downloader
	progress     *downloader.Progress
	config       *config.Config
	pool         *ants.Pool
	store.Storage
	*logrus.Logger
}

func (f *Fetcher) Fetch(title, url string, images []string) {
	title = strings.Trim(title, "!@#$%^&*()_+-=`~,.<>/?{}\\| ")

	f.Debugf("start download %s from url %s", title, url)

	if err := f.pool.Submit(func() {

		f.batchDownload(title, url, images)

	}); err != nil {
		f.WithError(err).WithFields(map[string]interface{}{
			"title": title, "url": url,
			"images": images,
		}).Errorln("add download task")
	}
}

func (f *Fetcher) batchDownload(title, url string, images []string) {
	cont := store.PostContent{
		Title:  title,
		Url:    url,
		Images: images,
	}

	if err := f.MkdirAll(&cont); err != nil {
		f.WithError(err).Errorln("mkdir all")
		return
	}

	for i, t := range images {
		idx := i
		target := t

		logger := f.WithFields(map[string]interface{}{
			"index": idx, "title": title, "url": url, "subUrl": target,
		})

		if err := f.pool.Submit(func() {
			f.progress.AddCur(1)

			f.fetchAndSave(logger, &cont, target)

		}); err != nil {
			logger.WithError(err).Errorln("add sub image download task")
			if err = f.SaveFailed(&cont, target); err != nil {
				logger.WithError(err).Errorln("save failed url")
			}
		} else {
			f.progress.AddMax(1)
		}
	}
}

func (f *Fetcher) fetchAndSave(logger *logrus.Entry, pc *store.PostContent, url string) {
	start := time.Now()
	f.Debugln("start download", url)

	fd, err := f.fetchHandler.Fetch(url)
	if err != nil {
		logger.WithError(err).Errorln("fetch image")
		return
	}

	err = f.Save(pc, &store.PostImage{
		Name:        fd.Name(),
		Data:        fd.Data,
		ContentType: fd.ContentType,
		Url:         url,
	})
	if err != nil {
		logger.WithError(err).Errorln("save image")
		return
	}

	f.Debugf("download %s size %s time %.2fs", url,
		fmtutil.DataSize(uint64(len(fd.Data))),
		time.Since(start).Seconds(),
	)
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
