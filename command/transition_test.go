package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockRequestBody struct {
	TransitionId string `json:"id"`
}

type TransitionRequest struct {
	Transition mockRequestBody `json:"transition"`
}

func createMockSendRequest(issueId, transitionId string) func(request *http.Request) (int, error) {
	return func(request *http.Request) (int, error) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return 400, err
		}

		body := &TransitionRequest{}

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
func TestTransitionAsExpected(t *testing.T) {
	prevState := client.SendRequestWithoutResp
	defer func() {
		client.SendRequestWithoutResp = prevState
	}()
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-123", "881")
	input := []byte(`{"issueId":"DEVEX-123","transitionId":"881"}`)
	actualEvent := transitionHandler(input)
	expEvent := flyte.Event{
		EventDef: transitionEventDef,
		Payload: transitionPayload{
			IssueId:      "DEVEX-123",
			TransitionId: "881",
			RequestURL:   "/rest/api/2/issue/DEVEX-123/transitions",
		},
	}
	assert.Equal(t, expEvent, actualEvent)
}

func TestTransitionFailure_IssueOrUserDoesNotExist(t *testing.T) {
	prevState := client.SendRequestWithoutResp
	defer func() {
		client.SendRequestWithoutResp = prevState
	}()
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-12333333", "881")
	input := []byte(`{"issueId":"DEVEX-12333333","transitionId":"881"}`)
	actual := transitionHandler(input)
	exp := flyte.Event{
		EventDef: transitionFailureEventDef,
		Payload: transitionFailurePayload{
			IssueId:      "DEVEX-12333333",
			TransitionId: "881",
			Error:        "issue or user does not exist",
		},
	}
	assert.Equal(t, exp, actual)
}

func TestTransitionFailure_TransitionDoesNotExist(t *testing.T) {
	prevState := client.SendRequestWithoutResp
	defer func() {
		client.SendRequestWithoutResp = prevState
	}()
	client.SendRequestWithoutResp = createMockSendRequest("DEVEX-123", "881")
	input := []byte(`{"issueId":"DEVEX-123","transitionId":"123456"}`)
	actual := transitionHandler(input)
	exp := flyte.Event{
		EventDef: transitionFailureEventDef,
		Payload: transitionFailurePayload{
			IssueId:      "DEVEX-123",
			TransitionId: "123456",
			Error:        "unsupported status code 500",
		},
	}
	assert.Equal(t, exp, actual)
}
