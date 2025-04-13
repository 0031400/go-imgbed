package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		ThumbnailDir string `yaml:"thumbnailDir"`
		OriginalDir  string `yaml:"originalDir"`
		PublicDir    string `yaml:"publicDir"`
		RandomNum    int    `yaml:"randomNum"`
		Token        string `yaml:"token"`
	} `yaml:"server"`
	WaterMark struct {
		font     string `yaml:"font"`
		text     string `yaml:"text"`
		size     int    `yaml:"size"`
		Position struct {
			x int `yaml:"x"`
			y int `yaml:"y"`
		}
	} `yaml:"watermark"`
	Thumbnail struct {
		Width  int `yaml:"width"`
		Height int `yaml:"height"`
	} `yaml:"thumbnail"`
	Quality int    `yaml:"quality"`
	LogFile string `yaml:"logFile"`
}

func (c *Config) Load() error {
	configFile := flag.String("c", "config.yaml", "the config file")
	flag.Parse()
	configBytes, err := os.ReadFile(*configFile)
	if err != nil {
		return errors.New("read config file fail\n" + err.Error())
	}
	err = yaml.Unmarshal(configBytes, &c)
	if err != nil {
		return errors.New("unmarshal config file fail\n" + err.Error())
	}
	return nil
}
