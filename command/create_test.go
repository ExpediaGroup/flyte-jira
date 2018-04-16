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
	"net/http"
	"reflect"
	"github.com/HotelsDotCom/flyte-jira/client"
	"testing"
)

func TestCreateTicketAsExpected(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusCreated, nil
	}

	var inputStruct = struct {
		Project   string `json:"project"`
		IssueType string `json:"issueType"`
		Title     string `json:"title"`
	}{"FLYTE", "Story", "test story"}
	input := toJson(inputStruct, t)

	actualEvent := createTicketHandler(input)
	expectedEvent := newCreateEvent("/browse/", "", "FLYTE", "Story", "test story")
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestCreateTicketFailure(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}

	var inputStruct = struct {
		Project   string `json:"project"`
		IssueType string `json:"issueType"`
		Title     string `json:"title"`
	}{"FLYTE", "Story", "test story"}
	input := toJson(inputStruct, t)

	actualEvent := createTicketHandler(input)
	expectedEvent := newCreateFailureEvent("Could not create ticket: ticketTitle='test story' : statusCode=400", "FLYTE", "Story", "test story")
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}
