package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Port   string `yaml:"port"`
	DBRoot string `yaml:"DBRoot"`
}

var configFile = "../port_config.yml"

func (c *conf) getConf(file string) *conf {

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
