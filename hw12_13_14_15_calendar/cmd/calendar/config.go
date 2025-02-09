package main

//nolint:depguard
import "github.com/BurntSushi/toml"

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	return
}

type Config struct {
	Logger LoggerConfig
	DB     DBConfig
	HTTP   HTTPConfig
}

type LoggerConfig struct {
	Level string
}

type DBConfig struct {
	InMemory  bool
	DSN       string
	Migration string
}

type HTTPConfig struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{}
}
