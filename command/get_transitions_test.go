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
	"github.com/HotelsDotCom/flyte-client/flyte"
	"net/http"
	"reflect"
	"testing"
)

func TestGetTransitionsWorkingAsExpected(t *testing.T) {
	initialFunc := client.SendRequest
	defer func() { client.SendRequest = initialFunc }()
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusOK, nil
	}
	var transitionRequest = struct {
		IssueId string `json:"issueId"`
	}{"DEVEX-567"}

	input := toJson(transitionRequest, t)

	actualEvent := getTransitionsHandler((input))
	expectedEvent := flyte.Event{
		EventDef: getTransitionsEventDef,
		Payload: transitionsSuccessPayload{
			Id:      "DEVEX-567",
			Results: nil,
		},
	}
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}

func TestGetTransitionsFailure(t *testing.T) {
	initialFunc := client.SendRequest
	defer func() { client.SendRequest = initialFunc }()
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		reqPath := request.URL.Path
		expReqPath := "/rest/api/2/issue/DEVEX-567/transitions"
		if reqPath != expReqPath {
			return http.StatusNotFound, nil
		}
		return http.StatusOK, nil
	}

	var transitionRequest = struct {
		IssueId string `json:"issueId"`
	}{"DEVEX-5677777"}

	input := toJson(transitionRequest, t)

	actualEvent := getTransitionsHandler((input))
	expectedEvent := flyte.Event{
		EventDef: getTransitionsEventDef,
		Payload: transitionsSuccessPayload{
			Id:      "DEVEX-5677777",
			Results: nil,
		},
	}
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}
