package model

type Config struct {
	Addr string   `yaml:"addr"`
	Auth AuthUser `yaml:"auth"`
}
