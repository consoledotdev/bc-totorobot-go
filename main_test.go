// Copyright 2021, Console Ltd https://console.dev
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	getListResponse = `{
		"id": "string",
		"web_id": 0,
		"name": "string",
		"contact": {
		  "company": "string",
		  "address1": "string",
		  "address2": "string",
		  "city": "string",
		  "state": "string",
		  "zip": "string",
		  "country": "string",
		  "phone": "string"
		},
		"permission_reminder": "string",
		"use_archive_bar": false,
		"campaign_defaults": {
		  "from_name": "string",
		  "from_email": "string",
		  "subject": "string",
		  "language": "string"
		},
		"date_created": "2019-08-24T14:15:22Z",
		"list_rating": 0,
		"email_type_option": true,
		"subscribe_url_short": "string",
		"subscribe_url_long": "string",
		"beamer_address": "string",
		"visibility": "pub",
		"double_optin": false,
		"has_welcome": false,
		"marketing_permissions": false,
		"modules": [
		  "string"
		],
		"stats": {
		  "member_count": 500,
		  "total_contacts": 0,
		  "unsubscribe_count": 0,
		  "cleaned_count": 0,
		  "member_count_since_send": 0,
		  "unsubscribe_count_since_send": 0,
		  "cleaned_count_since_send": 0,
		  "campaign_count": 0,
		  "campaign_last_sent": "2019-08-24T14:15:22Z",
		  "merge_field_count": 0,
		  "avg_sub_rate": 0,
		  "avg_unsub_rate": 0,
		  "target_sub_rate": 0,
		  "open_rate": 0,
		  "click_rate": 0,
		  "last_sub_date": "2019-08-24T14:15:22Z",
		  "last_unsub_date": "2019-08-24T14:15:22Z"
		},
		"_links": [
		  {
			"rel": "string",
			"href": "string",
			"method": "GET",
			"targetSchema": "string",
			"schema": "string"
		  }
		]
	  }`

	getSegmentResponse = `{
		"id": 0,
		"name": "string",
		"member_count": 200,
		"type": "saved",
		"created_at": "2019-08-24T14:15:22Z",
		"updated_at": "2019-08-24T14:15:22Z",
		"options": {
		"match": "any",
		"conditions": [
			null
		]
		},
		"list_id": "string",
		"_links": [
		{
			"rel": "string",
			"href": "string",
			"method": "GET",
			"targetSchema": "string",
			"schema": "string"
		}
		]
		}`
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "indexHandler"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIndexHandler404(t *testing.T) {
	req, err := http.NewRequest("GET", "/doesnotexist", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "404 - Not Found"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMailchimpStatsHandler(t *testing.T) {
	httpmock.Activate()                 // Blocks all HTTP requests
	defer httpmock.DeactivateAndReset() // Reset on exit

	// TODO: Mock Google Cloud Secrets Manager APIs
	// We don't actually want to block all HTTP requests (which is the default)
	// Because mocking Google APIs is a hassle, so we will allow those to return
	// This assumes you have the access to the secrets locally
	// See https://cloud.google.com/code/docs/vscode/secret-manager
	//
	// We'll mock Mailchimp and Basecamp
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	// Mock the API calls to Mailchimp
	// https://mailchimp.com/developer/marketing/api/lists/get-list-info/
	httpmock.RegisterResponder("GET", "https://us7.api.mailchimp.com/3.0/lists/testList1",
		httpmock.NewStringResponder(200, getListResponse))

	// https://mailchimp.com/developer/marketing/api/list-segments/get-segment-info/
	httpmock.RegisterResponder("GET", "https://us7.api.mailchimp.com/3.0/lists/string/segments/1",
		httpmock.NewStringResponder(200, getSegmentResponse))

	basecampChatbotURL, err := accessSecret("basecamp-chatbot-url/versions/latest")
	if err != nil {
		t.Errorf("Failed to get secret: %v", err)
	}

	// https://github.com/basecamp/bc3-api/blob/master/sections/chatbots.md#create-a-line
	httpmock.RegisterResponder("POST", basecampChatbotURL,
		httpmock.NewStringResponder(201, ""))

	req, err := http.NewRequest("GET", "/getMailchimpStats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getMailchimpStatsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMailchimpMethods(t *testing.T) {
	httpmock.Activate()                 // Blocks all HTTP requests
	defer httpmock.DeactivateAndReset() // Reset on exit

	// TODO: Mock Google Cloud Secrets Manager APIs
	// We don't actually want to block all HTTP requests (which is the default)
	// Because mocking Google APIs is a hassle, so we will allow those to return
	// This assumes you have the access to the secrets locally
	// See https://cloud.google.com/code/docs/vscode/secret-manager
	//
	// We'll just mock Mailchimp
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	// Mock the API calls to Mailchimp
	// https://mailchimp.com/developer/marketing/api/lists/get-list-info/
	httpmock.RegisterResponder("GET", "https://us7.api.mailchimp.com/3.0/lists/testList1",
		httpmock.NewStringResponder(200, getListResponse))

	memberCount := getMailchimpListMemberCount("testList1")
	log.Printf("Test: memberCount: %d", memberCount)

	if memberCount != 500 {
		t.Errorf("Test list count: got %q, want %q", memberCount, 500)
	}

	// https://mailchimp.com/developer/marketing/api/list-segments/get-segment-info/
	httpmock.RegisterResponder("GET", "https://us7.api.mailchimp.com/3.0/lists/string/segments/1",
		httpmock.NewStringResponder(200, getSegmentResponse))

	segmentCount := getMailchimpListSegmentMemberCount("testList1", "1")
	log.Printf("Test: segmentCount: %d", segmentCount)

	if segmentCount != 200 {
		t.Errorf("Test segment count: got %q, want %q", segmentCount, 200)
	}
}

func TestPostBasecampChat(t *testing.T) {
	httpmock.Activate()                 // Blocks all HTTP requests
	defer httpmock.DeactivateAndReset() // Reset on exit

	// TODO: Mock Google Cloud Secrets Manager APIs
	// We don't actually want to block all HTTP requests (which is the default)
	// Because mocking Google APIs is a hassle, so we will allow those to return
	// This assumes you have the access to the secrets locally
	// See https://cloud.google.com/code/docs/vscode/secret-manager
	//
	// We'll just mock Basecamp
	httpmock.RegisterNoResponder(httpmock.InitialTransport.RoundTrip)

	basecampChatbotURL, err := accessSecret("basecamp-chatbot-url/versions/latest")
	if err != nil {
		t.Errorf("Failed to get secret: %v", err)
	}

	// https://github.com/basecamp/bc3-api/blob/master/sections/chatbots.md#create-a-line
	httpmock.RegisterResponder("POST", basecampChatbotURL,
		httpmock.NewStringResponder(201, ""))

	postBasecampChat("Test")

	//info := httpmock.GetCallCountInfo()
	//t.Errorf("Tests: %+v", info)
}
