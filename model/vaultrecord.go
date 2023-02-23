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
	"net/url"
	"time"
)

type VaultRecordList struct {
	Items []VaultRecord `json:"items"`
}

const (
	WARNINGPERIOD_AT_EXPIRATION RecordWarningPeriod = "AT_EXPIRATION"
	WARNINGPERIOD_TWO_WEEKS     RecordWarningPeriod = "TWO_WEEKS"
	WARNINGPERIOD_ONE_MONTH     RecordWarningPeriod = "ONE_MONTH"
	WARNINGPERIOD_TWO_MONTHS    RecordWarningPeriod = "TWO_MONTHS"
	WARNINGPERIOD_THREE_MONTHS  RecordWarningPeriod = "THREE_MONTHS"
	WARNINGPERIOD_SIX_MONTHS    RecordWarningPeriod = "SIX_MONTHS"
	WARNINGPERIOD_NEVER         RecordWarningPeriod = "NEVER"
)

type RecordWarningPeriod string

type VaultRecord struct {
	Linkable
	AdditionalObjects *VaultRecordAdditionalObjects `json:"additionalObjects,omitempty"`

	UUID          string              `json:"uuid,omitempty"`
	Name          string              `json:"name"`
	URL           string              `json:"url,omitempty"`
	Username      string              `json:"username,omitempty"`
	Color         string              `json:"color,omitempty"` // see bottom of file for color values
	Filename      string              `json:"filename,omitempty"`
	Types         []string            `json:"types,omitempty"`
	EndDate       time.Time           `json:"endDate,omitempty" layout:"2006-01-02"` // Layout don't work for json, only url but kept as reference
	WarningPeriod RecordWarningPeriod `json:"warningPeriod,omitempty"`
}

// Custom marshal function to format time.Time enddate to "Y-m-d" string
func (vr *VaultRecord) MarshalJSON() ([]byte, error) {

	type Alias VaultRecord
	aux := &struct {
		EndDate string `json:"endDate,omitempty"`
		*Alias
	}{
		EndDate: vr.EndDate.Format("2006-02-01"),
		Alias:   (*Alias)(vr),
	}

	return json.Marshal(aux)
}

// Custom unmarshal function to parse "Y-m-d" enddate to a time.Time field
func (vr *VaultRecord) UnmarshalJSON(data []byte) error {

	type Alias VaultRecord
	var err error
	aux := &struct {
		EndDate string `json:"endDate"`
		*Alias
	}{
		Alias: (*Alias)(vr),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	vr.EndDate, err = time.Parse("2006-02-01", aux.EndDate)
	if err != nil {
		return err
	}

	return nil
}

type VaultRecordAdditionalObjects struct {
	Audit  *AuditAdditionalObject             `json:"audit,omitempty"`
	Secret *VaultRecordSecretAdditionalObject `json:"secret,omitempty"`
}

type VaultRecordSecretAdditionalObject struct {
	DType    string  `json:"$type"`
	Password *string `json:"password,omitempty"`
	Totp     *string `json:"totp,omitempty"`
	File     *[]byte `json:"file"`
	Comment  *string `json:"comment"`
}

func NewVaultRecord(name string, secrets *VaultRecordSecretAdditionalObject) *VaultRecord {
	secrets.DType = "vault.VaultRecordSecrets"

	return &VaultRecord{Linkable: Linkable{DType: "vault.VaultRecord"}, Name: name,
		AdditionalObjects: &VaultRecordAdditionalObjects{Secret: secrets}}
}

func (r *VaultRecord) CreatedAt() time.Time {
	return r.AdditionalObjects.Audit.CreatedAt
}

func (r *VaultRecord) CreatedBy() string {
	return r.AdditionalObjects.Audit.CreatedBy
}

func (r *VaultRecord) LastModifiedAt() time.Time {
	return r.AdditionalObjects.Audit.LastModifiedAt
}

func (r *VaultRecord) LastModifiedBy() string {
	return r.AdditionalObjects.Audit.LastModifiedBy
}

func (r *VaultRecord) Comment() *string {
	return r.AdditionalObjects.Secret.Comment
}

func (r *VaultRecord) Password() *string {
	return r.AdditionalObjects.Secret.Password
}

func (r *VaultRecord) File() *[]byte {
	return r.AdditionalObjects.Secret.File
}

type VaultRecordQueryParams struct {
	UUID         string                            `url:"uuid,omitempty"`
	Name         string                            `url:"name,omitempty"`
	Filename     string                            `url:"filename,omitempty"`
	URL          string                            `url:"url,omitempty"`
	Username     string                            `url:"username,omitempty"`
	Color        string                            `url:"color,omitempty"` // see below for color values
	NameContains string                            `url:"nameContains,omitempty"`
	Additional   *VaultRecordAdditionalQueryParams `url:"additional,omitempty"`
}

type VaultRecordAdditionalQueryParams struct {
	Audit  bool `url:"audit"`
	Secret bool `url:"secret"`
}

func (p VaultRecordAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}

type VaultRecordSearchQueryParams struct {
	UUID                         string    `url:"uuid,omitempty"`
	ID                           string    `url:"id,omitempty"`
	AccessibleByClient           string    `url:"accessibleByClient,omitempty"`
	Additional                   []string  `url:"additional,omitempty"`
	AccessibleByAccount          string    `url:"accessibleByAccount,omitempty"`
	AccessibleByAccountAsManager string    `url:"accessibleByAccountAsManager,omitempty"`
	Any                          bool      `url:"any,omitempty"`
	CreatedAfter                 time.Time `url:"createdAfter,omitempty"`
	CreatedBefore                time.Time `url:"createdBefore,omitempty"`
	ModifiedSince                time.Time `url:"modifiedSince,omitempty"`
	Exclude                      []string  `url:"exclude,omitempty"`
	Q                            string    `url:"q,omitempty"`
	Color                        string    `url:"color,omitempty"` // see below for color values
	ExpireWarningBeforeOrAt      time.Time `url:"expireWarningBeforeOrAt,omitempty" layout:"2006-01-02"`
	Filename                     string    `url:"filename,omitempty"`
	HasNoPolicy                  bool      `url:"hasNoPolicy,omitempty"`
	HasParent                    bool      `url:"hasParent,omitempty"`
	HasValidPolicy               bool      `url:"hasValidPolicy,omitempty"`
	Name                         string    `url:"name,omitempty"`
	NameContains                 string    `url:"nameContains,omitempty"`
	Parent                       string    `url:"parent,omitempty"`
	Secret                       string    `url:"secret,omitempty"`
	ShareExpiresBeforeOrAt       time.Time `url:"shareExpiresBeforeOrAt,omitempty"`
	Url                          string    `url:"url,omitempty"`
	Username                     string    `url:"username,omitempty"`
	Uuid                         string    `url:"uuid,omitempty"`
	Vault                        string    `url:"vault,omitempty"`
}

var VaultRecordColorNone = "NONE"
var VaultRecordColorGreen = "GREEN"
var VaultRecordColorRed = "RED"
var VaultRecordColorBlue = "BLUE"
var VaultRecordColorDark = "DARK"
