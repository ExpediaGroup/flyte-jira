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
	"github.com/ExpediaGroup/flyte-jira/client"
	"net/http"
	"reflect"
	"testing"
)

func TestDoTransitionAsExpected(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusNoContent, nil
	}

	var inputStruct = struct {
		IssueId      string `json:"issueId"`
		TransitionId string `json:"transitionId"`
	}{"DEVEX-123", "881"}
	input := toJson(inputStruct, t)

	actualEvent := doTransitionHandler(input)
	expectedEvent := doTransitionEvent(inputStruct)
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}

func TestDoTransitionFailure(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}

	var inputStruct = struct {
		IssueId      string `json:"issueId"`
		TransitionId string `json:"transitionId"`
	}{"DEVEX-123", "88661"}
	input := toJson(inputStruct, t)

	actualEvent := doTransitionHandler(input)
	expectedEvent := doTransitionFailureEvent(inputStruct, nil)
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %+v but got: %+v", expectedEvent, actualEvent)
	}
}
