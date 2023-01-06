package main

import (
	"fmt"
	"github.com/jchv/go-webview2"
	"path"
	"strings"
	"sync/atomic"
)

type GlobalProgress struct {
	Max atomic.Uint32
	Val atomic.Uint32
	w   webview2.WebView
}

var gp GlobalProgress
var saveFolder = `D:\temp`

var runJs func(js string)
var alert func(str string)

func main() {
	w := webview2.New(true)
	defer w.Destroy()
	gp.w = w
	w.SetSize(1000, 1200, webview2.HintNone)
	w.Init(combineJs())
	w.Navigate("https://t66y.com/index.php")
	runJs = func(js string) {
		w.Dispatch(func() {
			w.Eval(js)
		})
	}
	alert = func(str string) {
		w.Dispatch(func() {
			w.Eval(fmt.Sprintf("alert(%s);", str))
		})
	}
	if err := w.Bind("download", download); err != nil {
		panic(err)
	}
	if err := w.Bind("testFunc", func() {
		runJs(`console.log("Hello");`)
		runJs(`alert("Hello");`)
	}); err != nil {
		panic(err)
	}
	w.Run()
}

func download(title, url string, images []string) {
	title = strings.TrimSpace(title)
	fn, err := DownloadFunc(path.Join(saveFolder, title))
	if err != nil {
		alert(fmt.Sprintf("%s", err))
		return
	}
	go fn(url, images)
}

func AddMax(n uint32) {
	gp.Max.Add(n)
	SetPb()
}
func AddVal(n uint32) {
	gp.Val.Add(n)
	SetPb()
}

func SetPb() {
	m := gp.Max.Load()
	v := gp.Val.Load()
	if m == v {
		gp.Max.Store(0)
		gp.Val.Store(0)
		m = 0
		v = 0
	}
	if runJs != nil {
		runJs(fmt.Sprintf("setProgress(%d, %d);", m, v))
	}
}
