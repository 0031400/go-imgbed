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
		Cors         struct {
			Origins     []string `yaml:"origins"`
			Methods     []string `yaml:"methods"`
			Headers     []string `yaml:"headers"`
			Credentials bool     `yaml:"credentials"`
		} `yaml:"cors"`
	} `yaml:"server"`
	WaterMark struct {
		Enable   bool   `yaml:"enabled"`
		Font     string `yaml:"font"`
		Text     string `yaml:"text"`
		Size     int    `yaml:"size"`
		Pdi      int    `yaml:"pdi"`
		Position struct {
			X int `yaml:"x"`
			Y int `yaml:"y"`
		} `yaml:"position"`
		Color struct {
			R int `yaml:"r"`
			G int `yaml:"g"`
			B int `yaml:"b"`
			A int `yaml:"a"`
		} `yaml:"color"`
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
