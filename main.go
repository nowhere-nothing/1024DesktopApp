package main

import (
	"fmt"
	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"github.com/ncruces/zenity"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"sync"
	"unsafe"
	"webview_demo/bom"
	"webview_demo/config"
	"webview_demo/downloader"
	"webview_demo/js"
	"webview_demo/store"
)

type App struct {
	w webview2.WebView

	sync.WaitGroup

	logger *logrus.Logger
	Pool   *ants.Pool

	B  *bom.BOM
	P  *downloader.Progress
	C  *config.Config
	D  *downloader.Downloader
	S  store.Storage
	DH *store.DownloadHistory

	configSavePath string
}

func NewAPP(w webview2.WebView) *App {
	logger := logrus.New()
	a := &App{
		w:      w,
		logger: logger,
		B:      bom.NewBOM(),
		P:      &downloader.Progress{},
		S:      store.NewFsStorage(""),
	}
	return a
}

func (a *App) Init(confPath string) error {
	a.configSavePath = confPath
	if c, err := config.GetConfig(confPath); err != nil {
		return err
	} else {
		a.C = c
	}
	a.P.SetChangeHandler(func(x, y int32) {
		if x == y {
			a.P.Reset()
			x, y = 0, 0
		}
		a.logger.Printf("progress %d/%d", x, y+1)
		a.RunJS(fmt.Sprintf("setProgress(%d, %d)", x, y+1))
	})

	a.D = downloader.NewDownloader(&a.WaitGroup)

	if p, err := ants.NewPool(100,
		ants.WithLogger(nil),
		ants.WithNonblocking(true),
		ants.WithLogger(a.logger),
	); err != nil {
		return err
	} else {
		a.Pool = p
	}
	return nil
}

func (a *App) Close() {
	if err := config.SaveConfig(a.C, a.configSavePath); err != nil {
		a.logger.WithError(err).Errorln("save config")
	}
	a.Pool.Release()
	a.w.Dispatch(func() {
		a.w.Destroy()
	})
	a.Wait()
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

func (a *App) Bind(name string, f any) {
	if err := a.w.Bind(name, f); err != nil {
		panic(err)
	}
}

func (a *App) Run() {
	a.w.Run()
}

func selectSite() string {
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
	return url
}

func NewWebview(url string) webview2.WebView {
	w := webview2.New(true)
	w.SetSize(1000, 1200, webview2.HintNone)

	w.Init(js.InjectInitJS())

	w.Navigate(url)
	return w
}

func main() {
	app := NewAPP(NewWebview(selectSite()))
	if err := app.Init("./config.json"); err != nil {
		panic(err)
	}

	app.Bind("pickFolder", FolderPicker)
	app.Bind("emitEvent", app.B.EmitEvent)
	app.Bind("download", DownloadFunc(app))
	app.Bind("testFunc", testFunc(app))
	app.Bind("setSaveFolder", func() error {
		s, err := FolderPicker("选择文件存储目录", app.C.SavePath)
		if err != nil {
			return err
		}
		app.C.SavePath = s
		return nil
	})
	app.Bind("delayFuncs", func() {
		app.RunJS(js.InjectDelayJS())
	})
	app.Bind("injectDyn", func(fn string) error {
		data, err := js.DynLoadJS(fn)
		if err != nil {
			return err
		}
		app.RunJS(string(data))
		return nil
	})

	once := sync.Once{}
	app.B.AddEventListener(bom.DOMContentLoaded, func(eventType bom.EventType, event bom.Event) {
		once.Do(func() {
			//app.RunJS("sendSaveFolder();")
			app.RunJS("setProgress(0,0);")
		})
	})

	defer app.Close()
	app.Run()
}

func getChromium(w webview2.WebView) *edge.Chromium {
	browser := reflect.ValueOf(w).Elem().FieldByName("browser")
	browser = reflect.NewAt(browser.Type(), unsafe.Pointer(browser.UnsafeAddr())).Elem()
	return browser.Interface().(*edge.Chromium)
}
