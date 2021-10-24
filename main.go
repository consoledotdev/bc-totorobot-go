// Copyright 2021, Console Ltd https://console.dev
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hanzoai/gochimp3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	listID = "267911a165" // https://mailchimp.com/help/find-audience-id/
)

// Trace ID is used to track a request through the function calls
// It's set in the HTTP handler, then unset once the request completes
var (
	traceID string = ""
)

// Entry defines a log entry in Google Cloud logging format
// https://github.com/GoogleCloudPlatform/golang-samples/blob/fa7b610d56d1d8b7d2002ecc30d995f7e3874de9/run/logging-manual/main.go
type Entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
// https://github.com/GoogleCloudPlatform/golang-samples/blob/fa7b610d56d1d8b7d2002ecc30d995f7e3874de9/run/logging-manual/main.go
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}

func main() {
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func main",
		Component: "main",
		Trace:     traceID,
	})

	// Define HTTP server.
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/getMailchimpStats", getMailchimpStatsHandler)

	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println(Entry{
		Severity:  "NOTICE",
		Message:   fmt.Sprintf("Starting server on port %s", port),
		Component: "main",
		Trace:     traceID,
	})

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

// HANDLERS

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Set global trace ID for use in other function calls
	if traceID == "" {
		traceID = getTraceID(r)
	}

	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func index handler",
		Component: "indexHandler",
		Trace:     traceID,
	})

	// The / path matches everything that is not defined above
	// So if the path ins't /, 404
	if r.URL.Path != "/" {
		log.Println(Entry{
			Severity:  "NOTICE",
			Message:   fmt.Sprintf("Unknown path: %s", r.URL.Path),
			Component: "indexHandler",
			Trace:     getTraceID(r),
		})

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 - Not Found")
		return
	}
	fmt.Fprintf(w, "indexHandler")
	traceID = "" // Unset now the request has finished
}

func getMailchimpStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Set global trace ID for use in other function calls
	if traceID == "" {
		traceID = getTraceID(r)
	}

	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func getMailchimpStatsHandler",
		Component: "getMailchimpStatsHandler",
		Trace:     getTraceID(r),
	})

	memberCount := getMailchimpListMemberCount(listID)
	// https://us7.admin.mailchimp.com/lists/segments?id=518946
	confirmedCount := getMailchimpListSegmentMemberCount(listID, "3577267")
	unconfirmedCount := getMailchimpListSegmentMemberCount(listID, "3577271")

	// Construct Basecamp message
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "Construct Basecamp message",
		Component: "getMailchimpStatsHandler",
		Trace:     getTraceID(r),
	})
	var content strings.Builder
	p := message.NewPrinter(language.English)
	content.WriteString("<strong>Mailchimp Stats (go)</strong><ul>")
	content.WriteString(p.Sprintf("<li><strong>Confirmed subscribers:</strong> %d</li>", confirmedCount))
	content.WriteString(p.Sprintf("<li><strong>Unconfirmed members:</strong> %d</li>", unconfirmedCount))
	content.WriteString(p.Sprintf("<li><strong>Total list members:</strong> %d</li>", memberCount))
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   content.String(),
		Component: "getMailchimpStatsHandler",
		Trace:     getTraceID(r),
	})

	// Post to Basecamp
	postBasecampChat(content.String())

	fmt.Fprintf(w, "OK")

	traceID = "" // Unset now the request has finished
}

// INTERNAL METHODS

// Gets Google Cloud trace ID
// https://github.com/GoogleCloudPlatform/golang-samples/blob/fa7b610d56d1d8b7d2002ecc30d995f7e3874de9/run/logging-manual/main.go
func getTraceID(r *http.Request) string {
	// Derive the traceID associated with the current request.
	var trace string
	traceHeader := r.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")

	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("K_SERVICE"), traceParts[0])
	}

	return trace
}

// Accesses the payload for the given secret version. The version can be a
// version number as a string (e.g. "5") or an alias (e.g. "latest").
// E.g. `accessSecret("my-secret/versions/5")`
func accessSecret(name string) (string, error) {
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func accessSecret",
		Component: "accessSecret",
		Trace:     traceID,
	})

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Println(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("failed to create secretmanager client: %v", err),
			Component: "accessSecret",
			Trace:     traceID,
		})
		return "", err
	}
	defer client.Close()

	name = "projects/bc-totorobot-go/secrets/" + name

	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   fmt.Sprintf("Requesting secret %s", name),
		Component: "accessSecret",
		Trace:     traceID,
	})

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Println(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("failed to access secret version: %v", err),
			Component: "accessSecret",
			Trace:     traceID,
		})
		return "", err
	}

	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "Secret returned",
		Component: "accessSecret",
		Trace:     traceID,
	})

	secret := string(result.Payload.Data)
	//log.Printf("Plaintext: %s\n", secret)
	return secret, nil
}

