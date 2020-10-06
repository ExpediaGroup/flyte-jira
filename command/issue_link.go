package command

import (
	"encoding/json"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var (
	IssueCreateLinkCommand = flyte.Command{
		Name:         "IssueCreateLink",
		OutputEvents: []flyte.EventDef{linkEventDef, linkFailureEventDef},
		Handler:      issueCreateLinkHandler,
	}

	IssueGetLinkCommand = flyte.Command{
		Name:         "IssueGetLink",
		OutputEvents: []flyte.EventDef{linkEventDef, linkFailureEventDef},
		Handler:      issueGetLinkHandler,
	}

	IssueDeleteLinkCommand = flyte.Command{
		Name:         "IssueDeleteLink",
		OutputEvents: []flyte.EventDef{linkEventDef, linkFailureEventDef},
		Handler:      issueDeleteLinkHandler,
	}

	linkEventDef = flyte.EventDef{
		Name: "Link",
	}

	linkFailureEventDef = flyte.EventDef{
		Name: "LinkFailure",
	}
)

type (
	linkFailurePayload struct {
		linkRequest
		Error string `json:"error"`
	}

	linkRequest struct {
		LinkId       string `json:"linkId,omitempty"`
		InwardIssue  string `json:"inwardIssue,omitempty"`
		OutwardIssue string `json:"outwardIssue,omitempty"`
		LinkType     string `json:"linkType,omitempty"`
	}
)

func issueGetLinkHandler(input json.RawMessage) flyte.Event {
	req := linkRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling Issue Link Request [%s]: %s", input, err)
		return newLinkFailureEvent(req, err)
	}

	resp, err := client.GetLink(req.LinkId)
	if err != nil {
		log.Printf("Error fetching link %s: %s", req.LinkId, err)
		return newLinkFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: linkEventDef,
		Payload:  resp,
	}
}

func issueCreateLinkHandler(input json.RawMessage) flyte.Event {
	req := linkRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling Issue Link Request [%s]: %s", input, err)
		return newLinkFailureEvent(req, err)
	}

	if err := client.LinkIssues(req.InwardIssue, req.OutwardIssue, req.LinkType); err != nil {
		log.Printf("Error linking Issue %s to Issue %s with type %s: %s", req.InwardIssue, req.OutwardIssue, req.LinkType, err)
		return newLinkFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: linkEventDef,
		Payload:  req,
	}
}

func issueDeleteLinkHandler(input json.RawMessage) flyte.Event {
	req := linkRequest{}
	if err := json.Unmarshal(input, &req); err != nil {
		log.Printf("Error unmarshaling Issue Link Request [%s]: %s", input, err)
		return newLinkFailureEvent(req, err)
	}

	if err := client.DeleteLink(req.LinkId); err != nil {
		log.Printf("Error removing link %s: %s", req.LinkId, err)
		return newLinkFailureEvent(req, err)
	}

	return flyte.Event{
		EventDef: linkEventDef,
		Payload:  req,
	}
}

func newLinkFailureEvent(request linkRequest, err error) flyte.Event {
	return flyte.Event{
		EventDef: linkFailureEventDef,
		Payload: linkFailurePayload{
			request,
			err.Error(),
		},
	}
}
