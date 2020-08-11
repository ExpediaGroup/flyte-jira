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

var IssueCommentCommand = flyte.Command{
	Name:         "CommentIssue",
	OutputEvents: []flyte.EventDef{commentEventDef, commentFailureEventDef},
	Handler:      commentHandler,
}

func commentHandler(input json.RawMessage) flyte.Event {
	var handlerInput struct {
		Id      string `json:"id"`
		Comment string `json:"comment"`
	}

	if err := json.Unmarshal(input, &handlerInput); err != nil {
		err = fmt.Errorf("Could not marshal comment into json: %s", err)
		log.Println(err)
		return newCommentFailureEvent(err.Error(), "unknown", "unkown")
	}

	_, err := client.CommentIssue(handlerInput.Id, handlerInput.Comment)
	if err != nil {
		err = fmt.Errorf("Could not leave comment: %s", err)
		log.Println(err)
		return newCommentFailureEvent(err.Error(), handlerInput.Id, handlerInput.Comment)
	}

	return newCommentEvent(handlerInput.Id, handlerInput.Comment)
}

var commentEventDef = flyte.EventDef{
	Name: "Comment",
}

type commentSuccessPayload struct {
	Id      string `json:"id"`
	Comment string `json:"comment"`
}

var commentFailureEventDef = flyte.EventDef{
	Name: "CommentFailure",
}

type commentFailurePayload struct {
	Id      string `json:"id"`
	Comment string `json:"comment"`
	Error   string `json:"error"`
}

func newCommentFailureEvent(error, id, comment string) flyte.Event {
	return flyte.Event{
		EventDef: commentFailureEventDef,
		Payload: commentFailurePayload{
			Id:      id,
			Comment: comment,
			Error:   error,
		},
	}
}

func newCommentEvent(id, comment string) flyte.Event {
	return flyte.Event{
		EventDef: commentEventDef,
		Payload: commentSuccessPayload{
			Id:      id,
			Comment: comment,
		},
	}
}
