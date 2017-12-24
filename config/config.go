package config

import (
	"github.com/mdh67899/go-utils/file/content"
	"github.com/mdh67899/mail-provider/model"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func ParseConfig(filename string) *model.Config {
	if filename == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalln("config file", filename, ", failed:", err, ", try to mv cfg.example.yaml to cfg.yaml")
	}

	var cfg = model.Config{}

	configContent, err := content.ToTrimString(filename)
	if err != nil {
		log.Fatalln("read config file:", filename, "fail:", err)
	}

	err = yaml.Unmarshal([]byte(configContent), &cfg)
	if err != nil {
		log.Fatalln("parse config file:", filename, "fail:", err)
	}

	log.Println("read config file:", filename, "successfully")

	return &cfg
}
