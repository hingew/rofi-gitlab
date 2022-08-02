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
	"strings"
	"time"
)

const (
	Issues        string = "Issues"
	Pipelines            = "Pipelines"
	Boards               = "Boards"
	Project              = "Project"
	Environments         = "Environments"
	MergeRequests        = "Merge Requests"
	Logs                 = "Logs"
)

func main() {
	args := os.Args[1:]
	options := []string{Project, Issues, Pipelines, Boards, Logs, Environments, MergeRequests}

	err, config := config.Read()

	if err != nil {
		log.Fatal(err)
	}
	err, cache := data.ReadCache()

	if err != nil {
		log.Fatal(err)
	}
	if len(args) == 0 {
		projects := getGitlabProjects(config)
		for _, project := range projects {
			fmt.Println(project.Name)
		}
		return
	} else if config.SelectedProject == "" && len(args) > 0 {
		config.SelectedProject = args[0]
		config.Write()

		if err := config.Write(); err != nil {
			log.Fatal(err)
		}

		for _, option := range options {
			fmt.Println(option)
		}
		return
	} else if config.SelectedProject != "" && config.ShowIssue && len(args) > 0 {
		id := strings.Replace(strings.Fields(args[0])[0], "#", "", -1)
		open(config.BaseUrl, path.Join(config.SelectedProject, "-", "issues", id))

		config.SelectedProject = ""
		config.ShowIssue = false
		config.Write()
	} else {
		performAction(config, cache, args[0])
	}

}

func open(baseUrl string, subPath string) {
	exec.Command("xdg-open", path.Join(baseUrl, subPath)).Start()
}

func performAction(config *config.Config, cache *data.Cache, option string) {
	project := config.SelectedProject
	switch option {
	case Issues:
		issues := getProjectIssues(config, cache.ProjectID(project))
		config.ShowIssue = true
		config.Write()

		for _, issue := range issues {
			fmt.Println(issue.ShowTitle())
		}
	case Pipelines:
		open(config.BaseUrl, path.Join(project, "-", "pipelines"))
	case Boards:
		open(config.BaseUrl, path.Join(project, "-", "boards"))
	case Environments:
		open(config.BaseUrl, path.Join(project, "-", "environments"))
	case MergeRequests:
		open(config.BaseUrl, path.Join(project, "-", "merge_requests"))
	case Logs:
		open(config.BaseUrl, path.Join(project, "-", "logs"))
	case Project:
		open(config.BaseUrl, project)
	default:
		open(config.BaseUrl, project)
	}

}

func getProjectIssues(config *config.Config, projectID int) []data.Issue {
	err, cache := data.ReadCache()

	if err != nil {
		log.Fatal(err)
	}

	issues, exists := cache.Issues[config.SelectedProject]

	if !exists {
		err = get(config.BaseUrl+fmt.Sprintf("/api/v4/projects/%d/issues?state=opened", projectID), config.Token, &issues)
		if err != nil {
			log.Fatal(err)
		}

		cache.Issues[config.SelectedProject] = issues
		cache.Write()
	}

	return issues

}

func getGitlabProjects(config *config.Config) []data.Project {

	err, cache := data.ReadCache()

	if err != nil {
		log.Fatal(err)
	}

	// now := time.Now()
	// ttl := cache.Timestamp.Add(time.Second * time.Duration(config.TTL))

	// check if cache is valid
	// if ttl.Before(now) {
	if len(cache.Projects) == 0 {
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

func toRofiString(label string, value string) string {
	return value + `\0info\x1` + value
}
