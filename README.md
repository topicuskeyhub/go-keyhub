![Build](https://github.com/topicuskeyhub/go-keyhub/workflows/Build/badge.svg?branch=master)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/topicuskeyhub/go-keyhub?label=Release)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# go-keyhub - Topicus KeyHub API Client

### How to use
go.mod:
```
require (
  github.com/topicuskeyhub/go-keyhub v0.2.0
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
