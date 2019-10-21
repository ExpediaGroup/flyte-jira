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

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HotelsDotCom/flyte-jira/domain"
	"net/http"
	"strings"
)

// Must be initialised before using
var JiraConfig Config

type Config struct {
	Host     string
	Username string
	Password string
}

type Comment struct {
	Body string `json:"body"`
}

type IssueRequest struct {
	Fields RequestFields `json:"fields"`
}

type RequestFields struct {
	Project   ProjectRequest   `json:"project"`
	Summary   string           `json:"summary"`
	IssueType IssueTypeRequest `json:"issuetype"`
}

type ProjectRequest struct {
	Key string `json:"key"`
}

type IssueTypeRequest struct {
	Name string `json:"name"`
}

func CommentIssue(issueId, comment string) (domain.Issue, error) {
	var issue domain.Issue
	b, err := json.Marshal(Comment{comment})
	if err != nil {
		return issue, err
	}

	path := fmt.Sprintf("/rest/api/2/issue/%s/comment", issueId)
	request, err := constructPostRequest(path, string(b))
	if err != nil {
		return issue, err
	}

	statusCode, err := SendRequest(request, &issue)
	if statusCode != http.StatusCreated {
		err = fmt.Errorf("issueId=%s : statusCode=%d", issueId, statusCode)
		return domain.Issue{}, err
	}
	if err != nil {
		err = fmt.Errorf("issueId=%s : err=%v", issueId, err)
		return domain.Issue{}, err
	}

	return issue, nil
}

func GetIssueInfo(issueId string) (domain.Issue, error) {
	var issue domain.Issue
	path := fmt.Sprintf("/rest/api/2/issue/%s", issueId)

	request, err := constructGetRequest(path)
	if err != nil {
		return issue, err
	}

	statusCode, err := SendRequest(request, &issue)
	if statusCode != http.StatusOK {
		err = fmt.Errorf("issueId=%s : statusCode=%d", issueId, statusCode)
		return domain.Issue{}, err
	}
	if err != nil {
		err = fmt.Errorf("issueId=%s : err=%v", issueId, err)
		return domain.Issue{}, err
	}

	return issue, nil
}

func CreateIssue(project, issueType, title string) (domain.Issue, error) {
	var issue domain.Issue
	issueRequest := newCreateIssueRequest(project, issueType, title)
	b, err := json.Marshal(issueRequest)
	if err != nil {
		return issue, err
	}

	path := "/rest/api/2/issue/"
	request, err := constructPostRequest(path, string(b))
	if err != nil {
		return issue, err
	}
	statusCode, err := SendRequest(request, &issue)
	if statusCode != http.StatusCreated {
		err = fmt.Errorf("issueTitle='%s' : statusCode=%d", title, statusCode)
		return domain.Issue{}, err
	}
	if err != nil {
		err = fmt.Errorf("issueTitle=%s : err=%v", title, err)
		return domain.Issue{}, err
	}
	return issue, nil
}

func newCreateIssueRequest(projectKey, issueType, summary string) IssueRequest {
	project := ProjectRequest{projectKey}
	issue := IssueTypeRequest{issueType}

	fields := RequestFields{Project: project, Summary: summary, IssueType: issue}
	return IssueRequest{Fields: fields}
}

func constructGetRequest(path string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, getUrl(path), nil)
	if err != nil {
		return request, err
	}

	request.Header.Set("Accept", "application/json")
	if JiraConfig.Username != "" {
		request.SetBasicAuth(JiraConfig.Username, JiraConfig.Password)
	}

	return request, err
}

func constructPostRequest(path, data string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodPost, getUrl(path), bytes.NewBuffer([]byte(data)))
	if err != nil {
		return request, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	if JiraConfig.Username != "" {
		request.SetBasicAuth(JiraConfig.Username, JiraConfig.Password)
	}
	return request, err
}

func getUrl(path string) string {
	path = strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(JiraConfig.Host, "/"), path)
}
