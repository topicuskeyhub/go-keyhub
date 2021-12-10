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
	"net/url"
	"strconv"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
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
		err = fmt.Errorf("Could not create VaultRecord in Group %q. Error: %s", group.UUID, errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Created VaultRecord not found")
		}
	}

	return
}

// GetRecords Retrieve all vault records for a group including audit (secrets are not included)
func (s *VaultService) GetRecords(g *model.Group) (result []model.VaultRecord, err error) {
	result, err = s.List(g, nil, nil)
	return
}

// List Retrieve all vault records for a group (secrets are not included, default audit = true)
func (s *VaultService) List(group *model.Group, query *model.VaultRecordQueryParams, additional *model.VaultRecordAdditionalQueryParams) (records []model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	url, _ := url.Parse(group.Self().Href)

	if query == nil {
		query = &model.VaultRecordQueryParams{}
	}
	if additional == nil {
		additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}

	additionalParams := []string{}
	if additional.Audit {
		additionalParams = append(additionalParams, "audit")
	}
	query.Additional = additionalParams

	_, err = s.sling.New().Path(url.Path+"/vault/").Get("record").QueryStruct(query).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get VaultRecords of Group %q. Error: %s", group.UUID, errorReport.Message)
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
func (s *VaultService) GetRecord(group *model.Group, uuid uuid.UUID, options model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	result, err = s.GetByUUID(group, uuid, &options)
	return
}

//  Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) GetByUUID(group *model.Group, uuid uuid.UUID, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	url, _ := url.Parse(group.Self().Href)

	query := &model.VaultRecordQueryParams{UUID: uuid.String()}
	if additional == nil {
		additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}

	additionalParams := []string{}
	if additional.Audit {
		additionalParams = append(additionalParams, "audit")
	}
	if additional.Secret {
		additionalParams = append(additionalParams, "secret")
	}
	query.Additional = additionalParams

	_, err = s.sling.New().Path(url.Path+"/vault/").Get("record").QueryStruct(query).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get VaultRecord %q of Group %q. Error: %s", uuid.String(), group.UUID, errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("VaultRecord %q of Group %q not found", uuid.String(), group.UUID)
		}
	}

	return
}

//  Retrieve a vault record by ID for a certain group, including audit and secrets
func (s *VaultService) GetByID(group *model.Group, id int64, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	al := new(model.VaultRecord)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	url, _ := url.Parse(group.Self().Href)

	query := &model.VaultRecordQueryParams{}
	if additional == nil {
		additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}

	additionalParams := []string{}
	if additional.Audit {
		additionalParams = append(additionalParams, "audit")
	}
	if additional.Secret {
		additionalParams = append(additionalParams, "secret")
	}
	query.Additional = additionalParams

	_, err = s.sling.New().Path(url.Path+"/vault/record/").Get(idString).QueryStruct(query).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get VaultRecord %q of Group %q. Error: %s", idString, group.UUID, errorReport.Message)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("VaultRecord %q of Group %q not found", idString, group.UUID)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return
}

//  Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) Update(group *model.Group, vaultRecord *model.VaultRecord) (result *model.VaultRecord, err error) {
	al := new(model.VaultRecord)
	errorReport := new(model.ErrorReport)

	url, _ := url.Parse(vaultRecord.Self().Href)

	query := &model.VaultRecordQueryParams{}
	additionalParams := []string{}
	additionalParams = append(additionalParams, "audit")
	additionalParams = append(additionalParams, "secret")
	query.Additional = additionalParams

	if vaultRecord.AdditionalObjects.Audit != nil {
		vaultRecord.AdditionalObjects.Audit = nil
	}

	_, err = s.sling.New().Path(url.Path).Put("").BodyJSON(vaultRecord).QueryStruct(query).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not update VaultRecord %q of Group %q. Error: %s", vaultRecord.UUID, group.UUID, errorReport.Message)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return
}

//  Delete a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) DeleteByUUID(group *model.Group, uuid uuid.UUID) (err error) {
	errorReport := new(model.ErrorReport)

	vaultRecord, err := s.GetByUUID(group, uuid, nil)
	if err != nil {
		return err
	}

	url, _ := url.Parse(vaultRecord.Self().Href)

	_, err = s.sling.New().Path(url.Path).Delete("").Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not delete VaultRecord %q of Group %q. Error: %s", uuid.String(), group.UUID, errorReport.Message)
	}

	return
}

//  Delete a vault record by ID for a certain group, including audit and secrets
func (s *VaultService) DeleteByID(group *model.Group, id int64) (err error) {
	errorReport := new(model.ErrorReport)
	url, _ := url.Parse(group.Self().Href)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Path(url.Path+"/vault/record/").Delete(idString).Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not delete VaultRecord %q of Group %q. Error: %s", idString, group.UUID, errorReport.Message)
	}

	return
}
