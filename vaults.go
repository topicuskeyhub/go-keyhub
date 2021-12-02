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

func (s *VaultService) Create(group *model.Group, vaultRecord *model.VaultRecord) (result *model.VaultRecord, err error) {
	vaultRecords := new(model.VaultRecordList)
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)
	vaultRecords.Items = append(vaultRecords.Items, *vaultRecord)

	url, _ := url.Parse(group.Self().Href)
	additional := []string{}
	additional = append(additional, "secret")
	params := &model.VaultRecordQueryParams{Additional: additional}

	_, err = s.sling.New().Path(url.Path+"/vault/").Post("record").QueryStruct(params).BodyJSON(vaultRecords).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not create VaultRecord in Group '" + group.UUID + "'. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = errors.New("Created VaultRecord not found!")
		}
	}

	return
}

// GetRecords Retrieve all vault records for a group including audit (secrets are not included)
func (s *VaultService) GetRecords(g *model.Group) (result []model.VaultRecord, err error) {
	result, err = s.List(g, model.RecordOptions{})
	return
}

// List Retrieve all vault records for a group (secrets are not included)
func (s *VaultService) List(group *model.Group, options model.RecordOptions) (records []model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	url, _ := url.Parse(group.Self().Href)
	additional := []string{}
	if options.Audit {
		additional = append(additional, "audit")
	}
	params := &model.VaultRecordQueryParams{Additional: additional}
	_, err = s.sling.New().Path(url.Path+"/vault/").Get("record").QueryStruct(params).Receive(results, errorReport)

	if errorReport.Code > 0 {
		err = errors.New("Could not get VaultRecords of Group '" + group.UUID + "'. Error: " + errorReport.Message)
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
func (s *VaultService) GetRecord(group *model.Group, uuid string, options model.RecordOptions) (result *model.VaultRecord, err error) {
	result, err = s.Get(group, uuid, options)
	return
}

//  Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) Get(group *model.Group, uuid string, options model.RecordOptions) (result *model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)
	url, _ := url.Parse(group.Self().Href)

	additional := []string{}
	if options.Audit {
		additional = append(additional, "audit")
	}
	if options.Secret {
		additional = append(additional, "secret")
	}
	params := &model.VaultRecordQueryParams{UUID: uuid, Additional: additional}

	_, err = s.sling.New().Path(url.Path+"/vault/").Get("record").QueryStruct(params).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get VaultRecord '" + uuid + "' of Group '" + group.UUID + "'. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = errors.New("VaultRecord '" + uuid + "' of Group '" + group.UUID + "' not found!")
		}
	}

	return
}

func (s *VaultService) Delete(group *model.Group, uuid string) (err error) {
	errorReport := new(model.ErrorReport)

	vaultRecord, err := s.GetRecord(group, uuid, model.RecordOptions{})
	if err != nil {
		return err
	}

	url, _ := url.Parse(vaultRecord.Self().Href)

	_, err = s.sling.New().Path(url.Path).Delete("").Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not delete VaultRecord '" + uuid + "' of Group '" + group.UUID + "'. Error: " + errorReport.Message)
	}

	return
}
