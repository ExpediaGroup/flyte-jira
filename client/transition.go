package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type transition struct {
	Id string `json:"id"`
}

type transitionRequest struct {
	Transition transition `json:"transition"`
}

func Transition(issueId, transitionId string) (string, error) {
	path := fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueId)
	r := transitionRequest{
		Transition: transition{Id: transitionId},
	}
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	req, err := constructPostRequest(path, string(b))
	if err != nil {
		return "", err
	}
	resp, err := SendRequestWithoutResp(req)
	if err != nil {
		return req.URL.Path, err
	}

	switch resp {
	case http.StatusNoContent:
		err = nil
	case http.StatusBadRequest:
		err = errors.New("no transition specified")
	case http.StatusUnauthorized:
		err = errors.New("invalid permission to transition an issue")
	case http.StatusNotFound:
		err = errors.New("issue or user does not exist")
	default:
		err = fmt.Errorf("unsupported status code %d", resp)
	}

	return req.URL.Path, err
}
