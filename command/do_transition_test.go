package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type (
	mockRequestBody struct {
		TransitionId string `json:"id"`
	}

	DoTransitionRequest struct {
		Transition mockRequestBody `json:"transition"`
	}
)

func createMockSendRequest(issueId, transitionId string) func(request *http.Request) (int, error) {
	return func(request *http.Request) (int, error) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return 400, err
		}

		body := &DoTransitionRequest{}

		if err := json.Unmarshal(b, &body); err != nil {
			return 400, err
		}

		if body.Transition.TransitionId != transitionId {
			return 500, err
		}

		if issueId != "DEVEX-123" {
			return 404, err
		}

		return 204, nil
	}
}
func TestDoTransitionAsExpected(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-123", "881")
	input := []byte(`{"issueId":"DEVEX-123","transitionId":"881"}`)
	actual := doTransitionHandler(input)
	exp := flyte.Event{
		EventDef: doTransitionEventDef,
		Payload: doTransitionRequest{
			IssueId:      "DEVEX-123",
			TransitionId: "881",
		},
	}

	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}

func TestDoTransitionFailure_IssueOrUserDoesNotExist(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-12333333", "881")
	input := []byte(`{"issueId":"DEVEX-12333333","transitionId":"881"}`)
	actual := doTransitionHandler(input)
	exp := flyte.Event{
		EventDef: doTransitionFailureEventDef,
		Payload: doTransitionFailurePayload{
			IssueId:      "DEVEX-12333333",
			TransitionId: "881",
			Error:        "issue or user does not exist",
		},
	}

	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}

func TestDoTransitionFailure_TransitionDoesNotExist(t *testing.T) {
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-123", "881")
	input := []byte(`{"issueId":"DEVEX-123","transitionId":"123456"}`)
	actual := doTransitionHandler(input)
	exp := flyte.Event{
		EventDef: doTransitionFailureEventDef,
		Payload: doTransitionFailurePayload{
			IssueId:      "DEVEX-123",
			TransitionId: "123456",
			Error:        "unsupported status code 500",
		},
	}

	if !reflect.DeepEqual(actual, exp) {
		t.Errorf("Expected: %v but got: %v", exp, actual)
	}
}
