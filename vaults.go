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
	"github.com/dghubble/sling"
	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub/model"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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

	selfUrl, _ := url.Parse(group.Self().Href)
	params := &model.VaultRecordQueryParams{
		Additional: &model.VaultRecordAdditionalQueryParams{Secret: true},
	}

	_, err = s.sling.New().Path(selfUrl.Path+"/vault/").Post("record").QueryStruct(params).BodyJSON(vaultRecords).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not create VaultRecord in Group %q.", group.UUID)
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

	selfUrl, _ := url.Parse(group.Self().Href)

	if query == nil {
		query = &model.VaultRecordQueryParams{}
	}
	if additional != nil {
		query.Additional = additional
	} else if query.Additional == nil {
		query.Additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}

	searchRange := model.NewRange()
	for ok := true; ok; ok = searchRange.NextPage() {

		errorReport := new(model.ErrorReport)
		results := new(model.VaultRecordList)
		var response *http.Response
		response, err = s.sling.New().Path(selfUrl.Path+"/vault/").Get("record").QueryStruct(query).Add(searchRange.GetRequestRangeHeader()).Add(searchRange.GetRequestModeHeader()).Receive(results, errorReport)
		searchRange.ParseResponse(response)

		if errorReport.Code > 0 {
			err = errorReport.Wrap("Could not get VaultRecords of Group %q.", group.UUID)
		}
		if err == nil {
			if len(results.Items) > 0 {
				records = append(records, results.Items...)
			}
		}

	}

	return
}

func (s *VaultService) getMyClientId() (id int64, err error) {

	me := new(model.ClientApplication)

	errorReport := new(model.ErrorReport)

	_, err = s.sling.New().Get("/keyhub/rest/v1/client/me").Receive(&me, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("could not determine client details")
		return
	}
	if err != nil {
		err = fmt.Errorf("could not determine client details, error: %s", err.Error())
		return
	}

	id = me.Self().ID

	return
}

func (s *VaultService) FindByIDForClient(id int64, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	query := model.VaultRecordSearchQueryParams{
		ID: strconv.FormatInt(id, 10),
	}

	return s.findForClient(query, additional)

}

func (s *VaultService) FindByUUIDForClient(uuid uuid.UUID, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {

	query := model.VaultRecordSearchQueryParams{
		UUID: uuid.String(),
	}

	return s.findForClient(query, additional)
}

func (s *VaultService) findForClient(query model.VaultRecordSearchQueryParams, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	clientID, err := s.getMyClientId()
	if err == nil {
		query.AccessibleByClient = strconv.FormatInt(clientID, 10)
	}

	additionalParams := []string{}
	// If secrets are requested we need to do a new request so no need for audit data in search results
	if !additional.Secret {
		if additional.Audit {
			additionalParams = append(additionalParams, "audit")
		}
	}
	query.Additional = additionalParams

	_, err = s.sling.New().Get("/keyhub/rest/v1/vaultrecord/").QueryStruct(query).Receive(results, errorReport)

	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not find VaultRecord.")
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]

			if additional.Secret {
				// if secrets are requested, we need to retrieve the record again from the group url.
				// If not we can simply return the search result

				r, _ := regexp.Compile("^((.+)/group/([0-9]+))/vault/record/([0-9]+)")
				matches := r.FindStringSubmatch(result.Self().Href)
				// 0 = full url (unused)
				// 1 = group url
				// 2 = base url (unused)
				// 3 = group id
				// 4 = record id

				// group id
				gid := big.Int{}
				gid.SetString(matches[3], 10)

				// record id
				rid := big.Int{}
				rid.SetString(matches[4], 10)

				// Build a fake group we can use for retreiving a record without another rest call .
				fakegroup := model.NewEmptyGroup("Unknown")

				fakegroup.Links = append(
					fakegroup.Links,
					model.Link{
						Href: matches[1],
						Rel:  "self",
						ID:   gid.Int64(),
					},
				)

				return s.GetByID(fakegroup, rid.Int64(), additional)
			} else {
				return result, err
			}

		} else {
			err = fmt.Errorf("no VaultRecords found")
		}
	}

	return
}

// GetByUUID Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) GetByUUID(group *model.Group, uuid uuid.UUID, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	results := new(model.VaultRecordList)
	errorReport := new(model.ErrorReport)

	selfUrl, _ := url.Parse(group.Self().Href)

	query := &model.VaultRecordQueryParams{UUID: uuid.String()}
	if additional == nil {
		additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}
	query.Additional = additional

	_, err = s.sling.New().Path(selfUrl.Path+"/vault/").Get("record").QueryStruct(query).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get VaultRecord %q of Group %q.", uuid.String(), group.UUID)
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

// GetByID  Retrieve a vault record by ID for a certain group, including audit and secrets
func (s *VaultService) GetByID(group *model.Group, id int64, additional *model.VaultRecordAdditionalQueryParams) (result *model.VaultRecord, err error) {
	al := new(model.VaultRecord)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	selfUrl, _ := url.Parse(group.Self().Href)

	query := &model.VaultRecordQueryParams{}
	if additional == nil {
		additional = &model.VaultRecordAdditionalQueryParams{Audit: true}
	}
	query.Additional = additional

	_, err = s.sling.New().Path(selfUrl.Path+"/vault/record/").Get(idString).QueryStruct(query).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get VaultRecord %q of Group %q.", idString, group.UUID)
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

// Update Retrieve a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) Update(group *model.Group, vaultRecord *model.VaultRecord) (result *model.VaultRecord, err error) {
	al := new(model.VaultRecord)
	errorReport := new(model.ErrorReport)

	selfUrl, _ := url.Parse(vaultRecord.Self().Href)

	query := &model.VaultRecordQueryParams{
		Additional: &model.VaultRecordAdditionalQueryParams{
			Audit:  true,
			Secret: true,
		},
	}

	if vaultRecord.AdditionalObjects.Audit != nil {
		vaultRecord.AdditionalObjects.Audit = nil
	}

	_, err = s.sling.New().Path(selfUrl.Path).Put("").BodyJSON(vaultRecord).QueryStruct(query).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not update VaultRecord %q of Group %q.", vaultRecord.UUID, group.UUID)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return
}

// DeleteByUUID  Delete a vault record by uuid for a certain group, including audit and secrets
func (s *VaultService) DeleteByUUID(group *model.Group, uuid uuid.UUID) (err error) {
	errorReport := new(model.ErrorReport)

	vaultRecord, err := s.GetByUUID(group, uuid, nil)
	if err != nil {
		return err
	}

	selfUrl, _ := url.Parse(vaultRecord.Self().Href)

	_, err = s.sling.New().Path(selfUrl.Path).Delete("").Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not delete VaultRecord %q of Group %q.", uuid.String(), group.UUID)
	}

	return
}

// DeleteByID  Delete a vault record by ID for a certain group, including audit and secrets
func (s *VaultService) DeleteByID(group *model.Group, id int64) (err error) {
	errorReport := new(model.ErrorReport)
	selfUrl, _ := url.Parse(group.Self().Href)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Path(selfUrl.Path+"/vault/record/").Delete(idString).Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not delete VaultRecord %q of Group %q.", idString, group.UUID)
	}

	return
}
