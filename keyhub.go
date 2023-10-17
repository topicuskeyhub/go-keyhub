/* Licensed to the Apache Software Foundation (ASF) under one or more
   contributor license agreements.  See the NOTICE file distributed with
   this work for additional information regarding copyright ownership.
   The ASF licenses this file to You under the Apache License, Version 2.0
   (the "License"); you may not use this file except in compliance with
   the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License. */

package keyhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/coreos/go-oidc"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dghubble/sling"
)

const (
	/* KeyHub contract version supported by this client, set to 0 for latest */
	supportedContractVersion = 62

	/* KeyHub json mediatype */
	mediatype = "application/vnd.topicus.keyhub+json"
)

type Client struct {
	ID                 string
	Version            *VersionService
	Accounts           *AccountService
	Groups             *GroupService
	Systems            *SystemService
	ClientApplications *ClientApplicationService
	Vaults             *VaultService
	LaunchPadTile      *LaunchPadTileService
}

// khJsonBodyProvider encodes a JSON tagged struct value as a Body for requests.
// See https://golang.org/pkg/encoding/json/#MarshalIndent for details.
type khJsonBodyProvider struct {
	payload interface{}
}

func (p khJsonBodyProvider) ContentType() string {
	return ""
}

func (p khJsonBodyProvider) Body() (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(p.payload)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func NewClientDefault(issuer string, clientID string, clientSecret string) (*Client, error) {
	http.DefaultClient.Transport = http.DefaultTransport

	if http.DefaultClient.Timeout == 0 {
		http.DefaultClient.Timeout = time.Duration(time.Second * 10)
	}

	return NewClient(http.DefaultClient, issuer, clientID, clientSecret)
}

func NewClient(httpClient *http.Client, issuer string, clientID string, clientSecret string) (*Client, error) {

	var header_version string

	base := sling.New().Client(httpClient).Base(issuer)

	versionService := newVersionService(base.New().Set("Accept", "application/json").Set("Content-Type", "application/json"))
	version, err := versionService.Get()
	if err != nil {
		return nil, err
	}

	if supportedContractVersion > 0 {

		isContractVersionSupported := false
		for _, contractVersion := range version.ContractVersions {
			if supportedContractVersion == contractVersion {
				isContractVersionSupported = true
				break
			}
		}
		if !isContractVersionSupported {
			return nil, fmt.Errorf("KeyHub %v does not support api contract version %v", version.KeyhubVersion, supportedContractVersion)
		}

		header_version = strconv.FormatInt(supportedContractVersion, 10)
	} else {
		header_version = "latest"
	}

	versionedSling := base.New().
		Set("Accept", fmt.Sprintf("%v;version=%s", mediatype, header_version)).
		Set("Content-Type", fmt.Sprintf("%v;version=%s", mediatype, header_version))

	ctx := oidc.ClientContext(context.Background(), httpClient)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	var appClientConf = clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID},
		TokenURL:     provider.Endpoint().TokenURL + "?authVault=access",
	}
	oauth2Client := appClientConf.Client(ctx)
	oauth2Client.Timeout = httpClient.Timeout

	oauth2Sling := versionedSling.New().Client(oauth2Client)

	vaultClient := &http.Client{
		Transport: &Transport{
			Base: oauth2Client.Transport,
		},
	}
	return &Client{
		ID:                 clientID,
		Version:            versionService,
		Accounts:           newAccountService(oauth2Sling.New()),
		ClientApplications: newClientApplicationService(oauth2Sling.New()),
		Groups:             newGroupService(oauth2Sling.New()),
		Systems:            newSystemService(oauth2Sling.New()),
		LaunchPadTile:      newLaunchPadTileService(oauth2Sling.New()),
		Vaults:             newVaultService(versionedSling.New().Client(vaultClient)),
	}, nil
}
