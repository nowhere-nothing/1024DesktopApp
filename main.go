package main

import (
	"fmt"
	"github.com/jchv/go-webview2"
	"path"
	"strings"
	"sync"
)

var saveFolder = `D:\temp`

type App struct {
	w webview2.WebView
	sync.WaitGroup
	DownloadState
}

func NewAPP(w webview2.WebView) *App {
	a := &App{
		w: w,
	}
	a.onChange = a.setProgress
	return a
}

func (a *App) RunJS(js string) {
	a.w.Dispatch(func() {
		a.w.Eval(js)
	})
}
func (a *App) Alert(str string) {
	a.w.Dispatch(func() {
		a.w.Eval(fmt.Sprintf(`alert("%s");`, str))
	})
}

func (a *App) setProgress(m, c int32) {
	if m == c {
		m = 0
		c = 0
		a.Reset()
	}
	a.RunJS(fmt.Sprintf("setGlobalProgress(%d, %d)", m, c))
}

func (a *App) Bind(name string, f any) {
	if err := a.w.Bind(name, f); err != nil {
		panic(err)
	}
}
func (a *App) Run() {
	a.w.Run()
}

var app *App

func main() {
	w := webview2.New(true)
	defer w.Destroy()
	w.SetSize(1000, 1200, webview2.HintNone)
	w.Init(combineJs())
	w.Navigate("https://t66y.com/index.php")

	app = NewAPP(w)
	app.Bind("download", download)
	app.Bind("testFunc", func() {
		app.RunJS(`console.log("Hello");`)
		app.RunJS(`alert("Hello");`)
	})

	app.Run()
	app.Wait()
}

func download(title, url string, images []string) {
	title = strings.TrimSpace(title)
	fn, err := DownloadFunc(path.Join(saveFolder, title))
	if err != nil {
		app.Alert(fmt.Sprintf("%s", err))
		return
	}
	go fn(url, images)
}
