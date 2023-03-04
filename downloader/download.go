package downloader

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"io"
	"net/http"
	"strings"
	"sync"
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
const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

type Downloader struct {
	rawCli *http.Client
	cli    *resty.Client
	wg     *sync.WaitGroup
}

func NewDownloader(wg *sync.WaitGroup) *Downloader {
	client := resty.New()
	//client.SetProxy()
	//client.SetHeader("user-agent", ua)
	client.SetHeaders(headers)
	client.SetCloseConnection(true)
	client.SetTimeout(2 * time.Minute)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	return &Downloader{
		rawCli: &http.Client{},
		cli:    client,
		wg:     wg,
	}
}

var noDataError = errors.New("response no data")

func (d *Downloader) Fetch(url string) (*FetchData, error) {
	d.wg.Add(1)
	defer d.wg.Done()
	rsp, err := d.cli.R().Get(url)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http status %s", http.StatusText(rsp.StatusCode()))
	}
	data := rsp.Body()
	if len(data) == 0 {
		return nil, noDataError
	}

	fd := &FetchData{
		Data:        data,
		Url:         url,
		ContentType: rsp.Header().Get(contentType),
	}
	return fd, nil
}

func (d *Downloader) RawFetch(url string) (*FetchData, error) {
	d.wg.Add(1)
	defer d.wg.Done()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := d.rawCli.Do(req)
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
	return &FetchData{
		Data:        data,
		Url:         url,
		ContentType: resp.Header.Get(contentType),
	}, nil
}

type FetchData struct {
	Data        []byte
	Url         string
	ContentType string
}

func (fd *FetchData) Name() string {
	fn, ext := fd.split(fd.Url)
	if ext == "" {
		if t, err := filetype.Match(fd.Data); err == nil {
			ext = t.Extension
		} else if e := fd.oneOfExt(fd.ContentType); e != "" {
			ext = e
		} else {
			ext = ".unknown"
		}
	}
	return fn + ext
}

func (*FetchData) oneOfExt(mime string) string {
	var ext string
	filetype.Types.Range(func(key, value any) bool {
		kind := value.(types.Type)
		if kind.MIME.Value == mime {
			ext = kind.Extension
			return false
		}
		return false
	})
	return ext
}

func (*FetchData) split(url string) (string, string) {
	var fn, ext string
	ei := strings.LastIndex(url, ".")
	if ei != -1 {
		ext = url[ei:] // with dot
	}

	fi := strings.LastIndex(url, "/")
	if fi != -1 {
		if ei == -1 {
			fn = url[fi+1:]
		} else {
			fn = url[fi+1 : ei]
		}
	}
	return fn, ext
}
