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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-jira/domain"
)

var IssueInfoCommand = flyte.Command{
	Name:         "IssueInfo",
	OutputEvents: []flyte.EventDef{infoEventDef, infoFailureEventDef},
	Handler:      infoHandler,
}

func infoHandler(input json.RawMessage) flyte.Event {
	var id string
	issue := domain.Issue{}

	if err := json.Unmarshal(input, &id); err != nil {
		err := errors.New(fmt.Sprintf("Could not marshal issue id: %s", err))
		log.Println(err)
		return newInfoFailureEvent(err.Error(), "unkown")
	}

	issue, err := client.GetIssueInfo(id)
	if err != nil {
		err := errors.New(fmt.Sprintf("Could not get info: %v", err))
		log.Println(err)
		return newInfoFailureEvent(err.Error(), id)
	}
	return newInfoEvent(issue)
}

var infoEventDef = flyte.EventDef{
	Name: "Info",
}

var infoFailureEventDef = flyte.EventDef{
	Name: "InfoFailure",
}

func newInfoFailureEvent(err, id string) flyte.Event {
	return flyte.Event{
		EventDef: infoFailureEventDef,
		Payload:  infoFailurePayload{id, err},
	}
}

func newInfoEvent(t domain.Issue) flyte.Event {
	return flyte.Event{
		EventDef: infoEventDef,
		Payload: infoSuccessPayload{
			Id:          t.Key,
			Summary:     t.Fields.Summary,
			Status:      t.Fields.Status.Name,
			Description: t.Fields.Description,
			Assignee:    t.Fields.Assignee.Name,
		},
	}
}

type infoSuccessPayload struct {
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
}

type infoFailurePayload struct {
	Id    string `json:"id"`
	Error string `json:"error"`
}
