package store

import (
	"crypto/sha1"
	"encoding/hex"
)

type PostContent struct {
	Title  string
	Url    string
	Images []string
	hash   string
}

func (pc *PostContent) HashHex() string {
	if len(pc.hash) == 0 {
		pc.hash = HashHex([]byte(pc.Title + pc.Url))
	}
	return pc.hash
}

type PostImage struct {
	Name        string
	Url         string
	Data        []byte
	ContentType string
	hash        string
}

func (pi *PostImage) HashHex() string {
	if len(pi.hash) == 0 {
		pi.hash = HashHex(pi.Data)
	}
	return pi.hash
}

type Post struct {
	Title string
	Url   string
}

func HashHex(data []byte) string {
	s := sha1.Sum(data)
	r := hex.EncodeToString(s[:])
	return r
}
