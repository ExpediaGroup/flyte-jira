package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var (
	GetTransitions = flyte.Command{
		Name:         "GetTransitions",
		OutputEvents: []flyte.EventDef{getTransitionsEventDef, getTransitionsFailureEventDef},
		Handler:      getTransitionsHandler,
	}

	getTransitionsEventDef = flyte.EventDef{
		Name: "GetTransitions",
	}

	getTransitionsFailureEventDef = flyte.EventDef{
		Name: "GetTransitionsFailure",
	}
)

type (
	transitionRequest struct {
		IssueId string `json:"issueId"`
	}

	transitionsSuccessPayload struct {
		Id      string              `json:"id"`
		Results []client.Transition `json:"transitions"`
	}

	transitionsFailurePayload struct {
		Id    string `json:"id"`
		Error string `json:"error"`
	}
)

func getTransitionsHandler(input json.RawMessage) flyte.Event {
	req := transitionRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling getting transitions request [%s]: %s", input, err)
		return getTransitionsFailureEvent(req, err)
	}

	results, err := client.GetTransitions(req.IssueId)

	log.Printf("results %v", results.Transitions)

	if err != nil {
		log.Printf("Error getting transitions for issue %s: %s", req.IssueId, err)
		return getTransitionsFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: getTransitionsEventDef,
		Payload: transitionsSuccessPayload{
			Id:      req.IssueId,
			Results: results.Transitions,
		},
	}
}

func getTransitionsFailureEvent(request transitionRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: getTransitionsFailureEventDef,
		Payload: transitionsFailurePayload{
			request.IssueId,
			err.Error(),
		},
	}
}
