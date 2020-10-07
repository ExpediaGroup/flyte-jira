package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var (
	DoTransition = flyte.Command{
		Name:         "DoTransition",
		OutputEvents: []flyte.EventDef{doTransitionEventDef, doTransitionFailureEventDef},
		Handler:      doTransitionHandler,
	}

	doTransitionEventDef = flyte.EventDef{
		Name: "DoTransition",
	}

	doTransitionFailureEventDef = flyte.EventDef{
		Name: "DoTransitionFailure",
	}
)

type (
	doTransitionRequest struct {
		IssueId      string `json:"issueId"`
		TransitionId string `json:"transitionId"`
	}

	doTransitionFailurePayload struct {
		IssueId      string `json:"id"`
		TransitionId string `json:"transitionId"`
		Error        string `json:"error"`
	}
)

func doTransitionHandler(input json.RawMessage) flyte.Event {
	req := doTransitionRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling transition request [%s]: %s", input, err)
		return doTransitionFailureEvent(req, err)
	}

	err := client.DoTransition(req.IssueId, req.TransitionId)

	if err != nil {
		log.Printf("Error during a transition for issue %s: %s", req.IssueId, err)
		return doTransitionFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: doTransitionEventDef,
		Payload: doTransitionRequest{
			IssueId:      req.IssueId,
			TransitionId: req.TransitionId,
		},
	}

}

func doTransitionFailureEvent(request doTransitionRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: doTransitionFailureEventDef,
		Payload: doTransitionFailurePayload{
			request.IssueId,
			request.TransitionId,
			err.Error(),
		},
	}
}
