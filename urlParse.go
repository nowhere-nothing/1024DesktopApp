package main

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
)

func FileName(url string) string {
	i := strings.LastIndex(url, "/")
	if i == -1 {
		t := make([]byte, 0, 16)
		for _, v := range md5.Sum([]byte(url)) {
			t = append(t, v)
		}
		return base64.StdEncoding.EncodeToString(t)
	}
	return url[i+1:]
}
