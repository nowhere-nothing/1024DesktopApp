package downloader

import "sync/atomic"

type Progress struct {
	max      atomic.Int32
	cur      atomic.Int32
	onChange func(m, c int32)
}

func (d *Progress) SetChangeHandler(f func(x, y int32)) {
	d.onChange = f
}

func (d *Progress) Max() int32 {
	return d.max.Load()
}
func (d *Progress) Cur() int32 {
	return d.cur.Load()
}

func (d *Progress) AddMax(n int32) {
	d.max.Add(n)
	if d.onChange != nil {
		d.onChange(d.Max(), d.Cur())
	}
}
func (d *Progress) AddCur(n int32) {
	d.cur.Add(n)
	if d.onChange != nil {
		d.onChange(d.Max(), d.Cur())
	}
}
func (d *Progress) Reset() {
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
