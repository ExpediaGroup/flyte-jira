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

package command

import (
	"errors"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-jira/domain"
	"net/http"
	"reflect"
	"testing"
)

func TestSearchIssuesSuccess(t *testing.T) {
	initialSendRequest := client.SendRequest
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusOK, nil
	}
	defer func() { client.SendRequest = initialSendRequest }()

	actualEvent := searchIssuesHandler([]byte(`{"query": "project = FLYTE"}`))

	expectedEvent := newSearchSuccessEvent(SearchIssuesInput{"project = FLYTE", 0, 10}, 0, nil)

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestSearchIssuesFailure(t *testing.T) {
	initialSendRequest := client.SendRequest
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}
	defer func() { client.SendRequest = initialSendRequest }()

	actualEvent := searchIssuesHandler([]byte(`{"query": "project = FLYTE"}`))
	expectedEvent := newSearchFailureEvent(SearchIssuesInput{"project = FLYTE", 0, 10}, errors.New("Could not search for issues: query='project = FLYTE' : statusCode=400"))

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestSearchIssuesEmptyQuery(t *testing.T) {
	actualEvent := searchIssuesHandler([]byte(`{"query": ""}`))
	expectedEvent := newSearchFailureEvent(SearchIssuesInput{"", 0, 10}, errors.New("Empty query string"))

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestRequestError(t *testing.T) {
	initialSendRequest := client.SendRequest
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return -1, errors.New("request timed out")
	}
	defer func() { client.SendRequest = initialSendRequest }()

	actualEvent := searchIssuesHandler([]byte(`{"query": "project = FLYTE"}`))
	expectedEvent := newSearchFailureEvent(SearchIssuesInput{"project = FLYTE", 0, 10}, errors.New("Could not search for issues: query='project = FLYTE' : error=request timed out"))

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestIssueFormatting(t *testing.T) {
	initialSendRequest := client.SendRequest
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		p := client.SearchResult{
			TotalResults: 2,
			Issues:       []domain.Issue{createDummyIssue(), createDummyIssue()},
		}

		reflect.ValueOf(responseBody).Elem().Set(reflect.ValueOf(&p).Elem())

		return http.StatusOK, nil
	}
	defer func() { client.SendRequest = initialSendRequest }()

	actualEvent := searchIssuesHandler([]byte(`{"query": "project = FLYTE"}`))
	expectedEvent := newSearchSuccessEvent(SearchIssuesInput{"project = FLYTE", 0, 10}, 2, []domain.Issue{createDummyIssue(), createDummyIssue()})

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestInvalidInput(t *testing.T){
	actualEvent := searchIssuesHandler([]byte(`{]`))
	expectedEvent := flyte.NewFatalEvent(errors.New("input is not valid: invalid character ']' looking for beginning of object key string"))

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}
