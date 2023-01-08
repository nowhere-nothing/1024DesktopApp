package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

var headers = map[string]string{
	"accept":             "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
	"accept-encoding":    "gzip, deflate, br",
	"accept-language":    "zh,en-US;q=0.9,en;q=0.8,zh-CN;q=0.7,zh-TW;q=0.6",
	"cache-control":      "no-cache",
	"pragma":             "no-cache",
	"sec-ch-ua":          `"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`,
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": `"Windows"`,
	"sec-fetch-dest":     "image",
	"sec-fetch-mode":     "no-cors",
	"sec-fetch-site":     "cross-site",
	"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
}

const contentType = "content-type"
const contentLength = "content-length"

func FetchData(url string) ([]byte, error) {
	cli := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch data status: %s", resp.Status)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SaveBinary(dir, fn string, data []byte) error {
	fullPath := path.Join(dir, fn)
	return os.WriteFile(fullPath, data, 0644)
}

func SaveFromUrl(url, dir, fn string) error {
	cli := resty.New()
	//cli.SetProxy("")
	cli.SetOutputDirectory(dir)
	cli.SetRetryCount(3)
	cli.SetRetryWaitTime(5 * time.Second)

	_, err := cli.R().
		SetHeaders(headers).
		SetOutput(fn).
		Get(url)
	return err
}

var InvalidPathName = errors.New("invalid path name")

func CheckDirEmpty(dir string) (bool, error) {
	if dir == "" {
		return false, InvalidPathName
	}
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return false, err
	}
	n, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(n) == 0, nil
}

var DirNotEmpty = errors.New("dir is not empty")

func NewCli(dir string) *resty.Client {
	cli := resty.New()
	cli.SetOutputDirectory(dir)
	cli.SetRetryCount(3)
	cli.SetRetryWaitTime(5 * time.Second)
	return cli
}

// todo:refactory to support other storage
func Downloader(dir string) (func(origin string, images []string), error) {
	if ok, err := CheckDirEmpty(dir); err != nil {
		return nil, err
	} else if !ok {
		return nil, DirNotEmpty
	}

	log, err := NewLog(path.Join(dir, "log.txt"))
	if err != nil {
		return nil, err
	}

	return func(origin string, images []string) {
		app.Add(1)
		wg := sync.WaitGroup{}
		start := time.Now()
		defer func() {
			log.Printf("download use time [%s]", time.Since(start).String())
			log.Close()
			app.Done()
		}()
		log.Printf("start download from [%s]", origin)
		log.Printf("total images [%d]", len(images))

		app.AddMax(int32(len(images)))

		for _, i := range images {
			r := NewCli(dir).R()
			img := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer app.AddCur(1)
				log.Printf("start download image [%s]", img)
				fn := FileName(img)
				resp, err := r.SetHeaders(headers).SetOutput(fn).Get(img)
				if err != nil {
					log.Errorf("download image [%s] [%v]", img, err)
				} else {
					ct := resp.Header().Get(contentType)
					cl := resp.Header().Get(contentLength)

					log.Printf("download image [%s] done size [%s] fileType [%s]", img, cl, ct)
				}
			}()
		}
		wg.Wait()
	}, nil
}

type DownloadState struct {
	max      atomic.Int32
	cur      atomic.Int32
	onChange func(m, c int32)
}

func (d *DownloadState) Max() int32 {
	return d.max.Load()
}
func (d *DownloadState) Cur() int32 {
	return d.cur.Load()
}

func (d *DownloadState) AddMax(n int32) {
	d.max.Add(n)
	if d.onChange != nil {
		d.onChange(d.Max(), d.Cur())
	}
}
func (d *DownloadState) AddCur(n int32) {
	d.cur.Add(n)
	if d.onChange != nil {
		d.onChange(d.Max(), d.Cur())
	}
}
func (d *DownloadState) Reset() {
	d.max.Store(0)
	d.cur.Store(0)
}

type PicMeta struct {
	OriginTitle string
	OriginURL   string
	SrcUrl      string
	Data        []byte
}

func (p *PicMeta) Name() string {
	return ""
}

type Collection struct {
	OriginTitle string   `json:"originTitle"`
	OriginURL   string   `json:"originURL"`
	ImageURL    []string `json:"imageURL"`
}
