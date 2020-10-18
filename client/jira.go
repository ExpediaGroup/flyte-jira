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
	"errors"
	"fmt"
	"github.com/ExpediaGroup/flyte-jira/domain"
	"net/http"
	"strings"
)

// Must be initialised before using
var JiraConfig Config

type (
	Config struct {
		Host     string
		Username string
		Password string
	}

	Comment struct {
		Body string `json:"body"`
	}

	Issue struct {
		Fields Fields `json:"fields"`
	}

	Fields struct {
		Project   Project `json:"project"`
		Summary   string  `json:"summary"`
		IssueType Type    `json:"issuetype"`
	}

	Project struct {
		Key string `json:"key"`
	}

	Type struct {
		Name string `json:"name"`
	}

	SearchRequestType struct {
		Query      string   `json:"jql"`
		StartIndex int      `json:"startAt"`
		MaxResults int      `json:"maxResults"`
		Fields     []string `json:"fields"`
	}

	SearchResult struct {
		StartIndex   int            `json:"startAt"`
		MaxResults   int            `json:"maxResults"`
		TotalResults int            `json:"total"`
		Issues       []domain.Issue `json:"issues"'`
	}

	AssignRequest struct {
		Name string `json:"name,omitempty"`
	}

	LinkIssueRequest struct {
		LinkType IssueLink `json:"type"`
		Inward   LinkIssue `json:"inwardIssue"`
		Outward  LinkIssue `json:"outwardIssue"`
		Comment  `json:"comment"`
	}

	LinkIssue struct {
		Key string `json:"key"`
	}

	IssueLink struct {
		Name    string `json:"name"`
		Inward  string `json:"inward,omitempty"`
		Outward string `json:"outward,omitempty"`
	}
)

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
		return domain.Issue{}, fmt.Errorf("issueId=%s : statusCode=%d", issueId, statusCode)
	}
	if err != nil {
		return domain.Issue{}, fmt.Errorf("issueId=%s : err=%s", issueId, err)
	}

	return issue, nil
}

func CreateIssue(project, issueType, summary string) (domain.Issue, error) {
	var issue domain.Issue
	issueRequest := Issue{
		Fields: Fields{
			Project:   Project{Key: project},
			Summary:   summary,
			IssueType: Type{Name: issueType},
		}}
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
		return domain.Issue{}, fmt.Errorf("issueSummary='%s' : statusCode=%d", summary, statusCode)
	}
	if err != nil {
		err = fmt.Errorf("issueSummary=%s : err=%v", summary, err)
		return domain.Issue{}, err
	}
	return issue, nil
}

func SearchIssues(query string, startIndex int, maxResults int) (SearchResult, error) {
	var searchResult SearchResult

	requestBody := newSearchRequestBody(query, startIndex, maxResults)
	encodedBody, err := json.Marshal(requestBody)
	if err != nil {
		return searchResult, err
	}

	path := "/rest/api/2/search"
	request, err := constructPostRequest(path, string(encodedBody))
	if err != nil {
		return searchResult, err
	}

	statusCode, err := SendRequest(request, &searchResult)
	if err != nil {
		err := fmt.Errorf("query='%s' : error=%v", query, err)
		return searchResult, err
	}
	if statusCode != http.StatusOK {
		err := fmt.Errorf("query='%s' : statusCode=%d", query, statusCode)
		return searchResult, err
	}

	return searchResult, nil
}

func AssignIssue(issueId, username string) error {
	path := fmt.Sprintf("/rest/api/2/issue/%s/assignee", issueId)
	b, err := json.Marshal(&AssignRequest{username})
	if err != nil {
		return err
	}

	req, err := constructPutRequest(path, string(b))
	if err != nil {
		return err
	}

	httpCode, err := SendRequestWithoutResp(req)
	if err != nil {
		return err
	}

	switch httpCode {
	case http.StatusNoContent:
		err = nil
	case http.StatusBadRequest:
		err = errors.New("invalid user representation")
	case http.StatusUnauthorized:
		err = errors.New("invalid permission to assign to issue")
	case http.StatusNotFound:
		err = errors.New("issue or user does not exist")
	default:
		err = fmt.Errorf("unsupported status code %d", httpCode)
	}

	return err
}

func LinkIssues(inwardKey, outwardKey, linkType string) error {
	path := "/rest/api/2/issueLink"
	linkReq := LinkIssueRequest{
		LinkType: IssueLink{Name: linkType},
		Inward:   LinkIssue{inwardKey},
		Outward:  LinkIssue{outwardKey},
		Comment:  Comment{"Link related issues!"},
	}
	b, err := json.Marshal(linkReq)
	if err != nil {
		return err
	}

	httpReq, err := constructPostRequest(path, string(b))
	if err != nil {
		return err
	}

	httpCode, err := SendRequestWithoutResp(httpReq)
	if err != nil {
		return err
	}

	return checkHttpCode(httpCode, linkReq.Body)
}

func GetLink(linkId string) (*LinkIssueRequest, error) {
	path := fmt.Sprintf("/rest/api/2/issueLink/%s", linkId)
	httpReq, err := constructGetRequest(path)
	if err != nil {
		return nil, err
	}

	link := &LinkIssueRequest{}
	httpCode, err := SendRequest(httpReq, &link)

	return link, checkHttpCode(httpCode, linkId)
}

//TODO: get rid of all of those separate construct methods...
func DeleteLink(linkId string) error {
	path := fmt.Sprintf("/rest/api/2/issueLink/%s", linkId)
	httpReq, err := constructDeleteRequest(path)
	if err != nil {
		return err
	}

	httpCode, err := SendRequestWithoutResp(httpReq)
	return checkHttpCode(httpCode, linkId)
}

//TODO: change error msg to be more generic and use in other funcs
func checkHttpCode(httpCode int, in string) error {
	if 200 <= httpCode && httpCode <= 208 {
		return nil
	}

	var err error
	switch httpCode {
	case http.StatusBadRequest:
		err = fmt.Errorf("failed to create issue with comment %s", in)
	case http.StatusUnauthorized:
		err = errors.New("invalid permission to link issues")
	case http.StatusInternalServerError:
		err = errors.New("error occurred when creating link or comment")
	case http.StatusNotFound:
		err = errors.New("could not find issue or invalid link type specified")
	default:
		err = fmt.Errorf("unsupported status code %d found", httpCode)
	}

	return err
}

func newSearchRequestBody(query string, startIndex int, maxResults int) SearchRequestType {
	return SearchRequestType{
		query,
		startIndex,
		maxResults,
		[]string{"summary", "assignee", "labels", "status", "description", "priority"},
	}
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

func constructDeleteRequest(path string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodDelete, getUrl(path), nil)
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

func constructPutRequest(path, data string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodPut, getUrl(path), bytes.NewBuffer([]byte(data)))
	if err != nil {
		return request, err
	}

	request.Header.Set("Content-LinkType", "application/json")
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
