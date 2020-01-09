package main

import (
	"errors"
	"io/ioutil"

	"github.com/popu125/goImgClip/log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HotKey  string         `yaml:"hot_key"`
	LogFile string         `yaml:"log_file"`
	Targets []targetConfig `yaml:"targets"`
}

type targetConfig struct {
	Name         string    `yaml:"name"`
	Target       string    `yaml:"target"`
	TargetConfig yaml.Node `yaml:"config"`
}

func LoadConfig(path string) *Config {
	var conf Config
	l := log.GetLogger("config")
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	l.Info().Str("file", path).Msg("Loading config")

	data, err := ioutil.ReadFile(path)
	check(err)
	check(yaml.Unmarshal(data, &conf))

	// config pre-check
	if len(conf.Targets) == 0 {
		check(errors.New("there should be at least one target"))
	}

	return &conf
}
