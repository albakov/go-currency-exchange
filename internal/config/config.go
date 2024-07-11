package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	PathToDB string `toml:"abs_path_to_database"`
	CORS
}

type CORS struct {
	AccessControlAllowOrigin  string `toml:"access_control_allow_origin"`
	AccessControlAllowHeaders string `toml:"access_control_allow_headers"`
	AccessControlAllowMethods string `toml:"access_control_allow_methods"`
}

func MustNew() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c := &Config{}

	_, err = toml.DecodeFile(path.Join(wd, "config", "app.toml"), c)
	if err != nil {
		panic(err)
	}

	return c
}
