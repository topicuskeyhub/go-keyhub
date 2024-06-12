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
	"github.com/google/go-querystring/query"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/topicuskeyhub/go-keyhub/model"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func init() {
	httpmock.Activate()

	versions := makeRange(57, 71)

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/info", httpmock.NewJsonResponderOrPanic(200, model.NewVersionInfo("unknown", versions)))
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/.well-known/openid-configuration", httpmock.NewStringResponder(200, `{"authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorize","token_endpoint":"https://topicus-keyhub.com/login/oauth2/token","revocation_endpoint":"https://topicus-keyhub.com/login/oauth2/revoke","device_authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorizedevice","issuer":"https://topicus-keyhub.com","jwks_uri":"https://topicus-keyhub.com/login/oauth2/jwks.json","scopes_supported":["openid","profile","manage_account","provisioning","access_vault","group_admin","global_admin"],"response_types_supported":["code","id_token","code token","code id_token","id_token token","code id_token token"],"response_modes_supported":["fragment","query"],"grant_types_supported":["authorization_code","client_credentials","implicit","password","refresh_token","urn:ietf:params:oauth:grant-type:device_code"],"code_challenge_methods_supported":["plain","S256"],"token_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"revocation_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"request_object_signing_alg_values_supported":["RS256","none"],"ui_locales_supported":["nl-NL"],"service_documentation":"https://topicus-keyhub.com/docs","request_parameter_supported":true,"request_uri_parameter_supported":true,"authorization_response_iss_parameter_supported":true,"subject_types_supported":["public"],"userinfo_endpoint":"https://topicus-keyhub.com/login/oauth2/userinfo","end_session_endpoint":"https://topicus-keyhub.com/login/oauth2/logout","id_token_signing_alg_values_supported":["RS256"],"userinfo_signing_alg_values_supported":["RS256"],"display_values_supported":["page"],"claim_types_supported":["normal"],"claims_supported":["sub","name","given_name","family_name","middle_name","nickname","preferred_username","picture","email","email_verified","gender","birthdate","zoneinfo","locale","phone_number","phone_number_verified","address","updated_at"],"claims_parameter_supported":true}`))
	httpmock.RegisterResponder("POST", "https://topicus-keyhub.com/login/oauth2/token", httpmock.NewStringResponder(200, `{"access_token": "a"}`))

	accountlist := model.AccountList{}
	sum := int64(1)
	for sum < 100 {
		acc := model.NewAccount("user" + strconv.FormatInt(sum, 10))
		acc.UUID = uuid.NewString()
		acc.DisplayName = "user " + strconv.FormatInt(sum, 10)
		acc.Links = append(acc.Links, model.Link{ID: sum, Rel: "self"})
		accountlist.Items = append(accountlist.Items, *acc)
		sum += sum
	}

	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/account/", httpmock.NewJsonResponderOrPanic(206, accountlist))
	httpmock.RegisterResponder("GET", `=~^https://topicus-keyhub\.com/keyhub/rest/v1/account/(\d+)\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get ID from request
			id := httpmock.MustGetSubmatchAsUint(req, 1) // 1=first regexp submatch
			return httpmock.NewJsonResponse(200, accountlist.Items[id-1])
		})

	grouplist := model.GroupList{}
	sum = int64(1)
	for sum < 100 {
		gr := model.NewEmptyGroup("user " + strconv.FormatInt(sum, 10))
		gr.UUID = "17502bdc-7a9f-4c9d-b355-81c9e9d7a12e"
		gr.Links = append(gr.Links, model.Link{ID: sum, Rel: "self"})
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

func TestVersioning(t *testing.T) {

	_, err := NewClientDefault("https://topicus-keyhub.com", "clientid", "clientsecret")
	if err != nil {
		t.Fatalf("ERROR %s", err)
	}
}

func TestAccounts(t *testing.T) {

	client, err := NewClientDefault("https://topicus-keyhub.com", "clientid", "clientsecret")
	if err != nil {
		t.Fatalf("ERROR %s", err)
	}

	accounts, err := client.Accounts.List()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(accounts) == 0 {
		t.Fatalf("ERROR no accounts found")
	}

	account, err := client.Accounts.GetById(1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if account == nil {
		t.Fatalf("ERROR account with id 1 not found")
	}
	if account.Self().ID != 1 {
		t.Fatalf("ERROR account with id 1 not found")
	}
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

func verifyQueryParams(t *testing.T, queryParams interface{}, expected string) {

	r, err := query.Values(queryParams)
	if err != nil {
		t.Fatalf("Could not parse query params")
	}

	result := r.Encode()

	if result != expected {
		t.Fatalf("Parse error, result `%s` did not match expected `%s`", result, expected)
	}

}

func TestQueries(t *testing.T) {

	var q model.ServiceAccountQueryParams

	q = model.ServiceAccountQueryParams{}
	verifyQueryParams(t, q, "")

	CreatedAfter, _ := time.Parse("2006-01-02", "2023-01-04")
	CreatedBefore, _ := time.Parse("2006-01-02", "2023-01-04")
	ModifiedSince, _ := time.Parse("2006-01-02", "2023-01-04")

	q = model.ServiceAccountQueryParams{
		UUID:                         "51f0cb1d-5745-4512-8d0d-bb28e2449d3f",
		CreatedAfter:                 CreatedAfter,
		CreatedBefore:                CreatedBefore,
		ModifiedSince:                ModifiedSince,
		Additional:                   nil,
		Exclude:                      []int64{1000},
		Id:                           []int64{1001},
		CQLQuery:                     "Blaat",
		Active:                       "BOTH",
		GroupOnSystem:                1002,
		GroupOnSystemOwners:          1003,
		Name:                         "Name",
		NameContains:                 "Contains",
		NameDoesNotStartWith:         "NotStartWith",
		NameStartsWith:               "StartsWith",
		Password:                     1004,
		PasswordRotation:             "MANUAL",
		RequestedGroupOnSystemOwners: 1005,
		System:                       1006,
		TechnicalAdministrator:       1007,
		Username:                     "Username",
	}

	verifyQueryParams(t, q, "Username=Username&active=BOTH&createdAfter=2023-01-04T00%3A00%3A00Z&createdBefore=2023-01-04T00%3A00%3A00Z&createdBefore=2023-01-04T00%3A00%3A00Z&exclude=1000&groupOnSystem=1002&groupOnSystemOwners=1003&id=1001&name=Name&nameContains=Contains&nameDoesNotStartWith=NotStartWith&nameStartsWith=StartsWith&password=1004&passwordRotation=MANUAL&q=Blaat&requestedGroupOnSystemOwners=1005&system=1006&technicalAdministrator=1007&uuid=51f0cb1d-5745-4512-8d0d-bb28e2449d3f")

}
