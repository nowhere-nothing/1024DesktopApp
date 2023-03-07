package js

import (
	"bytes"
	_ "embed"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"webview_demo/config"
)

/*
//go:embed assets/init
var initJSDir embed.FS

func InjectInitJS() string {
	const basePath = "assets"
	var jsTmpl *template.Template

	data := mappingFiles(initJSDir, basePath+"/init")
	b := bytes.NewBuffer(make([]byte, 0, 1024*5))
	if err := jsTmpl.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}

func mappingFiles(folder embed.FS, dir string) map[string]string {
	m := make(map[string]string)
	entries, err := folder.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		b, err := folder.ReadFile(path.Join(dir, v.Name()))
		if err != nil {
			panic(err)
		}
		m[v.Name()] = string(b)
	}
	return m
}
*/

//go:embed jquery.min.js
var jq string

func GetJQ() string {
	return jq
}

func CompileTemplate(logger *logrus.Logger, c *config.Config, url string) (string, error) {
	td, err := os.ReadFile(c.Template)
	if err != nil {
		return "", err
	}

	var sc []byte
	for _, v := range c.SiteMap {
		if v.Url == url {
			sc, err = os.ReadFile(filepath.Join(c.ScriptDir, v.ScriptName))
			if err != nil {
				return "", err
			}
			break
		}
	}
	if len(sc) == 0 {
		logger.Warningln("site", url, "init script not found or empty")
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"loadFile": func(f string) (string, error) {
			d, err := os.ReadFile(filepath.Join(c.ScriptDir, f))
			if err != nil {
				return "", err
			}
			return string(d), nil
		},
	}).Parse(string(td))
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(make([]byte, 0, 1024*5))
	if err = tmpl.Execute(buf, map[string]any{
		"siteScript": string(sc),
	}); err != nil {
		return "", err
	}

	logger.Debugln("browser init script")
	for _, v := range strings.Split(buf.String(), "\n") {
		logger.Debug(v)
	}

	return buf.String(), nil
}
