package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"rofi-gitlab/config"
	"rofi-gitlab/data"
	"time"
)

const (
	Issues        string = "Issues"
	Pipelines            = "Pipelines"
	Boards               = "Boards"
	Project              = "Project"
	Environments         = "Environments"
	MergeRequests        = "Merge Requests"
)

func main() {
	args := os.Args[1:]

	err, config := config.Read()
	options := []string{Issues, Pipelines, Boards, Project, Environments, MergeRequests}

	if err != nil {
		log.Fatal(err)
	}

	if config.Choosen == "" && len(args) == 0 {
		projects := getGitlabProjects(config)
		for _, project := range projects {
			fmt.Println(project.Name)
		}
		return
	}

	if config.Choosen == "" && len(args) > 0 {
		config.Choosen = args[0]
		if err := config.Write(); err != nil {
			log.Fatal(err)
		}

		for _, option := range options {
			fmt.Println(option)
		}
		return
	}

	exec.Command("xdg-open", path.Join(config.BaseUrl, getPath(config.Choosen, args[0]))).Start()

	config.Choosen = ""
	config.Write()
}

func getPath(project string, option string) string {
	switch option {
	case Issues:
		return path.Join(project, "-", "issues")
	case Pipelines:
		return path.Join(project, "-", "pipelines")
	case Boards:
		return path.Join(project, "-", "boards")
	case Environments:
		return path.Join(project, "-", "environments")
	case MergeRequests:
		return path.Join(project, "-", "merge_requests")
	case Project:
		return project
	default:
		return project
	}

}

func getGitlabProjects(config *config.Config) []data.Project {

	err, cache := data.ReadCache()

	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	ttl := cache.Timestamp.Add(time.Second * time.Duration(config.TTL))

	// check if cache is valid
	if ttl.Before(now) {
		var projects []data.Project

		err = get(config.BaseUrl+"/api/v4/projects?simple=true&per_page=100", config.Token, &projects)

		if err != nil {
			log.Fatal(err)
		}

		cache.Projects = projects
		cache.Timestamp = time.Now()
		cache.Write()
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
