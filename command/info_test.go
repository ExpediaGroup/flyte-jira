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
	"fmt"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/ExpediaGroup/flyte-jira/domain"
	"net/http"
	"path"
	"reflect"
	"testing"
)

type infoTest struct {
	name       string
	rawJson    string
	expIssueId string
}

func TestGetInfoWorkingAsExpected(t *testing.T) {
	initialFunc := client.SendRequest
	defer func(){ client.SendRequest = initialFunc}()
	expectedIssueId := ""
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		reqPath := request.URL.Path
		issueId := path.Base(reqPath)
		if issueId != expectedIssueId {
			return http.StatusBadRequest, fmt.Errorf("expected issueId %s got %s", "DEVEX-553", issueId)
		}
		return http.StatusOK, nil
	}

	testCases := []infoTest{
		{"test-normal-input",
			`"DEVEX-553"`,
			"DEVEX-553",
		},
		{
			"test-url-input",
			`"http://test123.com/TeSt-75122"`,
			"TeSt-75122",
		},
		{
			"test-url-without-valid-base",
			`"https://jira.expedia.biz/browse/ELF-21462?jql=project%20%3D%20ELS%20AND%20resolution%20%3D%20Unresolved%20AND%20text%20~%20%22symantec%22%20ORDER%20BY%20priority%20DESC%2C%20updated%20DESC"`,
			"ELF-21462",
		},
		{"test-url-with-nested-path",
			`"https://jira.expedia.biz/servicedesk/customer/portal/518/FOOBARFOOBAR-7773"`,
			"FOOBARFOOBAR-7773",
		},
		{"test-external-url",
			`"https://jira.expedia.biz/browse/TEST-221"`,
			"TEST-221",
		},
		{"test-slack-url-input",
			`"<http://test123.com/ELS-790>"`,
			"ELS-790",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			in := tCase.rawJson
			expectedIssueId = tCase.expIssueId
			event := infoHandler([]byte(in))

			expectedEvent := newInfoEvent(domain.Issue{})
			if !reflect.DeepEqual(event, expectedEvent) {
				t.Errorf("Expected: %v but got: %v", expectedEvent, event)
			}
		})
	}
}

func TestGetInfoFailure(t *testing.T) {
	initialFunc := client.SendRequest
	defer func(){ client.SendRequest = initialFunc}()
	client.SendRequest = func(request *http.Request, responseBody interface{}) (int, error) {
		return http.StatusBadRequest, nil
	}
	input := `"Test-123"`
	actualEvent := infoHandler([]byte(input))

	// Issue empty because it's populated in Send request
	expectedEvent := newInfoFailureEvent("Test-123", fmt.Errorf("issueId=%s : statusCode=%d", "Test-123", 400))
	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}
