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
	"net/url"
	"time"

	"github.com/dghubble/sling"
)

type vaultItems struct {
	Items []VaultRecord `json:"items"`
}

type VaultRecord struct {
	UUID     string   `json:"uuid"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Username string   `json:"username"`
	Color    string   `json:"color"`
	Filename string   `json:"filename"`
	Types    []string `json:"types"`
	Links    []struct {
		ID   int    `json:"id"`
		Href string `json:"href"`
	} `json:"links"`
	AdditionalObjects struct {
		Audit struct {
			CreatedAt      time.Time `json:"createdAt"`
			CreatedBy      string    `json:"createdBy"`
			LastModifiedAt time.Time `json:"lastModifiedAt"`
			LastModifiedBy string    `json:"lastModifiedBy"`
		} `json:"audit"`
		Secret struct {
			Password string `json:"password"`
			File     []byte `json:"file"`
			Comment  string `json:"comment"`
		} `json:"secret"`
	} `json:"additionalObjects"`
}

func (r *VaultRecord) CreatedAt() time.Time {
	return r.AdditionalObjects.Audit.CreatedAt
}

func (r *VaultRecord) CreatedBy() string {
	return r.AdditionalObjects.Audit.CreatedBy
}

func (r *VaultRecord) LastModifiedAt() time.Time {
	return r.AdditionalObjects.Audit.LastModifiedAt
}

func (r *VaultRecord) LastModifiedBy() string {
	return r.AdditionalObjects.Audit.LastModifiedBy
}

func (r *VaultRecord) Comment() string {
	return r.AdditionalObjects.Secret.Comment
}

func (r *VaultRecord) Password() string {
	return r.AdditionalObjects.Secret.Password
}

func (r *VaultRecord) File() []byte {
	return r.AdditionalObjects.Secret.File
}

type VaultService struct {
	sling *sling.Sling
}

func newVaultService(sling *sling.Sling) *VaultService {
	return &VaultService{
		sling: sling,
	}
}

// GetRecords Retrieve all vault records for a group (secrets are not included)
func (s *VaultService) GetRecords(g *Group) ([]VaultRecord, error) {
	url, _ := url.Parse(g.Links[0].Href)

	additional := []string{}
	additional = append(additional, "audit")
	params := &vaultParams{Additional: additional}
	vaultRecords := &vaultItems{}
	resp, err := s.sling.New().Path(url.Path + "/").Path("vault/").Get("record").QueryStruct(params).ReceiveSuccess(vaultRecords)
	if err != nil {
		return nil, err
	}
	if code := resp.StatusCode; code < 200 || code > 299 {
		return nil, fmt.Errorf("Vault query returned with statuscode %d", code)
	}

	return vaultRecords.Items, nil
}

type vaultParams struct {
	UUID       string   `url:"uuid,omitempty"`
	Additional []string `url:"additional,omitempty"`
	Any        bool     `url:"any,omitempty"`
}

type RecordOptions struct {
	Audit  bool
	Secret bool
}

// GetRecord Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) GetRecord(group *Group, uuid string, options RecordOptions) (record *VaultRecord, err error) {
	url, _ := url.Parse(group.Links[0].Href)
	record = new(VaultRecord)
	additional := []string{}
	if options.Audit {
		additional = append(additional, "audit")
	}
	if options.Secret {
		additional = append(additional, "secret")
	}

	params := &vaultParams{UUID: uuid, Additional: additional}
	sl := s.sling.New().Set("Range", "items=0-0").Path(url.Path + "/").Path("vault/record").QueryStruct(params)

	vi := &vaultItems{}
	_, err = sl.ReceiveSuccess(vi)
	if err != nil {
		return
	}

	if len(vi.Items) == 1 {
		record = &vi.Items[0]
	}

	return
}

func (s *VaultService) Decode(resp *http.Response, v interface{}) error {
	fmt.Println("decoding")
	buf := make([]byte, 20000)
	defer resp.Body.Close()
	n, _ := resp.Body.Read(buf)
	fmt.Println("Body:")
	fmt.Println(string(buf[:n]))
	return nil
}
