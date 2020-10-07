/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package domain

type Issue struct {
	Fields Fields `json:"fields"`
	Key    string `json:"key"`
}

type Fields struct {
	Summary     string      `json:"summary"`
	Assignee    Assignee    `json:"assignee"`
	Labels      []string    `json:"labels"`
	Status      Status      `json:"status"`
	Description string      `json:"description"`
	Priority    Priority    `json:"priority"`
	Links       []IssueLink `json:"issuelinks"`
}

type Assignee struct {
	Self         string `json:"self"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type Status struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Labels struct {
	Labels string `json:"labels"`
}

type IssueLink struct {
	Id string `json:"id"`
}
type Priority struct {
	Self    string `json:"self"`
	IconURL string `json:"iconUrl"`
	Name    string `json:"name"`
	Id      string `json:"id"`
}
