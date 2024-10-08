package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// Jira credentials and API endpoint
var jiraURL = "https://[change].atlassian.net/rest/api/3/issue"
var email = "[change]"
var apiToken = os.Getenv("JIRA_AUTH_KEY")

// Jira project key where tickets should be created
var projectKey = "Atlas"

// Custom field ID for Epic Link
var epicLinkField = "customfield_10008" // Replace with your Jira instance's Epic Link field ID

// Ticket represents a Jira ticket
type Ticket struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IssueType   string `json:"issueType"`
	EpicKey     string `json:"epicKey"`
}

// CreateJiraTicket creates a Jira ticket
func CreateJiraTicket(title, description, issueType, epicKey string) {
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{
				"key": projectKey,
			},
			"summary":     title,
			"description": description,
			"issuetype": map[string]string{
				"name": issueType,
			},
			epicLinkField: epicKey,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", jiraURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(email, apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("Ticket '%s' created successfully.\n", title)
	} else {
		fmt.Printf("Failed to create ticket '%s'. Status Code: %d\n", title, resp.StatusCode)
	}
}

func main() {
	csvFile := flag.String("csv", "", "Path to the CSV file containing tickets")
	flag.Parse()

	if *csvFile == "" {
		fmt.Println("Please provide a CSV file using the -csv flag.")
		return
	}

	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read CSV file: %v\n", err)
		return
	}

	// Skip header
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			fmt.Println("Invalid record in CSV file. Each record must have a summary, description, issue type, and epic key.")
			continue
		}
		title := record[0]
		description := record[1]
		issueType := record[2]
		epicKey := record[3]
		CreateJiraTicket(title, description, issueType, epicKey)
	}
}
