package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mdh67899/mail-provider/model"
	"gopkg.in/yaml.v2"
)

type Config struct {
	sync.RWMutex
	Addr string         `yaml:"addr"`
	Auth model.AuthUser `yaml:"auth"`
}

var (
	Cfg = &Config{}
)

func ToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToTrimString(filePath string) (string, error) {
	str, err := ToString(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

func ParseConfig(filename string) {
	if filename == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalln("config file", filename, ", failed:", err, ", try to mv cfg.example.yaml to cfg.yaml")
	}

	var cfg *Config

	configContent, err := ToTrimString(filename)
	if err != nil {
		log.Fatalln("read config file:", filename, "fail:", err)
	}

	err = yaml.Unmarshal([]byte(configContent), &cfg)
	if err != nil {
		log.Fatalln("parse config file:", filename, "fail:", err)
	}

	Cfg.Lock()
	defer Cfg.Unlock()

	Cfg.Addr = cfg.Addr
	Cfg.Auth = cfg.Auth

	log.Println("read config file:", filename, "successfully")
}
