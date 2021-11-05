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

	"github.com/dghubble/sling"
	"github.com/topicuskeyhub/go-keyhub/model"
)

type AccountService struct {
	sling *sling.Sling
}

func newAccountService(sling *sling.Sling) *AccountService {
	return &AccountService{
		sling: sling.Path("/keyhub/rest/v1/account/"),
	}
}

func (s *AccountService) List() (accounts []model.Account, err error) {
	results := new(model.AccountList)
	errorReport := new(model.ErrorReport)

	_, err = s.sling.New().Get("").Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Accounts. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			accounts = results.Items
		} else {
			accounts = []model.Account{}
		}
	}

	return
}

func (s *AccountService) Get(uuid string) (a *model.Account, err error) {
	al := new(model.AccountList)
	errorReport := new(model.ErrorReport)

	params := &model.AccountQueryParams{UUID: uuid}

	_, err = s.sling.New().Get("").QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errors.New("Could not get Account '" + uuid + "'. Error: " + errorReport.Message)
	}
	if err == nil {
		if len(al.Items) > 0 {
			a = &al.Items[0]
		} else {
			err = errors.New("Account '" + uuid + "' not found!")
		}
	}

	return
}
