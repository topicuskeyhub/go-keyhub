![Build](https://github.com/topicuskeyhub/go-keyhub/workflows/Build/badge.svg?branch=master)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/topicuskeyhub/go-keyhub?label=Release)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# go-keyhub - Topicus KeyHub API Client

> [!CAUTION]
> DEPRECATED: please note that this package `github.com/topicuskeyhub/go-keyhub` has become deprecated in favor of `github.com/topicuskeyhub/sdk-go`
> 
> For more info see: https://github.com/topicuskeyhub/sdk-go/

### How to use

See the examples directory for more complete examples.

go.mod:
```
require (
  github.com/topicuskeyhub/go-keyhub v1.3.2
)
```

```go
import "github.com/topicuskeyhub/go-keyhub"

client, err := keyhub.NewClient(http.DefaultClient, issuer, clientid, clientsecret)
if err != nil {
    log.Fatalln("ERROR", err)
}

```

### How to develop
* Dependencies: `go mod download`
* Code formatting: `gofmt -s -w .`
