package command

import (
	"encoding/json"
	"fmt"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-jira/domain"
	"log"
)

var (
	searchSuccessEventDef = flyte.EventDef{Name: "SearchSuccess"}
	searchFailureEventDef = flyte.EventDef{Name: "SearchFailure"}
)

var SearchIssuesCommand = flyte.Command{
	Name:         "SearchIssues",
	OutputEvents: []flyte.EventDef{searchSuccessEventDef, searchFailureEventDef},
	Handler:      searchIssuesHandler,
}

func searchIssuesHandler(rawInput json.RawMessage) flyte.Event {

	input := SearchIssuesInput{"", 0, 10}
	if err := json.Unmarshal(rawInput, &input); err != nil {
		err := fmt.Errorf("input is not valid: %s", err)
		return flyte.NewFatalEvent(err)
	}

	if input.Query == "" {
		err := fmt.Errorf("Empty query string")
		return newSearchFailureEvent(input.Query, input.StartIndex, input.MaxResults, err.Error())
	}

	searchResult, err := client.SearchIssues(input.Query, input.StartIndex, input.MaxResults)
	if err != nil {
		err := fmt.Errorf("Could not search for issues: %v", err)
		log.Println(err)
		return newSearchFailureEvent(input.Query, input.StartIndex, input.MaxResults, err.Error())
	}

	return newSearchSuccessEvent(
		input.Query,
		input.StartIndex,
		input.MaxResults,
		searchResult.TotalResults,
		searchResult.Issues)
}

func newSearchSuccessEvent(query string, startIndex int, maxResults int, totalResults int, unformattedIssues []domain.Issue) flyte.Event {

	inputDetails := SearchIssuesInput{query, startIndex, maxResults}
	var issues []IssuePayload
	for _, issue := range unformattedIssues {
		formattedIssue := IssuePayload{
			Id:          issue.Key,
			Summary:     issue.Fields.Summary,
			Status:      issue.Fields.Status.Name,
			Description: issue.Fields.Description,
			Assignee:    issue.Fields.Assignee.Name,
		}
		issues = append(issues, formattedIssue)
	}
	return flyte.Event{
		EventDef: searchSuccessEventDef,
		Payload:  SearchSuccessOutput{inputDetails, totalResults, issues},
	}
}

func newSearchFailureEvent(query string, startIndex int, maxResults int, error string) flyte.Event {

	inputDetails := SearchIssuesInput{query, startIndex, maxResults}
	return flyte.Event{
		EventDef: searchFailureEventDef,
		Payload:  SearchFailureOutput{inputDetails, error},
	}
}

type SearchIssuesInput struct {
	Query      string `json:"query"`
	StartIndex int    `json:"startIndex"`
	MaxResults int    `json:"maxResults"`
}

type SearchSuccessOutput struct {
	SearchIssuesInput
	TotalResults int            `json:"total"`
	Issues       []IssuePayload `json:"issues"`
}

type SearchFailureOutput struct {
	SearchIssuesInput
	Error string `json:"error"`
}

type IssuePayload struct {
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
}
