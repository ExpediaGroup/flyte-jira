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
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"log"
	"regexp"
	"strings"
)

// Input struct presents input options for flyte command
type Input struct {
	Project     string   `json:"project"`
	IssueType   string   `json:"issuetype"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Inc         string   `json:"incident"` // ServiceNow incident
	Priority    string   `json:"priority"`
	Reporter    string   `json:"reporter"`
}

// CreateIssueCommand is a default command to create issue with minimum parameters
var CreateIssueCommand = flyte.Command{
	Name:         "CreateIssue",
	OutputEvents: []flyte.EventDef{createIssueEventDef, createIssueFailureEventDef},
	Handler:      createIssueHandler,
}

// CreateIncIssueCommand will require an incident argument specified to create jira issue. This is a custom setup for NOCBotV2 app,
// but can be reused anywhere for same purposes
var CreateIncIssueCommand = flyte.Command{
	Name:         "CreateIncIssue",
	OutputEvents: []flyte.EventDef{createIncIssueEventDef, createIncIssueFailureEventDef},
	Handler:      createIncIssueHandler,
}

func createIssueHandler(input json.RawMessage) flyte.Event {
	handlerInput := Input{}
	if err := json.Unmarshal(input, &handlerInput); err != nil {
		err := fmt.Errorf("Could not marshal create client issue input: %s", err)
		log.Println(err)
		return newCreateIssueFailureEvent(err.Error(), "unknown", "unknown", "unkown")
	}
	if (handlerInput.Summary == "" || handlerInput.Description == "") && (handlerInput.Project == "RCPSUP") {
		err := fmt.Errorf("Please provide both issue title & description. mandatory fields missing!  ")
		log.Println(err)
		return newCreateIssueFailureEvent(err.Error(), handlerInput.Project, handlerInput.Description, handlerInput.Summary)
	}
	updateSummaryForReactionBasedRequest(&handlerInput)
	issue, err := client.CreateIssue(handlerInput.Project, handlerInput.IssueType, handlerInput.Summary, handlerInput.Description, handlerInput.Priority, handlerInput.Reporter)
	if err != nil {
		err = fmt.Errorf("Could not create issue: %v", err)
		log.Println(err)
		return newCreateIssueFailureEvent(err.Error(), handlerInput.Project, handlerInput.IssueType, handlerInput.Summary)
	}
	return newCreateIssueEvent(fmt.Sprintf("%s/browse/%s", client.JiraConfig.Host, issue.Key), issue.Key, handlerInput.Project, handlerInput.IssueType, handlerInput.Summary, handlerInput.Description, handlerInput.Priority, handlerInput.Reporter)
}

// createIncIssueHandler handles CreateIncIssue NocBotV2 command and returns success/fail flyte.Event
func createIncIssueHandler(input json.RawMessage) flyte.Event {
	handlerInput := Input{}
	if err := json.Unmarshal(input, &handlerInput); err != nil {
		err := fmt.Errorf("could not marshal create client issue input: %s", err)
		log.Println(err)
		return newCreateIssueFailureEvent(err.Error(), "unknown", "unknown", "unknown")
	}
	log.Println(fmt.Sprintf("Create JIRA issue command recieved. Params: %+v", handlerInput))

	// Input validation
	if !incidentPattern(handlerInput.Inc) { // inc name should be valid
		return flyte.Event{
			EventDef: createIncIssueFailureEventDef,
			Payload: CreateIncIssueFailure{
				Message: fmt.Sprintf("%s: invalid incident number format", handlerInput.Inc),
			},
		}
	}

	if len(handlerInput.Summary) > 255 { // JIRA summary symbols limit is 255
		return flyte.Event{
			EventDef: createIncIssueFailureEventDef,
			Payload: CreateIncIssueFailure{
				Message: fmt.Sprintf("Too long summary (lenght %d). Limit is 255 symbols", len(handlerInput.Summary)),
			},
		}
	}

	issue, err := client.CreateCustomIssue(handlerInput.Project, handlerInput.IssueType, handlerInput.Summary,
		handlerInput.Description, handlerInput.Labels)
	if err != nil {
		err = fmt.Errorf("could not create issue: %v", err)
		log.Println(err)
		return newCreateIssueFailureEvent(err.Error(), handlerInput.Project, handlerInput.IssueType, handlerInput.Summary)
	}

	return flyte.Event{
		EventDef: createIncIssueEventDef,
		Payload: CreateIncIssueSuccess{
			ID:   issue.ID,
			Key:  issue.Key,
			Self: issue.Self,
		},
	}
}

var createIssueEventDef = flyte.EventDef{
	Name: "CreateIssue",
}

type createIssueSuccessPayload struct {
	Id          string `json:"id"`
	Url         string `json:"url"`
	Project     string `json:"project"`
	IssueType   string `json:"issuetype"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Reporter    string `json:"reporter"`
}

var createIssueFailureEventDef = flyte.EventDef{
	Name: "CreateIssueFailure",
}

type createIssueFailurePayload struct {
	Error       string `json:"error"`
	Project     string `json:"project"`
	IssueType   string `json:"issuetype"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

func newCreateIssueFailureEvent(err, project, issueType, summary string) flyte.Event {
	return flyte.Event{
		EventDef: createIssueFailureEventDef,
		Payload: createIssueFailurePayload{
			Error:     err,
			Project:   project,
			IssueType: issueType,
			Summary:   summary,
		},
	}
}

func newCreateIssueEvent(url, id, project, issueType, summary string, description string, priority string, reporter string) flyte.Event {
	return flyte.Event{
		EventDef: createIssueEventDef,
		Payload: createIssueSuccessPayload{
			Id:          id,
			Url:         url,
			Project:     project,
			IssueType:   issueType,
			Summary:     summary,
			Description: description,
			Priority:    priority,
			Reporter:    reporter,
		},
	}
}

// Types and functions for createIncIssue command
type CreateIncIssueSuccess struct { // to be returned into slack
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type CreateIncIssueFailure struct { // to be returned into slack
	Message string `json:"message"`
}

var createIncIssueEventDef = flyte.EventDef{
	Name: "CreateIncIssue",
}

var createIncIssueFailureEventDef = flyte.EventDef{
	Name: "CreateIncIssueFailure",
}

// incidentPattern checks if provided ServiceNow incident string matches the pattern
// INC1234567 or INC12345678. Returns bool
func incidentPattern(s string) bool {
	var p = regexp.MustCompile(`^INC[0-9]{7,8}$`)
	if p.MatchString(s) {
		return true
	} else {
		return false
	}
}

func updateSummaryForReactionBasedRequest(s *Input) {
	lenoftitle := 0
	if s.Summary == "ReactionBased" {
		s.Summary = s.Description
		if len(s.Description) > 70 {
			lenoftitle = 70
		} else {
			lenoftitle = len(s.Description)
		}
		title := s.Description[0:lenoftitle]
		title = strings.ReplaceAll(title, "\n", "")
		s.Summary = title + "..."
	}
}
