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
	"strconv"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub/model"
)

type GroupService struct {
	sling *sling.Sling
}

func newGroupService(sling *sling.Sling) *GroupService {
	return &GroupService{
		sling: sling.Path("/keyhub/rest/v1/group/"),
	}
}

func (s *GroupService) Create(group *model.Group) (result *model.Group, err error) {
	groups := new(model.GroupList)
	results := new(model.GroupList)
	errorReport := new(model.ErrorReport)
	groups.Items = append(groups.Items, *group)

	_, err = s.sling.New().Post("").BodyJSON(groups).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not create Group. Error: %s", errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Created Group not found")
		}
	}

	return
}

func (s *GroupService) CreateMembership(group *model.Group, list *model.GroupAccountList) (results *model.GroupAccountList, err error) {

	idString := strconv.FormatInt(group.Self().ID, 10)

	errorReport := new(model.ErrorReport)

	_, err = s.sling.New().Post(idString+"/account").BodyJSON(list).Receive(results, errorReport)

	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not create memberschip. Error: %s", errorReport.Message)
	}
	if err == nil {
		if len(results.Items) == 0 {
			err = fmt.Errorf("Created memberships not returned")
		}
	}

	return
}

func (s *GroupService) List() (groups []model.Group, err error) {
	results := new(model.GroupList)
	errorReport := new(model.ErrorReport)
	groups = []model.Group{}
	searchRange := model.NewRange()

	var response *http.Response

	for ok := true; ok; ok = searchRange.NextPage() {

		response, err = s.sling.New().Get("").Add(searchRange.GetRequestRangeHeader()).Add(searchRange.GetRequestModeHeader()).Receive(results, errorReport)
		searchRange.ParseResponse(response)

		if errorReport.Code > 0 {
			err = fmt.Errorf("Could not get Groups. Error: %s", errorReport.Message)
		}
		if err == nil {
			if len(results.Items) > 0 {
				groups = append(groups, results.Items...)
			}
		}

	}

	return
}

func (s *GroupService) GetByUUID(uuid uuid.UUID) (result *model.Group, err error) {
	results := new(model.GroupList)
	errorReport := new(model.ErrorReport)

	params := &model.GroupQueryParams{
		UUID:       uuid.String(),
		Additional: &model.GroupAdditionalQueryParams{Admins: true},
	}

	_, err = s.sling.New().Get("").QueryStruct(params).Receive(results, errorReport)

	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get Group %q. Error: %s", uuid.String(), errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Group %q not found", uuid.String())
		}
	}

	return
}

func (s *GroupService) GetById(id int64) (result *model.Group, err error) {
	al := new(model.Group)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	params := &model.GroupQueryParams{
		Additional: &model.GroupAdditionalQueryParams{Admins: true},
	}

	_, err = s.sling.New().Get(idString).QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get Group %q. Error: %s", idString, errorReport.Message)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("Group %q not found", idString)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return
}
