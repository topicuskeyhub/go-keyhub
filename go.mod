module github.com/topicuskeyhub/go-keyhub

go 1.21

retract (
	v1.3.1
	v1.3.2
	v1.3.3
	v1.3.4
)


require (
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dghubble/sling v1.4.0
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.14.0
	golang.org/x/net v0.0.0-20220531201128-c960675eff93
	golang.org/x/oauth2 v0.0.0-20220524215830-622c5d57e401
)

require github.com/gosimple/unidecode v1.0.1 // indirect

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0
	github.com/jarcoal/httpmock v1.2.0
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
