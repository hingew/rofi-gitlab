package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"rofi-gitlab/utils"
)

const configFileName string = "config.json"

type Config struct {
	BaseUrl         string
	Token           string
	TTL             int
}

func New() *Config {
	return &Config{BaseUrl: "", Token: "", TTL: 600}
}

func Read() (error, *Config) {
	var config *Config
	path := utils.Path(configFileName)
	jsonFile, err := os.Open(path)

	if err != nil {
		config = New()
		err = config.Write()

		return err, config
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)

	return nil, config
}

func (c *Config) Write() error {
	file, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		return err
	}

	err = os.MkdirAll(utils.Dir(), 0755)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(utils.Path(configFileName), file, 0644)

	return err
}
