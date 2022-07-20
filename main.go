package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	BaseUrl string
	Token   string
	TTL     int
}

type Cache struct {
	Projects  []Project `json:"projects"`
	Timestamp time.Time `json:"timestamp"`
}

type Project struct {
	Name string `json:"path_with_namespace"`
}

func main() {
	err, config := readConfig()

	if err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]

	options := [2]string{"issues", "pipelines"}
	switch len(args) {
	case 0:
		projects := getGitlabProjects(config)
		for _, project := range projects {
			fmt.Println(project.Name)
		}
	case 1:
		for _, option := range options {
			fmt.Println(option)
		}
	case 2:
		fmt.Println(fmt.Sprintf("%s/%s/-/%s", config.BaseUrl, args[0], args[1]))
	}

}

func getGitlabProjects(config *Config) []Project {
	err, cache := readCache()

	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	ttl := cache.Timestamp.Add(time.Second * time.Duration(config.TTL))

	// check if cache is valid
	if ttl.Before(now) {
		var projects []Project

		err = get(config.BaseUrl+"/api/v4/projects?simple=true&per_page=100", config.Token, &projects)

		if err != nil {
			log.Fatal(err)
		}

		cache.Projects = projects
		cache.Timestamp = time.Now()

		writeCache(*cache)
	}

	return cache.Projects
}

func get(url string, token string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
func readConfig() (error, *Config) {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		return err, nil
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	return nil, &config
}

func writeConfig(config Config) error {
	file, err := json.MarshalIndent(config, "", " ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.json", file, 0644)

	return err
}

func readCache() (error, *Cache) {
	jsonFile, err := os.Open("cache.json")

	if err != nil {
		return err, nil
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	var cache Cache

	json.Unmarshal(byteValue, &cache)

	return nil, &cache
}

func writeCache(cache Cache) error {
	file, err := json.MarshalIndent(cache, "", " ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("cache.json", file, 0644)

	return err
}
