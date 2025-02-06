package main

import "github.com/BurntSushi/toml"

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	return
}

type Config struct {
	Logger	LoggerConfig
	DB		DbConfig
	HTTP	HttpConfig
}

type LoggerConfig struct {
	Level string
}

type DbConfig struct {
	InMemory bool
	DSN       string
	Migration string
}

type HttpConfig struct {
	Host string
	Prot string
}

func NewConfig() Config {
	return Config{}
}
