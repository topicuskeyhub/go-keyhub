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
	"net/url"
	"strings"
)

const (
	GROUP_RIGHT_MANAGER = "MANAGER"
	GROUP_RIGHT_MEMBER  = "NORMAL"

	GROUP_EXT_ACCESS_NOT = "NOT_ALLOWED"
	GROUP_EXT_ACCESS_1W  = "ONE_WEEK"
	GROUP_EXT_ACCESS_2W  = "TWO_WEEKS"

	VAULT_RECOVERY_NONE     = "NONE"
	VAULT_RECOVERY_KEY_ONLY = "RECOVERY_KEY_ONLY"
	VAULT_RECOVERY_FULL     = "FULL"
)

type GroupList struct {
	Items []Group `json:"items"`
}

type GroupPrimer struct {
	Linkable
	Admin bool `json:"admin,omitempty"`

	UUID string `json:"uuid,omitempty"`
	Name string `json:"name"`
}

type Group struct {
	GroupPrimer

	AdditionalObjects *GroupAdditionalObjects `json:"additionalObjects,omitempty"`

	Description    string            `json:"description,omitempty"`
	ExtendedAccess string            `json:"extendedAccess,omitempty"`
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
	VaultRecovery string `json:"vaultRecovery,omitempty"`
}

// AddManager Add Account as Manager
func (g *Group) AddManager(account *Account) {
	g.addGroupAccount(account, GROUP_RIGHT_MANAGER)
}

// AddMember Add Account as Member
func (g *Group) AddMember(account *Account) {
	g.addGroupAccount(account, GROUP_RIGHT_MEMBER)
}

// AsPrimer Return Group with only Primer data
func (g *Group) AsPrimer() *Group {
	group := &Group{}
	group.GroupPrimer = g.GroupPrimer
	return group
}

// ToPrimer Convert to GroupPrimer
func (g *Group) ToPrimer() *GroupPrimer {
	groupPrimer := g.GroupPrimer
	return &groupPrimer
}

func (g *Group) DisableExtendedAccess() {
	g.ExtendedAccess = GROUP_EXT_ACCESS_NOT
}
func (g *Group) EnableExtendedAccess1W() {
	g.ExtendedAccess = GROUP_EXT_ACCESS_1W
}
func (g *Group) EnableExtendedAccess2W() {
	g.ExtendedAccess = GROUP_EXT_ACCESS_2W
}

func (g *Group) DisableVaultRecovery() {
	g.VaultRecovery = VAULT_RECOVERY_NONE
}
func (g *Group) KeyOnlyVaultRecovery() {
	g.VaultRecovery = VAULT_RECOVERY_KEY_ONLY
}
func (g *Group) FullVaultRecovery() {
	g.VaultRecovery = VAULT_RECOVERY_FULL
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

func NewEmptyGroup(name string) (result *Group) {
	return &Group{GroupPrimer: GroupPrimer{Linkable: Linkable{DType: "group.Group"}, Name: name}, ExtendedAccess: "NOT_ALLOWED"}
}

func NewGroup(name string, groupadmin *Account) (result *Group) {
	result = NewEmptyGroup(name)
	result.AuditConfig = NewGroupAuditConfig()
	result.AddManager(groupadmin)
	return
}

// GroupAccountList This struct needs a DType because it is used in an AdditionalObjects, which requires a $type to resolve which Java class is to be used
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
	UUID       string                      `url:"uuid,omitempty"`
	Additional *GroupAdditionalQueryParams `url:"additional,omitempty"`
}

type GroupAdditionalQueryParams struct {
	Audit  bool `url:"audit"`
	Admins bool `url:"admins"`
}

// EncodeValues Custom url encoder to convert bools to list
func (p GroupAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}

const (
	PRGRP_SECURITY_LEVEL_LOW    ProvisioningGroupSecurityLevel = "LOW"
	PRGRP_SECURITY_LEVEL_MEDIUM ProvisioningGroupSecurityLevel = "MEDIUM"
	PRGRP_SECURITY_LEVEL_HIGH   ProvisioningGroupSecurityLevel = "HIGH"
)

// Section: Group
func NewProvisioningGroup() *ProvisioningGroup {

	pg := ProvisioningGroup{
		Linkable: Linkable{
			DType: "group.ProvisioningGroup",
		},
		SecurityLevel:      PRGRP_SECURITY_LEVEL_HIGH,
		StaticProvisioning: false,
	}
	return &pg
}

// ProvisioningGroup instance of group.ProvisioningGroup
type ProvisioningGroup struct {
	Linkable

	GroupOnSystem      *GroupOnSystem                 `json:"groupOnSystem,omitempty"`
	Group              *Group                         `json:"group"`
	SecurityLevel      ProvisioningGroupSecurityLevel `json:"securityLevel"`
	StaticProvisioning bool                           `json:"staticProvisioning"`
}

func (p *ProvisioningGroup) SetSecurityLevelString(level string) error {

	switch strings.ToUpper(level) {
	case string(PRGRP_SECURITY_LEVEL_HIGH):
		p.SecurityLevel = PRGRP_SECURITY_LEVEL_HIGH
	case string(PRGRP_SECURITY_LEVEL_MEDIUM):
		p.SecurityLevel = PRGRP_SECURITY_LEVEL_MEDIUM
	case string(PRGRP_SECURITY_LEVEL_LOW):
		p.SecurityLevel = PRGRP_SECURITY_LEVEL_LOW
	default:
		return fmt.Errorf("value %s is not a valid level", level)
	}
	return nil

}

type ProvisioningGroupSecurityLevel string
