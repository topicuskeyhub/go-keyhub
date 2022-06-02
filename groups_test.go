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
	for sum < 10 {
		sumString := strconv.FormatInt(sum, 10)
		gr := keyhubmodel.NewEmptyGroup("group " + sumString)
		gr.Links = append(gr.Links, keyhubmodel.Link{ID: sum, Rel: "self", Href: "https://topicus-keyhub.com/keyhub/rest/v1/group/" + sumString})
		grouplist.Items = append(grouplist.Items, *gr)
		sum += sum
	}

	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/group/", httpmock.NewJsonResponderOrPanic(206, grouplist))
	httpmock.RegisterResponder("GET", `=~^https://topicus-keyhub.com/keyhub/rest/v1/group/(\d+)\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get ID from request
			id := httpmock.MustGetSubmatchAsUint(req, 1) // 1=first regexp submatch
			return httpmock.NewJsonResponse(200, grouplist.Items[id-1])
		})
}

func TestGroups(t *testing.T) {

	client, err := NewClientDefault("https://topicus-keyhub.com", "clientid", "clientsecret")
	if err != nil {
		t.Fatalf("ERROR %s", err)
	}

	groups, err := client.Groups.List()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(groups) == 0 {
		t.Fatalf("ERROR no groups found")
	}

	group, err := client.Groups.GetById(1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if group == nil {
		t.Fatalf("ERROR group with id 1 not found")
	}
	if group.Self().ID != 1 {
		t.Fatalf("ERROR group with id 1 not found")
	}
}
