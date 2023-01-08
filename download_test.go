package main

import "testing"

var target = `https://26img.com/i/vexmctsj.jpg`
var batchTarget = []string{
	"https://23img.com/i/2023/01/05/dn8iox.jpg",
	"https://23img.com/i/2023/01/05/do0wog.jpg",
	"https://23img.com/i/2023/01/05/dooofi.jpg",
	"https://23img.com/i/2023/01/05/dndgxd.jpg",
	"https://23img.com/i/2023/01/05/dne76g.jpg",
	"https://23img.com/i/2023/01/05/doemfq.jpg",
	"https://23img.com/i/2023/01/05/dn8n92.jpg",
	"https://23img.com/i/2023/01/05/dncpos.jpg",
	"https://23img.com/i/2023/01/05/dn9i4x.jpg",
	"https://23img.com/i/2023/01/05/dn9ook.jpg",
	"https://23img.com/i/2023/01/05/dnak60.jpg",
	"https://23img.com/i/2023/01/05/dnb5sw.jpg",
	"https://23img.com/i/2023/01/05/dnbtvd.jpg",
	"https://23img.com/i/2023/01/05/dnc8ag.jpg",
	"https://23img.com/i/2023/01/05/dnchvw.jpg",
	"https://23img.com/i/2023/01/05/do2sga.jpg",
	"https://23img.com/i/2023/01/05/dncvck.jpg",
	"https://23img.com/i/2023/01/05/doe0x7.jpg",
	"https://23img.com/i/2023/01/05/dn90lu.jpg",
	"https://23img.com/i/2023/01/05/doh90o.jpg",
	"https://23img.com/i/2023/01/05/dom3cc.jpg",
	"https://23img.com/i/2023/01/05/do3v7t.jpg",
	"https://23img.com/i/2023/01/05/donam1.jpg",
	"https://23img.com/i/2023/01/05/do4uzn.jpg",
	"https://23img.com/i/2023/01/05/dofgzb.jpg",
	"https://23img.com/i/2023/01/05/dofc7o.jpg",
	"https://23img.com/i/2023/01/05/dodtkr.jpg",
	"https://23img.com/i/2023/01/05/do38cx.jpg",
	"https://23img.com/i/2023/01/05/do4bng.jpg",
	"https://23img.com/i/2023/01/05/dooa21.jpg",
	"https://23img.com/i/2023/01/05/dont7h.jpg",
	"https://23img.com/i/2023/01/05/domiwn.jpg",
	"https://23img.com/i/2023/01/05/dogtaz.jpg",
	"https://23img.com/i/2023/01/05/do1lu4.jpg",
	"https://23img.com/i/2023/01/05/doeepy.jpg",
	"https://23img.com/i/2023/01/05/dnennh.jpg",
}

func TestDownload(t *testing.T) {
	data, err := FetchData(target)
	if err != nil {
		t.Fatal(err)
	}
	err = SaveBinary(".", "a.jpg", data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveFromUrl(t *testing.T) {
	const gzipPic = "https://23img.com/i/2023/01/05/dn8iox.jpg"
	err := SaveFromUrl(gzipPic, ".", "a.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFunc(t *testing.T) {
	fn, err := Downloader("./aaa")
	if err != nil {
		t.Fatal(err)
	}
	fn("https://t66y.com/htm_mob/2301/8/5475955.html", batchTarget)
}

func TestCombineJs(t *testing.T) {
	str := combineJs()
	t.Logf(str)
}
