package main

import (
	"fmt"
	"github.com/sqweek/dialog"
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
		// todo: trim special character
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

func setSaveFolder(p string) error {
	err := os.MkdirAll(p, 0644)
	if err != nil {
		return err
	}
	saveFolder = p
	return nil
}

func pickFolder(title, start string) (string, error) {
	return dialog.Directory().
		Title(title).
		SetStartDir(start).
		Browse()
}

type Event int
type EventType string
type EventHandler func(EventType, Event)

const (
	DOMContentLoaded EventType = "DOMContentLoaded"
)

type BOM struct {
	eventHandler map[EventType][]EventHandler
}

func NewBOM() *BOM {
	return &BOM{
		eventHandler: make(map[EventType][]EventHandler),
	}
}

func (b *BOM) AddEventListener(t EventType, f EventHandler) {
	b.eventHandler[t] = append(b.eventHandler[t], f)
}

func (b *BOM) EmitEvent(t EventType) {
	println("emit", t)
	for _, f := range b.eventHandler[t] {
		f(t, 0)
	}
}
