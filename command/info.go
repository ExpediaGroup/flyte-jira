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
	"github.com/ExpediaGroup/flyte-client/flyte"
	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/ExpediaGroup/flyte-jira/domain"
	"log"
	"regexp"
	"strings"
)

var (
	IssueInfoCommand = flyte.Command{
		Name:         "IssueInfo",
		OutputEvents: []flyte.EventDef{infoEventDef, infoFailureEventDef},
		Handler:      infoHandler,
	}

	infoEventDef = flyte.EventDef{
		Name: "Info",
	}

	infoFailureEventDef = flyte.EventDef{
		Name: "InfoFailure",
	}
)

type (
	infoSuccessPayload struct {
		Id          string `json:"id"`
		Summary     string `json:"summary"`
		Status      string `json:"status"`
		Description string `json:"description"`
		Assignee    string `json:"assignee"`
		Reporter    string `json:"reporter"`
		Components  string `json:"components"`
		Labels      string `json:"labels"`
		Priority    string `json:"priority"`
		Type        string `json:"type"`
	}

	infoFailurePayload struct {
		Id    string `json:"id"`
		Error string `json:"error"`
	}
)

func infoHandler(input json.RawMessage) flyte.Event {
	var in string
	if err := json.Unmarshal(input, &in); err != nil {
		log.Printf("Error unmarshaling input for IssueInfo: %s", err)
		return newInfoFailureEvent("", err)
	}

	//`\w+-\d+ should resolve any regex of type <KEY-NUMBER>
	re := regexp.MustCompile(`\w+-\d+`)
	issueId := re.FindString(in)

	issue, err := client.GetIssueInfo(issueId)
	if err != nil {
		log.Printf("Error fetching IssueInfo for %s: %s", issueId, err)
		return newInfoFailureEvent(issueId, err)
	}

	return newInfoEvent(issue)
}

func newInfoFailureEvent(issueId string, err error) flyte.Event {
	return flyte.Event{
		EventDef: infoFailureEventDef,
		Payload:  infoFailurePayload{issueId, err.Error()},
	}
}

func newInfoEvent(t domain.Issue) flyte.Event {
	components := t.Fields.Components
	component_str := ""
	labels_str := ""

	labels := t.Fields.Labels
	for i := range components {
		component_str = component_str + components[i].Name
		component_str = component_str + ","
	}

	for i := range labels {
		labels_str = labels_str + labels[i]
		labels_str = labels_str + ","
	}

	return flyte.Event{
		EventDef: infoEventDef,
		Payload: infoSuccessPayload{
			Id:          t.Key,
			Summary:     t.Fields.Summary,
			Status:      t.Fields.Status.Name,
			Description: t.Fields.Description,
			Assignee:    t.Fields.Assignee.EmailAddress,
			Reporter:    t.Fields.Reporter.EmailAddress,
			Components:  strings.TrimSuffix(component_str, ","),
			Labels:      strings.TrimSuffix(labels_str, ","),
			Priority:    t.Fields.Priority.Name,
			Type:        t.Fields.Type.Name,
		},
	}
}
