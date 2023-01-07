package main

import (
	"bytes"
	"embed"
	"path"
	"text/template"
)

//go:embed js/*
var jsDir embed.FS

//go:embed init.tmpl
var tmpl string
var jsTmpl = template.Must(template.New("").Parse(tmpl))

func combineJs() string {
	data := make(map[string]string)
	entries, err := jsDir.ReadDir("js")
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		b, err := jsDir.ReadFile(path.Join("js", v.Name()))
		if err != nil {
			panic(err)
		}
		data[v.Name()] = string(b)
	}
	// todo: template add func map

	b := bytes.NewBuffer(make([]byte, 0, 1024*5))
	if err := jsTmpl.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}
