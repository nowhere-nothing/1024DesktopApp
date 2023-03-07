package main

import (
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
	"webview_demo/config"
	"webview_demo/js"
)

func TestCompileTemplate(t *testing.T) {
	c, err := config.LoadConfig("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	str, err := js.CompileTemplate(logrus.New(), c, "https://t66y.com/index.php")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range strings.Split(str, "\n") {
		t.Log(v)
	}
}
