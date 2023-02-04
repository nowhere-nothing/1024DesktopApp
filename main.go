package main

import (
	"fmt"
	"github.com/jchv/go-webview2"
	"sync"
)

type App struct {
	w webview2.WebView
	sync.WaitGroup
	DownloadState
	once *sync.Once
}

func NewAPP(w webview2.WebView) *App {
	a := &App{
		w:    w,
		once: &sync.Once{},
	}
	a.onChange = a.setProgress
	return a
}

func (a *App) Once(f func()) {
	a.once.Do(f)
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
	bom := NewBOM()
	w := webview2.New(true)
	defer w.Destroy()
	w.SetSize(1000, 1200, webview2.HintNone)
	w.Init(combineJs())
	w.Navigate("https://t66y.com/index.php")

	app = NewAPP(w)
	app.Bind("emitEvent", bom.EmitEvent)
	app.Bind("download", downloadFunc(app))
	app.Bind("testFunc", testFunc(app))
	app.Bind("setSaveFolder", setSaveFolder)
	app.Bind("pickFolder", pickFolder)

	once := sync.Once{}
	bom.AddEventListener(DOMContentLoaded, func(eventType EventType, event Event) {
		once.Do(func() {
			app.RunJS("sendSaveFolder();")
		})
	})

	app.Run()
	app.Wait()
}
