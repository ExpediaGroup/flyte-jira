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
	"github.com/ExpediaGroup/flyte-jira/domain"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
	"net/url"
	"path"
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
	}

	infoFailurePayload struct {
		Id    string `json:"id"`
		Error string `json:"error"`
	}
)

func infoHandler(input json.RawMessage) flyte.Event {
	issueId, err := parseInput(input)
	if err != nil {
		log.Printf("Error parsing input for IssueInfo: %s", err)
		return newInfoFailureEvent(err.Error(), issueId)
	}

	issue, err := client.GetIssueInfo(issueId)
	if err != nil {
		err = fmt.Errorf("could not get info: %v", err)
		log.Print(err)
		return newInfoFailureEvent(err.Error(), issueId)
	}

	return newInfoEvent(issue)
}

func newInfoFailureEvent(err, id string) flyte.Event {
	return flyte.Event{
		EventDef: infoFailureEventDef,
		Payload: infoFailurePayload{
			id,
			err},
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

func parseInput(input json.RawMessage) (string, error) {
	log.Printf("Input raw: %s", input)
	var in string
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	// Slack places URLs between <> tags, if using it as the input method
	// then the tags need to be stripped first before processing the rest
	if strings.Contains(in, "<") {
		in = in[1:len(in)-1]
	}

	if !hasURLFormat(in) {
		log.Printf("Input is not url: %s", in)
		return in, nil
	}

	id := path.Base(in)
	if id == "." || id == "\\" {
		return "", fmt.Errorf("url format not supported for issueId: %s", in)
	}

	log.Printf("URL Base: %s", id)
	return id, nil
}

// Credit: https://stackoverflow.com/a/55551215
func hasURLFormat(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != "" && u.Path != ""
}
