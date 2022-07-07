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
	"strconv"

	"github.com/dghubble/sling"
	"github.com/google/uuid"
	"github.com/topicuskeyhub/go-keyhub/model"
)

// ClientService Service to manage client application configurations in keyhub
type ClientService struct {
	sling *sling.Sling
}

func newClientService(sling *sling.Sling) *ClientService {
	return &ClientService{
		sling: sling.Path("/keyhub/rest/v1/client/"),
	}
}

// Create a new client application in Keyhub
func (s *ClientService) Create(client *model.Client) (result *model.Client, err error) {
	clients := new(model.ClientList)
	results := new(model.ClientList)
	errorReport := new(model.ErrorReport)
	clients.Items = append(clients.Items, *client)

	_, err = s.sling.New().Post("").BodyJSON(clients).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not create Client. Error: %s", errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("Created Client not found")
		}
	}

	return
}

// List all available clients.
func (s *ClientService) List() (clients []model.Client, err error) {
	results := new(model.ClientList)
	errorReport := new(model.ErrorReport)

	_, err = s.sling.New().Get("").Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get Clients. Error: %s", errorReport.Message)
	}
	if err == nil {
		if len(results.Items) > 0 {
			clients = results.Items
		} else {
			clients = []model.Client{}
		}
	}

	return
}

// GetByUUID Retrieve a client by uuid
func (s *ClientService) GetByUUID(uuid uuid.UUID) (result *model.Client, err error) {
	al := new(model.ClientList)
	errorReport := new(model.ErrorReport)

	params := &model.ClientQueryParams{UUID: uuid.String()}
	params.Additional = []string{"secret", "audit"}
	_, err = s.sling.New().Get("").QueryStruct(params).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get Client %q. Error: %s", uuid, errorReport.Message)
	}
	if err == nil {
		if len(al.Items) > 0 {
			result = &al.Items[0]
		} else {
			err = fmt.Errorf("Client %q not found", uuid.String())
		}
	}

	return
}

// GetById Retrieve a client by keyhub id
func (s *ClientService) GetById(id int64) (result *model.Client, err error) {
	al := new(model.Client)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Get(idString).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get Client %q. Error: %s", idString, errorReport.Message)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("Client %q not found", idString)
		return
	}

	return al, nil
}
