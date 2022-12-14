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
	Projects        []Project            `json:"projects"`
	Timestamp       time.Time            `json:"timestamp"`
	Issues          map[string]IssueList `json:"issues"`
	SelectedProject string               `json:"selected_project"`
}

type IssueList struct {
	List      []Issue   `json:"issues"`
	Timestamp time.Time `json:"timestamp"`
}

func (c *Cache) ProjectID(path string) int {
	for _, project := range c.Projects {
		if project.Name == path {
			return project.Id
		}
	}

	return 0
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
