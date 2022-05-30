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

type GroupList struct {
	Items []Group `json:"items"`
}

type Group struct {
	Linkable
	AdditionalObjects *GroupAdditionalObjects `json:"additionalObjects,omitempty"`

	UUID           string `json:"uuid,omitempty"`
	Name           string `json:"name"`
	ExtendedAccess string `json:"extendedAccess"`
}

type GroupAdditionalObjects struct {
	Admins *GroupAccountList `json:"admins,omitempty"`
}

func NewEmptyGroup(name string) (result *Group) {
	return &Group{Linkable: Linkable{DType: "group.Group"}, Name: name, ExtendedAccess: "NOT_ALLOWED"}
}

func NewGroup(name string, groupadmin *Account) (result *Group) {
	gal := GroupAccountList{DType: "LinkableWrapper"}
	gra := NewGroupAccount(groupadmin, "MANAGER")
	gal.Items = append(gal.Items, *gra)
	gao := GroupAdditionalObjects{Admins: &gal}

	result = &Group{Linkable: Linkable{DType: "group.Group"}, Name: name, ExtendedAccess: "NOT_ALLOWED", AdditionalObjects: &gao}
	return
}

// GroupAccountList This struct needs a DType because it is used in an AdditionalObjects, which requires a $type to resolve which Java class is to be used
type GroupAccountList struct {
	DType string         `json:"$type,omitempty"`
	Items []GroupAccount `json:"items"`
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
	ga.Links = append(ga.Links, Link{ID: account.Links[0].ID, Rel: "self"})

	return ga
}

type GroupQueryParams struct {
	UUID       string   `url:"uuid,omitempty"`
	Additional []string `url:"additional,omitempty"`
}
