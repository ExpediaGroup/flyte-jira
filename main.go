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

package main

import (
	jira "github.com/ExpediaGroup/flyte-jira/client"
	"github.com/ExpediaGroup/flyte-jira/command"
	"github.com/HotelsDotCom/flyte-client/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	jira.JiraConfig = initializeConfig()

	hostUrl := getUrl(getEnv("FLYTE_API_URL"))

	packDef := flyte.PackDef{
		Name:    "Jira",
		HelpURL: getUrl("https://github.com/ExpediaGroup/flyte-jira/blob/master/README.md"),
		Commands: []flyte.Command{
			command.IssueInfoCommand,
			command.CreateIssueCommand,
			command.IssueCommentCommand,
			command.SearchIssuesCommand,
			command.IssueAssignCommand,
			command.IssueCreateLinkCommand,
			command.IssueGetLinkCommand,
			command.IssueDeleteLinkCommand,
		},
	}

	p := flyte.NewPack(packDef, client.NewClient(hostUrl, 10*time.Second))
	p.Start()

	select {}
}

func initializeConfig() jira.Config {
	return jira.Config{
		getEnv("JIRA_HOST"),
		getEnv("JIRA_USER"),
		getEnv("JIRA_PASSWORD"),
	}
}

func getEnv(env string) string {
	value := os.Getenv(env)
	if value == "" {
		log.Fatalf("%s env. variable is not set", env)
	}
	return value
}

func getUrl(urlString string) *url.URL {
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("%s is not a valid url", urlString)
	}
	return u
}