func getMailchimpListMemberCount(listID string) (MemberCount int) {
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func getMailchimpListMemberCount",
		Component: "getMailchimpListMemberCount",
		Trace:     traceID,
	})

	// Get Mailchimp API key
	apiKey, err := accessSecret("mailchimp-api-key/versions/latest")

	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Failed to get secret: %v", err),
			Component: "getMailchimpListMemberCount",
			Trace:     traceID,
		})
	}

	client := gochimp3.New(apiKey)

	// Fetch list
	log.Println(Entry{
		Severity:  "INFO",
		Message:   fmt.Sprintf("Get list: %s", listID),
		Component: "getMailchimpListMemberCount",
		Trace:     traceID,
	})

	list, err := client.GetList(listID, nil)
	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Failed to get list: %v", err),
			Component: "getMailchimpListMemberCount",
			Trace:     traceID,
		})
	}

	// Get list info
	// https://mailchimp.com/developer/api/marketing/lists/get-list-info/
	stats := list.Stats

	log.Println(Entry{
		Severity:  "INFO",
		Message:   fmt.Sprintf("Member count %d", stats.MemberCount),
		Component: "getMailchimpListMemberCount",
		Trace:     traceID,
	})

	return stats.MemberCount
}

func getMailchimpListSegmentMemberCount(listID string, SegmentID string) (MemberCount int) {
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func getMailchimpListSegmentMemberCount",
		Component: "getMailchimpListSegmentMemberCount",
		Trace:     traceID,
	})

	// Get Mailchimp API key
	apiKey, err := accessSecret("mailchimp-api-key/versions/latest")
	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Failed to get secret: %v", err),
			Component: "getMailchimpListMemberCount",
			Trace:     traceID,
		})
	}

	client := gochimp3.New(apiKey)

	// Fetch list
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   fmt.Sprintf("Get list %s", listID),
		Component: "getMailchimpListSegmentMemberCount",
		Trace:     traceID,
	})

	list, err := client.GetList(listID, nil)
	if err != nil {
		log.Fatalf("Failed to get list: %s", err)
	}

	// Get Segment info
	// https://mailchimp.com/developer/marketing/api/list-segments/get-segment-info/
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   fmt.Sprintf("Get segment %s", SegmentID),
		Component: "getMailchimpListSegmentMemberCount",
		Trace:     traceID,
	})

	segment, err := list.GetSegment(SegmentID, nil)
	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Failed to get segment: %s", err),
			Component: "getMailchimpListSegmentMemberCount",
			Trace:     traceID,
		})
	}

	log.Println(Entry{
		Severity:  "INFO",
		Message:   fmt.Sprintf("Member count %d", segment.MemberCount),
		Component: "getMailchimpListSegmentMemberCount",
		Trace:     traceID,
	})

	return segment.MemberCount
}

func postBasecampChat(content string) {
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "func postBasecampChat",
		Component: "postBasecampChat",
		Trace:     traceID,
	})

	// Create JSON payload
	postBody, _ := json.Marshal(map[string]string{
		"content": content,
	})

	// Create HTTP POST request
	log.Println(Entry{
		Severity:  "DEBUG",
		Message:   "Create HTTP POST request",
		Component: "postBasecampChat",
		Trace:     traceID,
	})

	// Get Basecamp Chatbot URL
	basecampChatbotURL, err := accessSecret("basecamp-chatbot-url/versions/latest")
	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Failed to get secret: %v", err),
			Component: "postBasecampChat",
			Trace:     traceID,
		})
	}

	requestBody := bytes.NewBuffer(postBody)
	req, err := http.Post(basecampChatbotURL, "application/json", requestBody)

	if err != nil {
		log.Fatalln(Entry{
			Severity:  "CRITICAL",
			Message:   fmt.Sprintf("Error with HTTP POST: %s", err),
			Component: "postBasecampChat",
			Trace:     traceID,
		})
	}

	defer req.Body.Close() // Close connection on function return

	log.Println(Entry{
		Severity:  "INFO",
		Message:   fmt.Sprintf("Response: %s", req.Status),
		Component: "postBasecampChat",
		Trace:     traceID,
	})
}
