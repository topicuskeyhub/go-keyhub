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

import "time"

type VaultRecordList struct {
	Items []VaultRecord `json:"items"`
}

type VaultRecord struct {
	Linkable
	AdditionalObjects *VaultRecordAdditionalObjects `json:"additionalObjects,omitempty"`

	UUID     string   `json:"uuid,omitempty"`
	Name     string   `json:"name"`
	URL      string   `json:"url,omitempty"`
	Username string   `json:"username,omitempty"`
	Color    string   `json:"color,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Types    []string `json:"types,omitempty"`
}

type VaultRecordAdditionalObjects struct {
	Audit  *VaultRecordAuditAdditionalObject  `json:"audit,omitempty"`
	Secret *VaultRecordSecretAdditionalObject `json:"secret,omitempty"`
}

type VaultRecordAuditAdditionalObject struct {
	CreatedAt      time.Time `json:"createdAt"`
	CreatedBy      string    `json:"createdBy"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	LastModifiedBy string    `json:"lastModifiedBy"`
}

type VaultRecordSecretAdditionalObject struct {
	Password string `json:"password,omitempty"`
	Totp     string `json:"totp,omitempty"`
	File     []byte `json:"file"`
	Comment  string `json:"comment"`
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

func (r *VaultRecord) Comment() string {
	return r.AdditionalObjects.Secret.Comment
}

func (r *VaultRecord) Password() string {
	return r.AdditionalObjects.Secret.Password
}

func (r *VaultRecord) File() []byte {
	return r.AdditionalObjects.Secret.File
}

type VaultRecordQueryParams struct {
	UUID       string   `url:"uuid,omitempty"`
	Additional []string `url:"additional,omitempty"`
	Any        bool     `url:"any,omitempty"`
}

type RecordOptions struct {
	Audit  bool
	Secret bool
}
