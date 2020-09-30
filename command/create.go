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
	"fmt"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
)

var CreateIssueCommand = flyte.Command{
	Name:         "CreateIssue",
	OutputEvents: []flyte.EventDef{createEventDef, createFailureEventDef},
	Handler:      createIssueHandler,
}

func createIssueHandler(input json.RawMessage) flyte.Event {
	var handlerInput struct {
		Project   string `json:"project"`
		IssueType string `json:"issuetype"`
		Title     string `json:"summary"`
	}

	if err := json.Unmarshal(input, &handlerInput); err != nil {
		err := fmt.Errorf("Could not marshal create client issue input: %s", err)
		log.Println(err)
		return newCreateFailureEvent(err.Error(), "unknown", "unknown", "unkown")
	}
	issue, err := client.CreateIssue(handlerInput.Project, handlerInput.IssueType, handlerInput.Title)
	if err != nil {
		err = fmt.Errorf("Could not create issue: %v", err)
		log.Println(err)
		return newCreateFailureEvent(err.Error(), handlerInput.Project, handlerInput.IssueType, handlerInput.Title)
	}
	return newCreateEvent(fmt.Sprintf("%s/browse/%s", client.JiraConfig.Host, issue.Key), issue.Key, handlerInput.Project, handlerInput.IssueType, handlerInput.Title)
}

var createEventDef = flyte.EventDef{
	Name: "CreateIssue",
}

type createSuccessPayload struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Project   string `json:"project"`
	IssueType string `json:"issuetype"`
	Title     string `json:"summary"`
}

var createFailureEventDef = flyte.EventDef{
	Name: "CreateIssueFailure",
}

type createFailurePayload struct {
	Error     string `json:"error"`
	Project   string `json:"project"`
	IssueType string `json:"issuetype"`
	Title     string `json:"summary"`
}

func newCreateFailureEvent(err, project, issueType, title string) flyte.Event {
	return flyte.Event{
		EventDef: createFailureEventDef,
		Payload: createFailurePayload{
			Error:     err,
			Project:   project,
			IssueType: issueType,
			Title:     title,
		},
	}
}

func newCreateEvent(url, id, project, issueType, title string) flyte.Event {
	return flyte.Event{
		EventDef: createEventDef,
		Payload: createSuccessPayload{
			Id:        id,
			Url:       url,
			Project:   project,
			IssueType: issueType,
			Title:     title,
		},
	}
}
