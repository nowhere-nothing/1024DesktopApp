package js

import (
	"bytes"
	"embed"
	"path"
	"text/template"
)

//go:embed assets/init
var initJSDir embed.FS

//go:embed assets/delay
var delayJSDir embed.FS

//go:embed assets/dyn
var dynJSDir embed.FS

//go:embed assets/init.tmpl
var tmpl string
var jsTmpl = template.Must(template.New("").Parse(tmpl))

//go:embed assets/delay.tmpl
var delayJSTmpl string
var delayTmpl = template.Must(template.New("").Parse(delayJSTmpl))

const basePath = "assets"

func InjectDelayJS() string {
	data := mappingFiles(delayJSDir, basePath+"/delay")
	b := bytes.NewBuffer(make([]byte, 0, 1024*5))
	if err := delayTmpl.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}

func InjectInitJS() string {
	// todo: template add func map
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

func DynLoadJS(fn string) ([]byte, error) {
	return dynJSDir.ReadFile(basePath + "/" + fn)
}
