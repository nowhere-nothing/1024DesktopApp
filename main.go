package main

import (
	"flag"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"webview_demo/config"
	"webview_demo/downloader"
	"webview_demo/store"
)

var configPath string
var genConfig string

const name = "picxus"

func main() {
	d, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	cp := filepath.Join(d, "Documents", name, "config.yaml")
	flag.StringVar(&configPath, "config", cp, "/path/to/config.yaml")
	flag.StringVar(&genConfig, "init", "", "generate config file template")
	flag.Parse()

	if len(genConfig) != 0 {
		if err := config.GenerateTemplate(genConfig); err != nil {
			panic(err)
		}
		return
	}

	c, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()

	l, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		logger.Errorln("parse log level", err, "use debug")
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(l)
	}

	if logger.IsLevelEnabled(logrus.DebugLevel) {
		logger.SetReportCaller(true)
	}

	if c.LogPath != "stdout" {
		logFile, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		logger.SetOutput(logFile)
	}

	var ds store.Storage
	if !strings.HasPrefix(c.Dsn, "file://") {
		err = store.InitDB(c.Dsn)
		if err != nil {
			logger.WithError(err).Errorln("init database dsn:", c.Dsn)
			dialog.Message("start up %v", err).Title("run error").Error()
			os.Exit(1)
		}
		ds = store.NewDBStorage()
	} else {
		ds = store.NewFsStorage(c.Dsn[7:])
	}

	logger.Debugln("use storage", c.Dsn)

	wg := sync.WaitGroup{}

	dl := downloader.NewDownloader(&wg)

	gPool, err := ants.NewPool(c.WorkPoolSize,
		ants.WithLogger(logger),
		ants.WithNonblocking(true),
	)
	if err != nil {
		logger.WithError(err).Errorln("init work pool size:", c.WorkPoolSize)
		dialog.Message("start up %v", err).Title("run error").Error()
		os.Exit(1)
	}
	defer gPool.Release()

	app := NewApp(logger, c)

	err = app.Run(gPool, dl, ds)
	if err != nil {
		logger.WithError(err).Errorln("run main app")
		dialog.Message("start up %v", err).Title("run error").Error()
		os.Exit(1)
	}
	app.Close()
	wg.Wait()
	if gPool.Running() != 0 {
		logger.Warningf("pool worker num %d running but progress exit", gPool.Running())
	}
}
