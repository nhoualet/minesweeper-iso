package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadConfig(filePath string, output interface{}) error {
	cfgFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("load config: couldn't read config file: %s", err)
	}
	err = yaml.UnmarshalStrict(cfgFile, output)
	if err != nil {
		return fmt.Errorf("load config: couldn't unmarshal config file: %s", err)
	}
	return nil
}

func SaveConfig(filePath string, data interface{}) error {
	d, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("save config: couldn't marshall data: %s", err.Error())
	}
	err = ioutil.WriteFile(filePath, d, 0644)
	if err != nil {
		return fmt.Errorf("save config: couldn't write data: %s", err.Error())
	}
	return nil
}
