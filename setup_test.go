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
	"github.com/jarcoal/httpmock"
	keyhubmodel "github.com/topicuskeyhub/go-keyhub/model"
)

func SetupTest() {
	httpmock.Activate()

	// Exact URL match
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/keyhub/rest/v1/info", httpmock.NewJsonResponderOrPanic(200, keyhubmodel.NewVersionInfo("unknown", []int{2147483647, 54, 53})))
	httpmock.RegisterResponder("GET", "https://topicus-keyhub.com/.well-known/openid-configuration", httpmock.NewStringResponder(200, `{"authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorize","token_endpoint":"https://topicus-keyhub.com/login/oauth2/token","revocation_endpoint":"https://topicus-keyhub.com/login/oauth2/revoke","device_authorization_endpoint":"https://topicus-keyhub.com/login/oauth2/authorizedevice","issuer":"https://topicus-keyhub.com","jwks_uri":"https://topicus-keyhub.com/login/oauth2/jwks.json","scopes_supported":["openid","profile","manage_account","provisioning","access_vault","group_admin","global_admin"],"response_types_supported":["code","id_token","code token","code id_token","id_token token","code id_token token"],"response_modes_supported":["fragment","query"],"grant_types_supported":["authorization_code","client_credentials","implicit","password","refresh_token","urn:ietf:params:oauth:grant-type:device_code"],"code_challenge_methods_supported":["plain","S256"],"token_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"revocation_endpoint_auth_methods_supported":["client_secret_basic","client_secret_post"],"request_object_signing_alg_values_supported":["RS256","none"],"ui_locales_supported":["nl-NL"],"service_documentation":"https://topicus-keyhub.com/docs","request_parameter_supported":true,"request_uri_parameter_supported":true,"authorization_response_iss_parameter_supported":true,"subject_types_supported":["public"],"userinfo_endpoint":"https://topicus-keyhub.com/login/oauth2/userinfo","end_session_endpoint":"https://topicus-keyhub.com/login/oauth2/logout","id_token_signing_alg_values_supported":["RS256"],"userinfo_signing_alg_values_supported":["RS256"],"display_values_supported":["page"],"claim_types_supported":["normal"],"claims_supported":["sub","name","given_name","family_name","middle_name","nickname","preferred_username","picture","email","email_verified","gender","birthdate","zoneinfo","locale","phone_number","phone_number_verified","address","updated_at"],"claims_parameter_supported":true}`))
	httpmock.RegisterResponder("POST", "https://topicus-keyhub.com/login/oauth2/token", httpmock.NewStringResponder(200, `{"access_token": "a", "vaultSession": "b"}`))
}
