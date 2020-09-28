package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/ExpediaGroup/flyte-jira/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
)

func TestLinkCreateIsSuccessful(t *testing.T) {
	initialFunc := client.SendRequestWithoutResp
	defer func() { client.SendRequestWithoutResp = initialFunc }()
	in := []byte(`{
        "inwardIssue": "TEST-123",
        "outwardIssue": "Test-321",
        "linkType": "Depends"
        }`)

	expHttpReq := client.LinkIssueRequest{
		LinkType: client.IssueLink{Name: "Depends"},
		Inward:   client.LinkIssue{"TEST-123"},
		Outward:  client.LinkIssue{"Test-321"},
		Comment:  client.Comment{"Link related issues!"},
	}
	client.SendRequestWithoutResp = func(request *http.Request) (int, error) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return 400, err
		}

		body := client.LinkIssueRequest{}
		if err := json.Unmarshal(b, &body); err != nil {
			return 400, err
		}

		if !reflect.DeepEqual(expHttpReq, body) {
			return http.StatusNotFound, fmt.Errorf("expHttpRequest: %v actual: %v", expHttpReq, body)
		}

		return http.StatusOK, nil
	}
	actualEvent := issueCreateLinkHandler(in)
	expectedEvent := flyte.Event{
		EventDef: linkEventDef,
		Payload: linkRequest{
			InwardIssue:  "TEST-123",
			OutwardIssue: "Test-321",
			LinkType:     "Depends",
		},
	}

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}

func TestLinkGetIsSuccessful(t *testing.T) {
	initialFunc := client.SendRequest
	defer func() { client.SendRequest = initialFunc }()
	in := []byte(`{"linkId": "1223"}`)

	expLinkId := "1223"
	client.SendRequest = func(request *http.Request, respB interface{}) (responseCode int, err error) {
		reqPath := request.URL.Path
		linkId := path.Base(reqPath)
		if linkId != expLinkId {
			return http.StatusBadRequest, fmt.Errorf("expected issueId %s got %s", "DEVEX-553", linkId)
		}

		switch v := respB.(type) {
		case **client.LinkIssueRequest:
			{
				(*v).Inward = client.LinkIssue{"TEST-123"}
				(*v).Outward = client.LinkIssue{"TEST-321"}
				(*v).LinkType = client.IssueLink{Name: "Depends"}
			}
		}

		return http.StatusOK, nil
	}

	actualEvent := issueGetLinkHandler(in)
	expectedEvent := flyte.Event{
		EventDef: linkEventDef,
		Payload: &client.LinkIssueRequest{
			Inward:   client.LinkIssue{"TEST-123"},
			Outward:  client.LinkIssue{"TEST-321"},
			LinkType: client.IssueLink{Name: "Depends"},
		},
	}

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}

func TestLinkDeleteIsSuccessful(t *testing.T) {
	initialFunc := client.SendRequestWithoutResp
	defer func() { client.SendRequestWithoutResp = initialFunc }()
	in := []byte(`{"linkId": "1223"}`)

	expLinkId := "1223"
	client.SendRequestWithoutResp = func(request *http.Request) (int, error) {
		reqPath := request.URL.Path
		linkId := path.Base(reqPath)
		if linkId != expLinkId {
			return http.StatusBadRequest, fmt.Errorf("expected issueId %s got %s", "DEVEX-553", linkId)
		}

		return http.StatusOK, nil
	}

	actualEvent := issueDeleteLinkHandler(in)
	expectedEvent := flyte.Event{
		EventDef: linkEventDef,
		Payload: linkRequest{
			LinkId: "1223",
		},
	}

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Errorf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}
