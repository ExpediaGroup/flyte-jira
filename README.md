## Overview

[![Build Status](https://github.com/ExpediaGroup/flyte-jira/workflows/Build/badge.svg?branch=master&event=push)](https://github.com/ExpediaGroup/flyte-jira/actions?query=workflow:"Build"branch:"master")
[![Docker Stars](https://img.shields.io/docker/stars/expediagroup/flyte-jira.svg)](https://hub.docker.com/r/expediagroup/flyte-jira)
[![Docker Pulls](https://img.shields.io/docker/pulls/expediagroup/flyte-jira.svg)](https://hub.docker.com/r/expediagroup/flyte-jira)

The Jira pack provides the ability to create issues, comment on issues and
to get info about issues.


## Build & Run
### Command Line
To build and run from the command line:
* Clone this repo
* Run `go build`
* Run `FLYTE_API_URL=http://.../ JIRA_HOST=https://... JIRA_USER=... JIRA_PASSWORD=... ./flyte-jira`
* Fill in this command with the relevant API url, Jira host, Jira user and Jira password environment variables

### Docker
To build and run from docker
* Run `docker build -t flyte-jira .`
* Run `docker run -e FLYTE_API_URL=http://.../ -e JIRA_HOST=https://... -e JIRA_USER=... -e JIRA_PASSWORD=... flyte-jira`
* All of these environment variables need to be set

## Commands
This pack provides the following commands: `CommentIssue`, `IssueInfo`, `CreateIssue`, `IssueAssign`, `IssueCreateLink`, `IssueGetLink`, `IssueDeleteLink`
### issueInfo command
This command returns information about a specific issue.
#### Input
This command's input is the id of the desired issue. The issue can either be a URL with an issueId base or an issueId in string representation.
E.g
```
"input": "TEST-123",
"input": "http://foo.bar/TEST-123"
```

Because slack automatically places embeded URLs in between `< >` tags, then the following is also accepted:
`"input": "<http://foo.bar/TEST-123>"`
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
    "assignee": "jsmith@expediagroup.com",
    "reporter": "test@expediagroup.com",
    "priority": "Medium",
    "components": "Compute Platform",
    "labels": "feature-request",
    "type": "Support"
}
```
##### InfoFailure event
This contains the id of the issue and the error.
```
"payload": {
    "id" : "TEST-123",
    "error": "Could not get info on TEST-123: status code 400",
}
```

### CreateIssue command
This command creates a Jira issue.
#### Input
This command inputs are the project that the issue should be created under, the issue type and the summary (title).
```
"fields": {
    "project":  
       {
          "key": "TEST"
       },
    "issuetype": 
       {
          "name": "Story"
       },
    "summary": "Fix csetcd bug"
    }
```
#### Output
This command can return either a `CreateIssue` event or a `CreateIssueFailure` event.
##### CreateIssue event
This is the success event, it contains the id of the issue and the url of the issue along with the input (project,
issuetype & summary) It returns them in the form:
```
"payload": {
    "id": "TEST-123",
    "url": "https://localhost:8100/browse/TEST-123",
    "project": "TEST",
    "issuetype": "Story",
    "summary": "Fix csetcd bug",
    "priority": "Medium",
    "reporter": "songupta@expediagroup.com",
    "description": "This is issue description"
}
```
##### CreateIssueFailure event
This contains the error if the issue cannot be created along with the input (project, issuetype & summary):
```
"payload": {
    "error": "Cannot create issue: Fix csetcd bug: status code 400",
    "project": "TEST",
    "issuetype": "Story",
    "summary": "Fix csetcd bug"
}
```

### CreateIncIssue command
This command creates a Jira issue in a target project (specified in a flow).
Required a ServiceNow incident and labels (optional) as an argument 
#### Input
This command inputs are the project that the issue should be created under, the issue type and the summary (title).
```
"input": {
          "incident": "INC1234567",
          "issuetype": "Story",
          "project": "SRO",
          "summary": "Test",
          "description": "This JIRA issue was created by Flyte",
          "labels": [
            "TEST-LABEL"
          ]
        }
```
#### Output
This command can return either a `CreateIncIssue` event or a `CreateIncIssueFailure` event.
##### CreateIncIssue event
This is the success event. It returns them in the form:
```
"payload": {
    "id": "10000",
    "key": "FLYTE-1",
    "self": "https://jira.expedia.biz/rest/api/2/issue/10000",
}
```
##### CreateIncIssueFailure event
This contains the error message if the issue cannot be created:
```
"payload": {
    "message": "Cannot create issue: Fix csetcd bug: status code 400",
}
```

### CommentIssue command
This command comments on an issue.
#### Input
This commands input is the id of the issue and the comment to be added.
```
"input": {
    "id": "TEST-123",
    "comment": "Added to backlog"
    }
```
#### Output
This command can return either a `Comment` event or a `CommentFailure` event. 
##### Comment event
This is the success event, it contains the id of the issue, the comment and the status.
Status is the status code returned when a comment is left successfully.
```
"payload": {
    "id": "TEST-123",
    "comment": "Added to backlog"
}
```
##### CommentFailure event
This returns the error if a issue cannot be commented on successfully. It contains, the issue id, the comment and
the error:
```
"payload": {
    "id": "TEST-123",
    "comment": "Added to backlog",
    "error": "Could not comment on issue: status code 400"
}
```

### SearchIssues command
This command searches issues using [JQL queries](https://confluence.atlassian.com/jirasoftwareserver/advanced-searching-939938733.html).
#### Input
This command's inputs are the query string, the index of the first element to be retrieved in the list of results and the number of elements to be retrieved.
```
"input": {
    "query": "project = Flyte", // required
    "startIndex": 0,            // optional, default: 0
    "maxResults": 10            // optional, default: 10
}
```
#### Output
This command can return either a `SearchSuccess` event or a `SearchFailure` event. 
##### SearchSuccess event
This is the success event, it contains the values given as input for the command, the total number of possible results and the issues retrieved.
```
"payload": {
    "query": "project = Flyte",
    "startIndex": 0,
    "maxResults": 10,
    "total": 85,
    "issues":[
        {
            "id": "TEST-123",
            "summary": "Fix client race condition",
            "status": "In Progress",
            "description": "The client experiences.....",
            "assignee": "jsmith",
        }
        ... 
    ]
}
```
##### SearchFailure event
This returns the error if the jql query is not valid. It contains the values given as input for the command and the error:
```
"payload": {
    "query": "project = Flyte",
    "startIndex": 0,
    "maxResults": 10,
    "error": "Could not search for issues: statusCode=400"
}
```
---
[issue-assign]: https://docs.atlassian.com/software/jira/docs/api/REST/7.6.1/#api/2/issue-assign
### IssueAssign command
Assign a user to a JIRA issue

### Input
The input is a `json` object with a `username` and an `issueId` fields. According to the Jira [API docs][issue-assign], the username field can be left empty to unassign the issue, which is why it can be omitted from the input. Likewise, a "-1" string name will `auto-assign` the issue.
`input JSON object`:
```json
{
  "username": "test-123", // optional, nil name == unassign
  "issueId": "ISSUE-01" // required
}
```
#### Output
The output consists of two events: an `AssignEvent` or an `AssignFailureEvent`

#### AssignEvent
If the assignment is successful, an event will be propagated back with a payload consisting of the initial request paramters.
```json
"payload": {
  "username":"test-123",
  "issueId": "ISSUE-01"
}
```

#### AssignFailureEvent
If the assignment is unsuccessful, an `assignFailureEvent` will come back with a payload consisting of the initial request and the error message.
```json
"payload": {
  "username": "foo",
  "issueId": "ISSUE01",
  "error": "Unauthorised"
}
```

---
### Links

Jira offers the possibility to manage links between issues. The commands that are available to do so are `IssueGetLink`, `IssueDeleteLink` and `IssueCreateLink`. While their functionality is self-explanatory, the input/output varies depending on the operation.

The two event types in the case of links are:
1. `Link` --> propagated on success
2. `LinkFailure` --> propagated on failure

#### IssueGetLink

will return a link object that links two issues. 

#### Input:
The input is a `json` object with a `linkId` field, where the `linkId` is id of the respective link object. Generally, those ids can be obtained by looking at a card information (`IssueInfo` now supports link information!).
```json
{
 "linkId": "12351245"
}
```

#### Output:
In case of a success, the output will contain the link information:
```json
{
 "inwardIssue": "<issue-key>",
 "outwardIssue": "<issue-key>",
 "linkType": {
   "Name": "Depends"
  },
 "Comment": "Link related issues!" 
}
```

#### IssueCreateLink
Based on 2 project keys and and a link type, create a link with that type between the 2 issues.

#### Input:
```json
{
 "inwardIssue": "<issue-key>",
  "outwardIssue": "<issue-key>",
  "linkType": "<type>",
}
```

#### Output:
The output is the initial request or a failure event if unsuccessful.

#### IssueDeleteLink
Based on a `linkId`, the link between 2 issues is deleted.

#### Input:
```json
{
 "linkId": "12341234"
}
```

#### Output:
The initial request if successful or a failure event if unsuccessful.

---
