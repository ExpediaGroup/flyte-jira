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
	"net/http"
	"github.com/HotelsDotCom/flyte-jira/domain"
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

type TicketRequest struct {
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

func CommentTicket(ticketId, comment string) (domain.Ticket, error) {
	var ticket domain.Ticket
	b, err := json.Marshal(Comment{comment})
	if err != nil {
		return ticket, err
	}

	path := fmt.Sprintf("/rest/api/2/issue/%s/comment", ticketId)
	request, err := constructPostRequest(path, string(b))
	if err != nil {
		return ticket, err
	}

	statusCode, err := SendRequest(request, &ticket)
	if statusCode != http.StatusCreated {
		err = fmt.Errorf("ticketId=%s : statusCode=%d", ticketId, statusCode)
		return domain.Ticket{}, err
	}
	if err != nil {
		err = fmt.Errorf("ticketId=%s : err=%v", ticketId, err)
		return domain.Ticket{}, err
	}

	return ticket, nil
}

func GetTicketInfo(ticketId string) (domain.Ticket, error) {
	var ticket domain.Ticket
	path := fmt.Sprintf("/rest/api/2/issue/%s", ticketId)

	request, err := constructGetRequest(path)
	if err != nil {
		return ticket, err
	}

	statusCode, err := SendRequest(request, &ticket)
	if statusCode != http.StatusOK {
		err = fmt.Errorf("ticketId=%s : statusCode=%d", ticketId, statusCode)
		return domain.Ticket{}, err
	}
	if err != nil {
		err = fmt.Errorf("ticketId=%s : err=%v", ticketId, err)
		return domain.Ticket{}, err
	}

	return ticket, nil
}

func CreateTicket(project, issueType, title string) (domain.Ticket, error) {
	var ticket domain.Ticket
	ticketRequest := newCreateTicketRequest(project, issueType, title)
	b, err := json.Marshal(ticketRequest)
	if err != nil {
		return ticket, err
	}

	path := "/rest/api/2/issue/"
	request, err := constructPostRequest(path, string(b))
	if err != nil {
		return ticket, err
	}
	statusCode, err := SendRequest(request, &ticket)
	if statusCode != http.StatusCreated {
		err = fmt.Errorf("ticketTitle='%s' : statusCode=%d", title, statusCode)
		return domain.Ticket{}, err
	}
	if err != nil {
		err = fmt.Errorf("ticketTitle=%s : err=%v", title, err)
		return domain.Ticket{}, err
	}
	return ticket, nil
}

func newCreateTicketRequest(projectKey, issueType, summary string) TicketRequest {
	project := ProjectRequest{projectKey}
	issue := IssueTypeRequest{issueType}

	fields := RequestFields{Project: project, Summary: summary, IssueType: issue}
	return TicketRequest{Fields: fields}
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
