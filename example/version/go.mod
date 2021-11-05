module github.com/topicuskeyhub/go-keyhub/example-version

go 1.17

require github.com/topicuskeyhub/go-keyhub v0.2.1
require github.com/google/uuid v1.3.0

require (
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/dghubble/sling v1.4.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect

	github.com/pquerna/cachecontrol v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

replace github.com/topicuskeyhub/go-keyhub => ../..
