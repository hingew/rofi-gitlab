package data

import "fmt"

type Project struct {
	Name string `json:"path_with_namespace"`
	Id   int    `json:"id"`
}

type Issue struct {
	Title string `json:"title"`
	Id    int    `json:"iid"`
}

func (i *Issue) ShowTitle() string {
	return fmt.Sprintf("#%d %s", i.Id, i.Title)
}
