package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	SavePath string `json:"save_path"`
}

func GetConfig(path string) (*Config, error) {
	ff, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(ff)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if len(data) == 0 {
		return conf, nil
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func SaveConfig(c *Config, path string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
