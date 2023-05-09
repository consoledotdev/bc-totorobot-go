module bc-totorobot-go

go 1.17

// Fix struct defintion for GetSegment ID
// https://github.com/hanzoai/gochimp3/pull/59
replace github.com/hanzoai/gochimp3 => github.com/davidmytton/gochimp3 v0.0.0-20211022095840-cabfd0e08b5e

require (
	cloud.google.com/go/secretmanager v1.10.1
	github.com/hanzoai/gochimp3 v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.9.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
)

require (
	cloud.google.com/go/compute v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.13.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/s2a-go v0.1.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/api v0.118.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require github.com/jarcoal/httpmock v1.3.0
