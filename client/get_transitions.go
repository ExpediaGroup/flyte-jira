package client

import (
	"errors"
	"fmt"
	"net/http"
)

type TransitionsResult struct {
	Transitions []TransitionObj `json:"transitions"`
}
type TransitionObj struct {
	TransitionId   string `json:"id"`
	TransitionName string `json:"name"`
}

func GetTransitions(issueId string) (TransitionsResult, error) {
	result := TransitionsResult{Transitions: []TransitionObj{}}
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueId)
	request, err := constructGetRequest(path)
	if err != nil {
		return result, err
	}
	statusCode, err := SendRequest(request, &result)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		return result, errors.New("Issue does not exist")
	}
	return result, nil
}
