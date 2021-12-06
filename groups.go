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
	"strconv"

	"github.com/dghubble/sling"
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
		err = errors.New("Could not create Group. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = errors.New("Created Group not found!")
		}
	}

	return
}

func (s *GroupService) List() (groups []model.Group, err error) {
	results := new(model.GroupList)
	errorReport := new(model.ErrorReport)

	_, err = s.sling.New().Get("").Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Groups. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			groups = results.Items
		} else {
			groups = []model.Group{}
		}
	}

	return
}

func (s *GroupService) Get(uuid string) (result *model.Group, err error) {
	result, err = s.GetByUUID(uuid)
	return
}

func (s *GroupService) GetByUUID(uuid string) (result *model.Group, err error) {
	results := new(model.GroupList)
	errorReport := new(model.ErrorReport)

	additional := []string{}
	additional = append(additional, "admins")
	params := &model.GroupQueryParams{UUID: uuid, Additional: additional}

	_, err = s.sling.New().Get("").QueryStruct(params).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Group '" + uuid + "'. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = errors.New("Group '" + uuid + "' not found!")
		}
	}

	return
}

func (s *GroupService) GetById(id int64) (result *model.Group, err error) {
	al := new(model.Group)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	additional := []string{}
	additional = append(additional, "admins")
	params := &model.GroupQueryParams{Additional: additional}

	_, err = s.sling.New().Get(idString).QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Group '" + idString + "'. Error: " + errorReport.Message)
		return
	}
	if err == nil && al == nil {
		err = errors.New("Group '" + idString + "' not found!")
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = al
	return
}
