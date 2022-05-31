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

package model

import (
	"encoding/json"
	"fmt"
)

const (
	GROUP_RIGHT_MANAGER = "MANAGER"
	GROUP_RIGHT_MEMBER  = "NORMAL"
)

type GroupList struct {
	Items []Group `json:"items"`
}

type Group struct {
	Linkable
	AdditionalObjects *GroupAdditionalObjects `json:"additionalObjects,omitempty"`

	UUID           string            `json:"uuid,omitempty"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	ExtendedAccess string            `json:"extendedAccess"`
	AuditConfig    *GroupAuditConfig `json:"auditConfig,omitempty"`

	RotatingPasswordRequired  bool `json:"rotatingPasswordRequired,omitempty"`
	RecordTrail               bool `json:"recordTrail,omitempty"`
	PrivateGroup              bool `json:"privateGroup,omitempty"`
	HideAuditTrail            bool `json:"hideAuditTrail,omitempty"`
	ApplicationAdministration bool `json:"applicationAdministration,omitempty"`
	Auditor                   bool `json:"auditor,omitempty"`
	SingleManaged             bool `json:"singleManaged,omitempty"`

	AuthorizingGroupProvisioning *Group `json:"authorizingGroupProvisioning,omitempty"`
	AuthorizingGroupMembership   *Group `json:"authorizingGroupMembership,omitempty"`
	AuthorizingGroupAuditing     *Group `json:"authorizingGroupAuditing,omitempty"`
	NestedUnder                  *Group `json:"nestedUnder,omitempty"`
	//Classification               string `json:"classification,omitempty"`               //group_GroupClassificationPrimer{...}
	VaultRecovery string `json:"vaultRecovery,omitempty"` //[ NONE, RECOVERY_KEY_ONLY, FULL ]
}

func (g *Group) AddManager(account *Account) {

	g.addGroupAccount(account, GROUP_RIGHT_MANAGER)

}
func (g *Group) AddMember(account *Account) {

	g.addGroupAccount(account, GROUP_RIGHT_MEMBER)

}

func (g *Group) addGroupAccount(account *Account, groupRight string) {

	// Check if additionalObjects is set
	if g.AdditionalObjects == nil {
		g.AdditionalObjects = &GroupAdditionalObjects{}
	}

	// Check if Admin list is set
	if g.AdditionalObjects.Admins == nil {
		g.AdditionalObjects.Admins = &GroupAccountList{DType: "LinkableWrapper"}
	}

	// Add account as Group Manager
	g.AdditionalObjects.Admins.Items = append(g.AdditionalObjects.Admins.Items, *NewGroupAccount(account, groupRight))

}

type GroupAdditionalObjects struct {
	Admins *GroupAccountList `json:"admins,omitempty"`
}

func NewGroup(name string, groupadmin *Account) (result *Group) {
	result = &Group{Linkable: Linkable{DType: "group.Group"}, Name: name, ExtendedAccess: "NOT_ALLOWED"}
	result.AuditConfig = NewGroupAuditConfig()
	result.AddManager(groupadmin)
	return
}

// This struct needs a DType because it is used in an AdditionalObjects, which requires a $type
// to resolve which Java class is to be used
type GroupAccountList struct {
	DType string         `json:"$type,omitempty"`
	Items []GroupAccount `json:"items"`
}

type GroupAuditConfig struct {
	Linkable
	DType             string         `json:"$type,omitempty"`
	Months            MonthSelection `json:"months"`
	Permissions       []string       `json:"permissions,omitempty"`
	AdditionalObjects struct{}       `json:"additionalObjects,omitempty"`
}

func (gac *GroupAuditConfig) UnmarshalJSON(data []byte) error {

	type Alias GroupAuditConfig
	aux := &struct {
		Months []string `json:"months"`
		*Alias
	}{
		Alias: (*Alias)(gac),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	fmt.Println(aux.Months)
	gac.Months.Enable(aux.Months...)

	return nil
}

func (gac *GroupAuditConfig) MarshalJSON() ([]byte, error) {

	type Alias GroupAuditConfig
	aux := &struct {
		Months []string `json:"months"`
		*Alias
	}{
		Months: gac.Months.ToList(),
		Alias:  (*Alias)(gac),
	}

	return json.Marshal(aux)
}

func NewGroupAuditConfig() *GroupAuditConfig {

	gac := &GroupAuditConfig{DType: "group.GroupAuditConfig"}

	return gac

}

type GroupAccount struct {
	Linkable
	AdditionalObjects map[string]interface{} `json:"additionalObjects,omitempty"`

	UUID        string `json:"uuid,omitempty"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Rights      string `json:"rights"`
}

func NewGroupAccount(account *Account, rights string) *GroupAccount {
	ga := &GroupAccount{Linkable: Linkable{DType: "group.GroupAccount"}, Rights: rights}
	ga.Links = append(ga.Links, Link{ID: account.Links[0].ID, Rel: "self", Type: "auth.AccountPrimer", Href: account.Links[0].Href})
	ga.Username = account.Username
	ga.DisplayName = account.DisplayName
	ga.UUID = account.UUID
	return ga

}

type GroupQueryParams struct {
	UUID       string   `url:"uuid,omitempty"`
	Additional []string `url:"additional,omitempty"`
}
