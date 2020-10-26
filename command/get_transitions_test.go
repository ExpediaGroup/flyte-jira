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
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestGetTransitionsWorkingAsExpected(t *testing.T) {
	initialFunc := client.SendRequest
	defer func() { client.SendRequest = initialFunc }()
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		path := request.URL.Path
		var subStr = strings.Split(path, "/")
		if subStr[5] != "DEVEX-567" {
			return http.StatusNotFound, errors.New("requested jira issue not found")
		}
		return http.StatusOK, nil
	}
	input := []byte(`{"issueId":"DEVEX-567"}`)
	actualEvent := getTransitionsHandler(input)
	expectedEvent := flyte.Event{
		EventDef: getTransitionsEventDef,
		Payload: transitionsSuccessPayload{
			Id:      "DEVEX-567",
			Results: []client.TransitionObj{},
		},
	}
	assert.Equal(t, expectedEvent, actualEvent)
}

func TestGetTransitionsFailure(t *testing.T) {
	initialFunc := client.SendRequest
	defer func() { client.SendRequest = initialFunc }()
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		reqPath := request.URL.Path
		expReqPath := "/rest/api/2/issue/DEVEX-567/transitions"
		if reqPath != expReqPath {
			return http.StatusNotFound, errors.New("Issue does not exist")
		}
		return http.StatusOK, nil
	}
	input := []byte(`{"issueId":"DEVEX-5677777"}`)
	actualEvent := getTransitionsHandler(input)
	expectedEvent := flyte.Event{
		EventDef: getTransitionsFailureEventDef,
		Payload: transitionsFailurePayload{
			Id:    "DEVEX-5677777",
			Error: "Issue does not exist",
		},
	}
	assert.Equal(t, expectedEvent, actualEvent)
}
