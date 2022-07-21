package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"rofi-gitlab/utils"
	"time"
)

const cacheFileName string = "cache.json"

type Cache struct {
	Projects  []Project `json:"projects"`
	Timestamp time.Time `json:"timestamp"`
}

func newCache() *Cache {
	return &Cache{}
}

func ReadCache() (error, *Cache) {
	jsonFile, err := os.Open(utils.Path(cacheFileName))

	if err != nil {
		cache := &Cache{}
		err = cache.Write()
		return err, cache
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	var cache Cache

	json.Unmarshal(byteValue, &cache)

	return nil, &cache
}

func (c *Cache) Write() error {
	file, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(utils.Path(cacheFileName), file, 0644)

	return err
}
