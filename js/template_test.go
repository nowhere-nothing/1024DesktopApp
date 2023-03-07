package js

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
	"webview_demo/config"
)

func TestCompileTemplate(t *testing.T) {
	c, err := config.LoadConfig("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	str, err := CompileTemplate(logrus.New(), c, "https://t66y.com/index.php")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range strings.Split(str, "\n") {
		t.Log(v)
	}
}

func TestEnv(t *testing.T) {
	for _, v := range os.Environ() {
		t.Log(v)
	}
}
