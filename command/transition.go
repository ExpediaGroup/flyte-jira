package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var (
	Transition = flyte.Command{
		Name:         "Transition",
		OutputEvents: []flyte.EventDef{transitionEventDef, transitionFailureEventDef},
		Handler:      transitionHandler,
	}

	transitionEventDef = flyte.EventDef{
		Name: "Transition",
	}

	transitionFailureEventDef = flyte.EventDef{
		Name: "TransitionFailure",
	}
)

type transitionRequest struct {
	IssueId      string `json:"issueId"`
	TransitionId string `json:"transitionId"`
}

type transitionPayload struct {
	IssueId      string `json:"issueId"`
	TransitionId string `json:"transitionId"`
	RequestURL   string `json:"requestURL"`
}

type transitionFailurePayload struct {
	IssueId      string `json:"issueId"`
	TransitionId string `json:"transitionId"`
	Error        string `json:"error"`
}

func transitionHandler(input json.RawMessage) flyte.Event {
	req := transitionRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling transition request [%s]: %s", input, err)
		return transitionFailureEvent(req, err)
	}

	err, reqURL := client.Transition(req.IssueId, req.TransitionId)

	if err != nil {
		log.Printf("Error during a transition for issue %s: %s", req.IssueId, err)
		return transitionFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: transitionEventDef,
		Payload: transitionPayload{
			IssueId:      req.IssueId,
			TransitionId: req.TransitionId,
			RequestURL:   reqURL,
		},
	}

}

func transitionFailureEvent(request transitionRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: transitionFailureEventDef,
		Payload: transitionFailurePayload{
			request.IssueId,
			request.TransitionId,
			err.Error(),
		},
	}
}
