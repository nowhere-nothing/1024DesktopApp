package main

import (
	"bytes"
	"embed"
	"path"
	"text/template"
)

//go:embed js/*
var jsDir embed.FS
var initFiles = []string{
	"viewer.min.css",
	"viewer.min.js",
	"init.js",
}

func combineJs() string {
	data := make(map[string]string)
	for _, v := range initFiles {
		b, err := jsDir.ReadFile(path.Join("js", v))
		if err != nil {
			panic(err)
		}
		data[v] = string(b)
	}
	// todo: template add func map

	b := bytes.NewBuffer(make([]byte, 0, 1024*5))
	if err := jsTmpl.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}

//go:embed init.tmpl
var tmpl string
var jsTmpl = template.Must(template.New("").Parse(tmpl))
