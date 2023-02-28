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
	"net/http"
	"net/url"
	"strconv"
)

type ServiceAccountService struct {
	sling *sling.Sling
}

// NewServiceAccountService Create new ServiceAccountService
func NewServiceAccountService(sling *sling.Sling) *ServiceAccountService {
	return &ServiceAccountService{
		sling: sling.Path("/keyhub/rest/v1/serviceaccount/"),
	}
}

// GetByUUID Get Service account by UUID
func (s *ServiceAccountService) GetByUUID(uuid uuid.UUID) (result *model.ServiceAccount, err error) {
	list := new(model.ServiceAccountList)
	errorReport := new(model.ErrorReport)

	params := &model.AccountQueryParams{UUID: uuid.String()}

	_, err = s.sling.New().Get("").QueryStruct(params).Receive(list, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get ServiceAccount %q.", uuid)
	}
	if err == nil {
		if len(list.Items) > 0 {
			result = &list.Items[0]
		} else {
			err = fmt.Errorf("Account %q not found", uuid.String())
		}
	}

	return
}

// GetById Get Service account by ID
func (s *ServiceAccountService) GetById(id int64) (result *model.ServiceAccount, err error) {
	sa := new(model.ServiceAccount)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Get(idString).Receive(sa, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get ServiceAccount %q.", idString)
		return
	}
	if err == nil && sa == nil {
		err = fmt.Errorf("Account %q not found", idString)
		return
	}

	//use an intermediate variable so sling can fill that variable with the json results. When request was succesful we use the variable as return value.
	result = sa
	return
}

// List all Service Accounts
func (s *ServiceAccountService) List() (list []model.ServiceAccount, err error) {
	list = []model.ServiceAccount{}

	searchRange := model.NewRange()

	var response *http.Response

	for ok := true; ok; ok = searchRange.NextPage() {

		errorReport := new(model.ErrorReport)
		results := new(model.ServiceAccountList)
		response, err = s.sling.New().Get("").Add(searchRange.GetRequestRangeHeader()).Add(searchRange.GetRequestModeHeader()).Receive(results, errorReport)
		searchRange.ParseResponse(response)

		if errorReport.Code > 0 {
			err = errorReport.Wrap("Could not get Groups.")
		}
		if err == nil {
			if len(results.Items) > 0 {
				list = append(list, results.Items...)
			}
		}

	}

	return
}

// Create  Create a serviceaccount
func (s *ServiceAccountService) Create(serviceaccount *model.ServiceAccount) (result *model.ServiceAccount, err error) {
	serviceAccounts := new(model.ServiceAccountList)
	results := new(model.ServiceAccountList)
	errorReport := new(model.ErrorReport)
	serviceAccounts.Items = append(serviceAccounts.Items, *serviceaccount)

	_, err = s.sling.New().Post("").BodyJSON(serviceAccounts).Receive(results, errorReport)

	if errorReport.Code > 0 {
		fmt.Println("Wrap", errorReport.Message)
		//apiErr := model.NewKeyhubApiError(*errorReport, "Could not create ServiceAccount in System %q.", serviceaccount.System.Name)
		err = errorReport.Wrap("Could not create ServiceAccount in System %q.", serviceaccount.System.Name)

		return
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Could not create ServiceAccount")
		}
	}

	return
}

// Update  Update service account
func (s *ServiceAccountService) Update(serviceAccount *model.ServiceAccount) (result *model.ServiceAccount, err error) {
	updated := new(model.ServiceAccount)
	errorReport := new(model.ErrorReport)

	selfUrl, _ := url.Parse(serviceAccount.Self().Href)

	_, err = s.sling.New().Path(selfUrl.Path).Put("").BodyJSON(serviceAccount).Receive(updated, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not update ServiceAccount %s", serviceAccount.UUID)
		return
	}
	result = updated
	return
}

// Delete Delete a service account by object
func (s *ServiceAccountService) Delete(serviceAccount *model.ServiceAccount) (err error) {
	errorReport := new(model.ErrorReport)

	selfUrl, _ := url.Parse(serviceAccount.Self().Href)

	_, err = s.sling.New().Path(selfUrl.Path).Delete("").Receive(nil, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not delete ServiceAccount %q", serviceAccount.UUID)
	}

	return
}

// DeleteByUUID  Delete a service account by uuid for a certain group
func (s *ServiceAccountService) DeleteByUUID(uuid uuid.UUID) (err error) {
	serviceAccount, err := s.GetByUUID(uuid)
	if err != nil {
		return err
	}

	return s.Delete(serviceAccount)
}

// DeleteByID  Delete a service account by ID
func (s *ServiceAccountService) DeleteByID(id int64) (err error) {
	serviceAccount, err := s.GetById(id)
	if err != nil {
		return err
	}
	return s.Delete(serviceAccount)
}
