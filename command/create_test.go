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
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateIssueAsExpected(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusCreated, nil
	}
	input := []byte(`{"project":"FLYTE","issuetype":"Story", "summary": "test story","description": "test description", "priority": "Medium", "reporter": "songupta"}`)
	actualEvent := createIssueHandler(input)
	expectedEvent := newCreateIssueEvent("/browse/", "", "FLYTE", "Story", "test story", "test description", "Medium", "songupta")
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestCreateIssueFailure(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}
	input := []byte(`{"project":"FLYTE","issuetype":"Story", "summary": "test story"}`)
	actualEvent := createIssueHandler(input)
	expectedEvent := newCreateIssueFailureEvent("Could not create issue: issueSummary='test story' : statusCode=400", "FLYTE", "Story", "test story")
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestCreateCustomIssueAsExpected(t *testing.T) {
	client.SendCustomRequest = func(request *http.Request) ([]byte, error) {
		return []byte(`{}`), nil
	}
	input := []byte(`{"project":"FLYTE","issuetype":"Story", "summary": "test story", "incident":"INC1234567"}`)
	actualEvent := createIncIssueHandler(input)
	expectedEvent := flyte.Event{
		EventDef: createIncIssueEventDef,
		Payload: CreateIncIssueSuccess{
			ID:   "",
			Key:  "",
			Self: "",
		},
	}
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}
