package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var (
	GetStatuses = flyte.Command{
		Name:         "GetStatuses",
		OutputEvents: []flyte.EventDef{getStatusesEventDef, getStatusesFailureEventDef},
		Handler:      getStatusesHandler,
	}

	getStatusesEventDef = flyte.EventDef{
		Name: "GetStatuses",
	}

	getStatusesFailureEventDef = flyte.EventDef{
		Name: "GetStatusesFailure",
	}
)

type (
	getStatusesPayload struct {
		statusRequest
		Error string `json:"error"`
	}

	statusRequest struct {
		ProjectId string `json:"projectId"`
	}
)

func getStatusesHandler(input json.RawMessage) flyte.Event {
	req := statusRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling getting statuses request [%s]: %s", input, err)
		return getStatusesFailureEvent(req, err)
	}

	if err := client.GetStatuses(req.ProjectId); err != nil {
		log.Printf("Error getting statuses for project %s: %s", req.ProjectId, err)
		return getStatusesFailureEvent(req, err)
	}
	return flyte.Event{
		EventDef: getStatusesEventDef,
		Payload:  req,
	}
}

func getStatusesFailureEvent(request statusRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: getStatusesFailureEventDef,
		Payload: getStatusesPayload{
			request,
			err.Error(),
		},
	}
}
