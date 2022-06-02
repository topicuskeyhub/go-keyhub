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
	"net/http"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	keyhubmodel "github.com/topicuskeyhub/go-keyhub/model"
)

func init() {
	SetupTest()

	grouplist := keyhubmodel.GroupList{}
	sum := int64(1)
	for sum < 2 {
		sumString := strconv.FormatInt(sum, 10)
		gr := keyhubmodel.NewEmptyGroup("group " + sumString)
		gr.Links = append(gr.Links, keyhubmodel.Link{ID: sum, Rel: "self", Href: "https://topicus-keyhub.com/keyhub/rest/v1/group/" + sumString})
		grouplist.Items = append(grouplist.Items, *gr)
		sum += sum
	}

	httpmock.RegisterResponder("GET", `=~^https://topicus-keyhub.com/keyhub/rest/v1/group/(\d+)\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get ID from request
			id := httpmock.MustGetSubmatchAsUint(req, 1) // 1=first regexp submatch
			return httpmock.NewJsonResponse(200, grouplist.Items[id-1])
		})

	recordlist := keyhubmodel.VaultRecordList{}
	sum = int64(1)
	for sum < 10 {
		sumString := strconv.FormatInt(sum, 10)
		vr := keyhubmodel.NewVaultRecord("record "+strconv.FormatInt(sum, 10), &keyhubmodel.VaultRecordSecretAdditionalObject{Password: &sumString})
		vr.Links = append(vr.Links, keyhubmodel.Link{ID: sum, Rel: "self", Href: "https://topicus-keyhub.com/keyhub/rest/v1/group/1/vault/record/" + sumString})
		recordlist.Items = append(recordlist.Items, *vr)
		sum += sum
	}

	httpmock.RegisterResponder("GET", `=~^https://topicus-keyhub.com/keyhub/rest/v1/group/(\d+)/vault/record\z`, httpmock.NewJsonResponderOrPanic(206, recordlist))
	httpmock.RegisterResponder("GET", `=~^https://topicus-keyhub.com/keyhub/rest/v1/group/(\d+)/vault/record/(\d+)\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get ID from request
			// groupId := httpmock.MustGetSubmatchAsUint(req, 1) // 1=first regexp submatch
			recordId := httpmock.MustGetSubmatchAsUint(req, 2) // 2=second regexp submatch
			return httpmock.NewJsonResponse(200, recordlist.Items[recordId-1])
		})
}

func TestVaultRecords(t *testing.T) {

	client, err := NewClientDefault("https://topicus-keyhub.com", "clientid", "clientsecret")
	if err != nil {
		t.Fatalf("ERROR %s", err)
	}

	group, err := client.Groups.GetById(1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	records, err := client.Vaults.List(group, nil, &keyhubmodel.VaultRecordAdditionalQueryParams{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(records) == 0 {
		t.Fatalf("ERROR no records found")
	}

	record, err := client.Vaults.GetByID(group, 1, &keyhubmodel.VaultRecordAdditionalQueryParams{Secret: true})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
	if record == nil {
		t.Fatalf("ERROR record with id 1 not found")
	}
	if record.Self().ID != 1 {
		t.Fatalf("ERROR record with id 1 not found")
	}
}
