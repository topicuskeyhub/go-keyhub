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
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/topicuskeyhub/go-keyhub/model"
)

func init() {
	httpmock.Activate()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/info", httpmock.NewJsonResponderOrPanic(200, model.NewVersionInfo("unknown", []int{2147483647, 54, 53})))
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/.well-known/openid-configuration", httpmock.NewStringResponder(200, `{"authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorize","token_endpoint":"https://topicus-keyhub.com/login/oauth2/token","revocation_endpoint":"https://topicus-keyhub.com/login/oauth2/revoke","device_authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorizedevice","issuer":"https://topicus-keyhub.com","jwks_uri":"https://topicus-keyhub.com/login/oauth2/jwks.json","scopes_supported":["openid","profile","manage_account","provisioning","access_vault","group_admin","global_admin"],"response_types_supported":["code","id_token","code token","code id_token","id_token token","code id_token token"],"response_modes_supported":["fragment","query"],"grant_types_supported":["authorization_code","client_credentials","implicit","password","refresh_token","urn:ietf:params:oauth:grant-type:device_code"],"code_challenge_methods_supported":["plain","S256"],"token_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"revocation_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"request_object_signing_alg_values_supported":["RS256","none"],"ui_locales_supported":["nl-NL"],"service_documentation":"https://topicus-keyhub.com/docs","request_parameter_supported":true,"request_uri_parameter_supported":true,"authorization_response_iss_parameter_supported":true,"subject_types_supported":["public"],"userinfo_endpoint":"https://topicus-keyhub.com/login/oauth2/userinfo","end_session_endpoint":"https://topicus-keyhub.com/login/oauth2/logout","id_token_signing_alg_values_supported":["RS256"],"userinfo_signing_alg_values_supported":["RS256"],"display_values_supported":["page"],"claim_types_supported":["normal"],"claims_supported":["sub","name","given_name","family_name","middle_name","nickname","preferred_username","picture","email","email_verified","gender","birthdate","zoneinfo","locale","phone_number","phone_number_verified","address","updated_at"],"claims_parameter_supported":true}`))
	httpmock.RegisterResponder("POST", "https://topicus-keyhub.com/login/oauth2/token", httpmock.NewStringResponder(200, `{"access_token": "a"}`))

	al := model.AccountList{}
	acc := model.Account{UUID: "4aaaa08c-023f-4b59-8164-a77013d95b90", Username: "user", DisplayName: "first user"}
	al.Items = append(al.Items, acc)
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/account/", httpmock.NewJsonResponderOrPanic(206, al))
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/account/1", httpmock.NewJsonResponderOrPanic(200, acc))

	gl := model.GroupList{}
	gr := model.Group{UUID: "17502bdc-7a9f-4c9d-b355-81c9e9d7a12e", Name: "first group", ExtendedAccess: "NOT_ALLOWED"}
	gl.Items = append(gl.Items, gr)
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/group/", httpmock.NewJsonResponderOrPanic(206, gl))
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/group/1", httpmock.NewJsonResponderOrPanic(200, gr))
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
}
