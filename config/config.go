package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Dsn              string      `yaml:"dsn"`
	LogPath          string      `yaml:"log_path"`
	LogLevel         string      `yaml:"log_level"`
	WorkPoolSize     int         `yaml:"work_pool_size"`
	ProgressFunction string      `yaml:"progress_function"`
	Template         string      `yaml:"template"`
	ScriptDir        string      `yaml:"script_dir"`
	SiteMap          []*SiteItem `yaml:"site_map"`
	ViewPort         `yaml:"view_port"`
}

type ViewPort struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type SiteItem struct {
	Url        string `yaml:"url"`
	ScriptName string `yaml:"script_name"`
}

func LoadConfig(p string) (*Config, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	fmt.Println("loaded config from", p)
	fmt.Print(string(data))
	return c, nil
}

func GenerateTemplate(p string) error {
	data, err := yaml.Marshal(&Config{
		Dsn:              "[root:root@tcp(localhost:3306)/db | file://save/to/path]",
		LogPath:          "['./log' | 'stdout']",
		LogLevel:         "debug",
		Template:         "./js/assets/init.tmpl",
		ScriptDir:        "./js/assets",
		ProgressFunction: "setProgress",
		SiteMap: []*SiteItem{
			{Url: "host_name1", ScriptName: `"init script file1"`},
			{Url: "host_name2", ScriptName: `"init script file2"`},
		},
		ViewPort: ViewPort{
			Height: 600,
			Width:  800,
		},
	})
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}
