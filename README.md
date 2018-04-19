## Overview

![Build Status](https://travis-ci.org/HotelsDotCom/flyte-jira.svg?branch=master)
+[![Docker Stars](https://img.shields.io/docker/stars/hotelsdotcom/flyte-jira.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-jira)
+[![Docker Pulls](https://img.shields.io/docker/pulls/hotelsdotcom/flyte-jira.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-jira)

The Jira pack provides the ability to create tickets, comment on tickets and
to get info about tickets.


## Build & Run
### Command Line
To build and run from the command line:
* Clone this repo
* Run `dep ensure` (must have [dep](https://github.com/golang/dep) installed)
* Run `go build`
* Run `FLYTE_API_URL=http://.../ JIRA_HOST=https://... JIRA_USER=... JIRA_PASSWORD=... ./flye-jira`
* Fill in this command with the relevant API url, jira host, jira user and jira password environment variables

### Docker
To build and run from docker
* Run `docker build -t flye-jira .`
* Run `docker run -e FLYTE_API_URL=http://.../ -e JIRA_HOST=https://... -e JIRA_USER=... -e JIRA_PASSWORD=... flye-jira`
* All of these environment variables need to be set

## Commands
This pack provides three commands: `CommentTicket`, `TicketInfo` and `CreateTicket`.
### ticketInfo command
This command returns information about a specific ticket.
#### Input
This commands input is the id of the desired ticket:
```
"input" : "TEST-123",
```
#### Output
This command can either return an `Info` event or an `InfoFailure` event.
##### Info event
This is the success event, it contains the Id, Summary, Status, Description and Assignee. It returns them in the form:
```
"payload": {
    "id": "TEST-123",
    "summary": "Fix client race condition",
    "status": "In Progress",
    "description": "The client experiences.....",
    "assignee": "jsmith",
}
```
##### InfoFailure event
This contains the id of the ticket and the error.
```
"payload": {
    "id" : "TEST-123",
    "error": "Could not get info on TEST-123: status code 400",
}
```

### CreateTicket command
This command creates a jira ticket.
#### Input
This commands input is the project the ticket should be created under, the issue type and the title.
```
"input": {
    "project": "TEST",
    "issue_type": "Story",
    "title": "Fix csetcd bug"
    }
```
#### Output
This command can return either a `CreateTicket` event or a `CreateTicketFailure` event. 
##### CreatedTicket event
This is the success event, it contains the id of the ticket and the url of the ticket along with the input(project, 
issue_type & title) It returns them in the form:
```
"payload": {
    "id": "TEST-123",
    "url": "https://localhost:8100/browse/TEST-123",
    "project": "TEST",
    "issue_type": "Story",
    "title": "Fix csetcd bug"
}
```
##### CreateTicketFailure event
This contains the error if the ticket cannot be created along with the input (project, issue_type & title):
```
"payload": {
    "error": "Cannot create ticket: Fix csetcd bug: status code 400",
    "project": "TEST",
    "issue_type": "Story",
    "title": "Fix csetcd bug"
}
```

### CommentTicket command
This command comments on an issue.
#### Input
This commands input is the id of the ticket and the comment to be added. 
```
"input": {
    "id": "TEST-123",
    "comment": "Added to backlog"
    }
```
#### Output
This command can return either a `Comment` event or a `CommentFailure` event. 
##### Comment event
This is the success event, it contains the id of the ticket, the comment and the status. 
Status is the status code returned when a comment is left successfully.
```
"payload": {
    "id": "TEST-123",
    "comment": "Added to backlog"
}
```
##### CommentFailure event
This returns the error if a ticket cannot be commented on successfully. It contains, the ticket id, the comment and
the error:
```
"payload": {
    "id": "TEST-123",
    "comment": "Added to backlog",
    "error": "Could not comment on ticket: status code 400"
}
```
