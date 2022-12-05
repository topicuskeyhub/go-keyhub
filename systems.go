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
	"net/url"
	"strconv"
)

type SystemService struct {
	sling *sling.Sling
}

func newSystemService(sling *sling.Sling) *SystemService {
	return &SystemService{
		sling: sling.Path("/keyhub/rest/v1/system/"),
	}
}

func (s *SystemService) FindGroupOnSystem(system *model.ProvisionedSystem, query *model.GroupOnSystemQueryParams) (results *model.GroupOnSystemList, err error) {

	results = &model.GroupOnSystemList{}

	selfUrl, _ := url.Parse(system.Self().Href)

	if query == nil {
		query = &model.GroupOnSystemQueryParams{}
	}
	if query.Additional == nil {
		query.Additional = &model.GroupOnSystemAdditionalQueryParams{Audit: false}
	}

	searchRange := model.NewRange()

	for ok := true; ok; ok = searchRange.NextPage() {

		errorReport := new(model.ErrorReport)
		pageList := new(model.GroupOnSystemList)
		response, _ := s.sling.New().Path(selfUrl.Path+"/").Get("group").QueryStruct(query).
			Add(searchRange.GetRequestRangeHeader()).
			Add(searchRange.GetRequestModeHeader()).
			Receive(pageList, errorReport)

		searchRange.ParseResponse(response)

		if errorReport.Code > 0 {
			err = errorReport.Wrap("could not get GroupsOnSystem for System %s.", system.UUID)
			return nil, err
		}

		results.Items = append(results.Items, pageList.Items...)

	}

	return

}

func (s *SystemService) CreateGroupOnSystem(groupOnSystem *model.GroupOnSystem) (result *model.GroupOnSystem, err error) {

	list := new(model.GroupOnSystemList)
	results := new(model.GroupOnSystemList)
	errorReport := new(model.ErrorReport)
	groupId := strconv.FormatInt(groupOnSystem.System.Self().ID, 10)
	groupOnSystem.System = nil

	list.Items = append(list.Items, *groupOnSystem)

	_, err = s.sling.New().Post(groupId+"/group").BodyJSON(list).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not create GroupOnSystem.")
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Created GOS not found")
		}
	}
	return
}

func (s *SystemService) GetGroupOnSystem(system *model.ProvisionedSystem, groupId int64, additional *model.GroupOnSystemAdditionalQueryParams) (result *model.GroupOnSystem, err error) {

	al := new(model.GroupOnSystem)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(system.Self().ID, 10)
	groupIdString := strconv.FormatInt(groupId, 10)

	if additional == nil {
		additional = &model.GroupOnSystemAdditionalQueryParams{Audit: false, ProvGroups: true}
	}

	params := &model.GroupOnSystemQueryParams{
		Additional: additional,
	}
	_, err = s.sling.New().Get(idString+"/group/"+groupIdString).QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get GroupOnSystem \"%s/%s\".", idString, groupIdString)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("System %q not found", idString)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return

}

func (s *SystemService) DeleteGroupOnSystem(groupOnSystem *model.GroupOnSystem) (err error) {

	errorReport := new(model.ErrorReport)
	gosId := strconv.FormatInt(groupOnSystem.Self().ID, 10)
	groupId := strconv.FormatInt(groupOnSystem.System.Self().ID, 10)

	params := struct {
		System bool `url:"system"`
	}{
		System: true,
	}

	var result interface{}

	_, err = s.sling.New().Delete(groupId+"/group/"+gosId).QueryStruct(params).Receive(result, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("could not delete GroupOnSystem.")
	}
	return
}

/*
func (s *SystemService) Create(system *model.ProvisionedSystem) (result *model.ProvisionedSystem, err error) {

		return
	}

	func (s *SystemService) List() (systems []model.ProvisionedSystem, err error) {
		return
	}
*/
func (s *SystemService) GetByUUID(uuid uuid.UUID) (system *model.ProvisionedSystem, err error) {
	results := new(model.ProvisionedSystemList)
	errorReport := new(model.ErrorReport)

	params := &model.GroupQueryParams{
		UUID:       uuid.String(),
		Additional: &model.GroupAdditionalQueryParams{Admins: false},
	}

	_, err = s.sling.New().Get("").QueryStruct(params).Receive(results, errorReport)

	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get System %q.", uuid.String())
	}
	if err == nil {
		if len(results.Items) > 0 {
			system = &results.Items[0]
		} else {
			err = fmt.Errorf("System %q not found", uuid.String())
		}
	}

	return

}

func (s *SystemService) GetById(id int64) (system *model.ProvisionedSystem, err error) {
	al := new(model.ProvisionedSystem)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	params := &model.GroupQueryParams{
		Additional: &model.GroupAdditionalQueryParams{Admins: false},
	}
	_, err = s.sling.New().Get(idString).QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get Group %q.", idString)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("System %q not found", idString)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	system = al
	return

}
