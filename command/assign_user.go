package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"log"
)

var (
	IssueAssignCommand = flyte.Command{
		Name:         "IssueAssign",
		OutputEvents: []flyte.EventDef{assignEventDef, assignFailureEventDef},
		Handler:      assignIssueHandler,
	}

	assignEventDef = flyte.EventDef{
		Name: "Assign",
	}

	assignFailureEventDef = flyte.EventDef{
		Name: "AssignFailure",
	}
)

type (
	assignFailurePayload struct {
		assignRequest
		Error string `json:"error"`
	}

	assignRequest struct {
		IssueId  string `json:"issueId"`
		Username string `json:"username,omitempty"`
	}
)

func assignIssueHandler(input json.RawMessage) flyte.Event {
	req := assignRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling Issue Assign Request [%s]: %s", input, err)
		return newAssignFailureEvent(req, err)
	}

	if err := client.AssignIssue(req.IssueId, req.Username); err != nil {
		log.Printf("Error assigning Issue %s to User %s: %s", req.IssueId, req.Username, err)
		return newAssignFailureEvent(req, err)
	}
	return flyte.Event{
		EventDef: assignEventDef,
		Payload:  req,
	}
}

func newAssignFailureEvent(request assignRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: assignFailureEventDef,
		Payload: assignFailurePayload{
			request,
			err.Error(),
		},
	}
}
