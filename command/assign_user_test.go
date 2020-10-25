package command

import (
	"encoding/json"
	"errors"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/ExpediaGroup/flyte-client/flyte"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type mockReqBody struct {
	Name string `json:"name,omitempty"`
}

func createMockSendReq(testName string) func(request *http.Request) (int, error) {
	return func(request *http.Request) (int, error) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return 400, err
		}

		body := &mockReqBody{}
		if err := json.Unmarshal(b, &body); err != nil {
			return 400, err
		}

		if body.Name != testName {
			return 404, errors.New("invalid name")
		}

		return 204, nil
	}
}
func TestSuccessfulCommand(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendReq("test-123")
	input := []byte(`{"issueId":"foo","username":"test-123"}`)
	actual := assignIssueHandler(input)
	exp := flyte.Event{
		EventDef: assignEventDef,
		Payload: assignRequest{
			IssueId:  "foo",
			Username: "test-123",
		},
	}

	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}

/*
 When the name is nil, issue will be unassigned
 https://docs.atlassian.com/software/jira/docs/api/REST/7.6.1/#api/2/issue-assign
*/
func TestSuccessfulCommandWithNoUser(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendReq("")
	input := []byte(`{"issueId":"foo"}`)
	actual := assignIssueHandler(input)
	exp := flyte.Event{
		EventDef: assignEventDef,
		Payload: assignRequest{
			IssueId: "foo",
		},
	}

	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}

func TestFailedCommand(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendReq("")
	input := []byte(`{"issueId":"foo", "username":"test-123"}`)
	actual := assignIssueHandler(input)
	exp := newAssignFailureEvent(
		assignRequest{
			IssueId:  "foo",
			Username: "test-123",
		},
		errors.New("invalid name"))
	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}
