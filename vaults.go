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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/sling"
	"github.com/topicuskeyhub/go-keyhub/model"
)

type VaultService struct {
	sling *sling.Sling
}

func newVaultService(sling *sling.Sling) *VaultService {
	return &VaultService{
		sling: sling,
	}
}

// GetRecords Retrieve all vault records for a group (secrets are not included)
func (s *VaultService) GetRecords(g *model.Group) (result []model.VaultRecord, err error) {
	result, err = s.List(g)
	return
}

// GetRecords Retrieve all vault records for a group (secrets are not included)
func (s *VaultService) List(g *model.Group) (records []model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	url, _ := url.Parse(g.Self().Href)
	additional := []string{}
	additional = append(additional, "audit")
	params := &model.VaultRecordQueryParams{Additional: additional}

	_, err = s.sling.New().Path(url.Path+"/vault/").Get("record").QueryStruct(params).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Vault of Group '" + g.UUID + "'. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			records = results.Items
		} else {
			records = []model.VaultRecord{}
		}
	}

	return
}

// GetRecord Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) GetRecord(group *model.Group, uuid string, options model.RecordOptions) (record *model.VaultRecord, err error) {
	url, _ := url.Parse(group.Self().Href)
	record = new(model.VaultRecord)

	additional := []string{}
	if options.Audit {
		additional = append(additional, "audit")
	}
	if options.Secret {
		additional = append(additional, "secret")
	}

	params := &model.VaultRecordQueryParams{UUID: uuid, Additional: additional}
	sl := s.sling.New().Set("Range", "items=0-0").Path(url.Path + "/").Path("vault/record").QueryStruct(params)

	vi := &model.VaultRecordList{}
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
