module bc-totorobot-go

go 1.17

// Fix struct defintion for GetSegment ID
// https://github.com/hanzoai/gochimp3/pull/59
replace github.com/hanzoai/gochimp3 => github.com/davidmytton/gochimp3 v0.0.0-20211022095840-cabfd0e08b5e

require (
	cloud.google.com/go/secretmanager v1.1.0
	github.com/hanzoai/gochimp3 v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.7
	google.golang.org/genproto v0.0.0-20220201184016-50beb8ab5c44
)

require (
	cloud.google.com/go v0.100.2 // indirect
	cloud.google.com/go/compute v0.1.0 // indirect
	cloud.google.com/go/iam v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	google.golang.org/api v0.66.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.40.1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

require github.com/jarcoal/httpmock v1.0.8
