package main

import (
	"fmt"
	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"github.com/ncruces/zenity"
	"os"
	"reflect"
	"sync"
	"unsafe"
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
	urls := []string{
		"https://t66y.com/index.php", // .tpc_cont img
		"https://www.4kup.net/",      // #gallery img
		"http://ryuryu.tw/#/portal/signup",
	}
	url, err := zenity.List("选择网址", urls, zenity.Title("1024"))
	if err != nil {
		os.Exit(1)
	}
	if url == "" {
		url, err = zenity.Entry("输入网址", zenity.Title("1024"))
		if err != nil {
			os.Exit(1)
		}
	}

	bom := NewBOM()
	w := webview2.New(true)
	defer w.Destroy()
	w.SetSize(1000, 1200, webview2.HintNone)
	w.Init(injectInitJS())
	w.Navigate(url)

	app = NewAPP(w)
	app.Bind("emitEvent", bom.EmitEvent)
	app.Bind("download", downloadFunc(app))
	app.Bind("testFunc", testFunc(app))
	app.Bind("setSaveFolder", setSaveFolder)
	app.Bind("pickFolder", pickFolder)
	app.Bind("delayFuncs", func() {
		app.RunJS(injectDelayJS())
	})

	once := sync.Once{}
	bom.AddEventListener(DOMContentLoaded, func(eventType EventType, event Event) {
		once.Do(func() {
			//app.RunJS("sendSaveFolder();")
		})
	})

	app.Run()
	app.Wait()
}

func getChromium(w webview2.WebView) *edge.Chromium {
	browser := reflect.ValueOf(w).Elem().FieldByName("browser")
	browser = reflect.NewAt(browser.Type(), unsafe.Pointer(browser.UnsafeAddr())).Elem()
	return browser.Interface().(*edge.Chromium)
}
