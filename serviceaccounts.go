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
	"github.com/topicuskeyhub/go-keyhub/model"
	"net/http"
)

type ServiceAccountService struct {
	sling *sling.Sling
}

func NewServiceAccountService(sling *sling.Sling) *ServiceAccountService {
	return &ServiceAccountService{
		sling: sling.Path("/keyhub/rest/v1/serviceaccount/"),
	}
}

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
