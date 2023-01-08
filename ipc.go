package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var saveFolder string

func downloadFunc(app *App) func(title, url string, images []string) {
	return func(title, url string, images []string) {
		if saveFolder == "" {
			app.Alert("图片保存路径没有设置")
			return
		}
		title = strings.TrimSpace(title)
		fn, err := Downloader(path.Join(saveFolder, title))
		if err != nil {
			app.Alert(fmt.Sprintf("%s", err))
			return
		}
		go fn(url, images)
	}
}

func testFunc(app *App) func() {
	return func() {
		app.RunJS(`console.log("Hello");`)
		app.RunJS(`alert("Hello");`)
	}
}

func initSaveFolder(p string) {
	saveFolder = p
}

func setSaveFolder(p string) error {
	err := os.MkdirAll(p, 0644)
	if err != nil {
		return err
	}
	saveFolder = p
	return nil
}
