module bc-totorobot-go

go 1.17

// Fix struct defintion for GetSegment ID
// https://github.com/hanzoai/gochimp3/pull/59
replace github.com/hanzoai/gochimp3 => github.com/davidmytton/gochimp3 v0.0.0-20211022095840-cabfd0e08b5e

require (
	cloud.google.com/go/secretmanager v1.8.0
	github.com/hanzoai/gochimp3 v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.7
	google.golang.org/genproto v0.0.0-20221010155953-15ba04fc1c0e
)

require (
	cloud.google.com/go/compute v1.10.0 // indirect
	cloud.google.com/go/iam v0.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.6.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20221012135044-0b7e1fb9d458 // indirect
	golang.org/x/oauth2 v0.0.0-20221006150949-b44042a4b9c1 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	google.golang.org/api v0.99.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.50.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

require github.com/jarcoal/httpmock v1.2.0
