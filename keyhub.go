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
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dghubble/sling"
)

const (
	/* KeyHub contract version supported by this client */
	supportedContractVersion = 40

	/* KeyHub json mediatype */
	mediatype = "application/vnd.topicus.keyhub+json"
)

type Client struct {
	ID      string
	Version *VersionService
	Groups  *GroupService
	Vaults  *VaultService
}

func NewClient(httpClient *http.Client, issuer string, clientID string, clientSecret string) (*Client, error) {
	if httpClient.Timeout == 0 {
		httpClient.Timeout = time.Duration(time.Second * 10)
	}

	base := sling.New().Base(issuer)

	versionService := newVersionService(base.New().Client(httpClient).Set("Accept", mediatype))
	version, err := versionService.Get()
	if err != nil {
		return nil, err
	}
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

	versionedSling := base.New().Set("Accept", fmt.Sprintf("%v;version=%v", mediatype, supportedContractVersion))

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
	oauth2Client.Timeout = time.Duration(time.Second * 10)

	oauth2Sling := versionedSling.New().Client(oauth2Client)

	vaultClient := &http.Client{
		Transport: &Transport{
			Base: oauth2Client.Transport,
		},
	}

	return &Client{
		ID:      clientID,
		Version: versionService,
		Groups:  newGroupService(oauth2Sling.New()),
		Vaults:  newVaultService(versionedSling.New().Client(vaultClient)),
	}, nil
}
