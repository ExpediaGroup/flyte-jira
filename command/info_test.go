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
	"github.com/ExpediaGroup/flyte-jira/domain"
	"net/http"
	"reflect"
	"testing"
)

type infoTest struct {
	name string
	rawJson string
}

func TestGetInfoWorkingAsExpected(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusOK, nil
	}

	testCases := []infoTest {
		{"test-normal-input", `"Test"`},
		{"test-url-input", `"http://test123.com/Test"`},
		{"test-slack-url-input", `"<http://test123.com/Test>"`},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			in := tCase.rawJson
			event := infoHandler([]byte(in))

			expectedEvent := newInfoEvent(domain.Issue{})
			if !reflect.DeepEqual(event, expectedEvent) {
				t.Errorf("Expected: %v but got: %v", expectedEvent, event)
			}
		})
	}
}

func TestGetInfoFailure(t *testing.T) {
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}
	input := toJson("Test", t)
	actualEvent := infoHandler(input)

	// Issue empty because it's populated in Send request
	expectedEvent := newInfoFailureEvent("could not get info: issueId=Test : statusCode=400", "Test")
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}
