package main

import (
	"errors"
	"fmt"
	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"github.com/ncruces/zenity"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"os"
	"reflect"
	"strings"
	"sync"
	"unsafe"
	"webview_demo/bom"
	"webview_demo/config"
	"webview_demo/downloader"
	"webview_demo/js"
	"webview_demo/store"
)

type App struct {
	win      webview2.WebView
	logger   *logrus.Logger
	Bom      *bom.BOM
	Progress *downloader.Progress
	Config   *config.Config
}

func NewApp(l *logrus.Logger, c *config.Config) *App {
	return &App{
		logger:   l,
		Progress: &downloader.Progress{},
		Config:   c,
		Bom:      bom.NewBOM(),
	}
}

func (a *App) RunJS(str string) {
	a.win.Dispatch(func() {
		a.win.Eval(str)
	})
}

func (a *App) Bind(name string, f any) {
	ft := reflect.TypeOf(f)
	sb := strings.Builder{}
	sb.WriteString("in:[")
	for i := 0; i < ft.NumIn(); i++ {
		it := ft.In(i)
		sb.WriteString(it.Name())
		sb.WriteString(" ")
	}
	sb.WriteString("] out:[")
	for i := 0; i < ft.NumOut(); i++ {
		ot := ft.Out(i)
		sb.WriteString(ot.Name())
		sb.WriteString(" ")
	}
	sb.WriteString("]")

	a.logger.Debugln("bind to js function", name, sb.String())

	if err := a.win.Bind(name, f); err != nil {
		a.logger.WithError(err).Errorln("bind function", name)
		dialog.Message("start up %v", err).Title("run error").Error()
		os.Exit(1)
	}
}

func (a *App) bindFunctions(f *Fetcher) {
	a.Bind("pickFolder", FolderPicker)
	a.Bind("emitEvent", a.Bom.EmitEvent)
	a.Bind("download", f.Fetch)
	a.Bind("testFunction", testFunc(a))
	a.Bind("loadJQ", func() {
		a.RunJS(js.GetJQ())
	})
	a.Bind("loadJs", func(fn string) error {
		a.logger.Infoln("browser load file", fn)
		data, err := os.ReadFile(fn)
		if err != nil {
			return err
		}
		a.RunJS(string(data))
		return nil
	})
	once := sync.Once{}
	a.Bom.AddEventListener(bom.DOMContentLoaded, func(et bom.EventType, e bom.Event) {
		once.Do(func() {
			a.RunJS(fmt.Sprintf(a.Config.ProgressFunction, 0, 0))
		})
	})
}

func (a *App) Run(pool *ants.Pool, dl *downloader.Downloader, ds store.Storage) error {
	a.Progress.SetChangeHandler(func(x, y int32) {
		if x == y {
			a.Progress.Reset()
			x, y = 0, -1
		}
		a.RunJS(fmt.Sprintf("%s(%d,%d);", a.Config.ProgressFunction, x, y+1))
	})

	urls := make([]string, 0, len(a.Config.SiteMap))
	for _, v := range a.Config.SiteMap {
		urls = append(urls, v.Url)
	}
	site, err := zenity.List("选择网址", urls, zenity.Title(name))
	if err != nil {
		if errors.Is(zenity.ErrCanceled, err) {
			os.Exit(0)
		}
		return err
	}
	if site == "" {
		return fmt.Errorf("site not select")
	}

	j, err := js.CompileTemplate(a.logger, a.Config, site)
	if err != nil {
		return err
	}
	a.win = webview2.New(true)
	a.win.SetSize(a.Config.Width, a.Config.Height, webview2.HintNone)
	a.win.Init(j)
	a.win.Navigate(site)

	a.bindFunctions(&Fetcher{
		fetchHandler: dl,
		progress:     a.Progress,
		config:       a.Config,
		Logger:       a.logger,
		pool:         pool,
		Storage:      ds,
	})

	a.win.Run()
	return nil
}

func (a *App) Close() {
	a.win.Destroy()
}

func getChromium(w webview2.WebView) *edge.Chromium {
	browser := reflect.ValueOf(w).Elem().FieldByName("browser")
	browser = reflect.NewAt(browser.Type(), unsafe.Pointer(browser.UnsafeAddr())).Elem()
	return browser.Interface().(*edge.Chromium)
}
