package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	resty "github.com/go-resty/resty/v2"
)

// Issue struct for keeping issue data
type Issue struct {
	Key     string `json:"key"`
	Updated string `json:"updated"`
}

// Result struct for keeping Jira search results
type Result struct {
	// Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

func main() {
	fmt.Println("Run queries to Jira account...")
	client := resty.New()
	userEmail := os.Getenv("JIRA_USER_EMAIL")
	jiraToken := os.Getenv("JIRA_TOKEN")
	jiraAccount := os.Getenv("JIRA_ACCOUNT")

	client.SetBasicAuth(userEmail, jiraToken)

	count := 0
	total := 0
	maxResults := 15
	jql := `status = \"Merged / Done\" AND issuetype in (Bug, \"Change Request\", Story, Task, Sub-task) AND project in (DT) AND updated < -60d ORDER BY updated Asc`

	data, err := getResults(client, jiraAccount, getPayload(count, maxResults, jql), &count)

	if err == nil {
		count = data.MaxResults
		total = data.Total

		fmt.Println(count)
		fmt.Println(total)

		for count < total {
			getResults(client, jiraAccount, getPayload(count, maxResults, jql), &count)
		}
	} else {
		fmt.Println(err)
	}

	// Explore response object
}

func getPayload(startAt int, maxResults int, jql string) string {

	payload := `{
		"expand": [
		  "names",
		  "schema",
		  "operations"
		],
		"jql": "` + jql + `",
		"maxResults": ` + strconv.Itoa(maxResults) + `,
		"fieldsByKeys": false,
		"fields": [
			"updated"
		],
		"startAt":` + strconv.Itoa(startAt) + `
	  }`
	return payload
}

func getResults(client *resty.Client, jiraAccount string, payload string, count *int) (*Result, error) {
	data := &Result{
		Issues: []Issue{},
	}

	resp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(payload).
		// SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		Post("https://" + jiraAccount + "/rest/api/3/search")

	if err == nil {
		marsErr := json.Unmarshal(resp.Body(), data)
		printResponseOrError(resp, marsErr)
		*count += data.MaxResults

		fmt.Println("Count ", *count)
		return data, nil
	}

	fmt.Println(err)
	return nil, err
}

func printResponseOrError(resp *resty.Response, err error) {
	if err == nil {
		printResponse(resp)
	} else {
		fmt.Println("Error      :", err)
	}
}
func printResponse(resp *resty.Response) {
	fmt.Println("Post Response Info:")
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Time       :", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())

	data := &Result{
		Issues: []Issue{},
	}
	err := json.Unmarshal(resp.Body(), data)
	printBodyOrError(data, err)
}

func printBodyOrError(data *Result, err error) {
	if err == nil {
		fmt.Println("Body       :", data.Issues)
	} else {
		fmt.Println(err)
	}
}
