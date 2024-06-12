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
	"fmt"
	"strings"

	"github.com/dghubble/sling"
	"github.com/topicuskeyhub/go-keyhub/model"
)

type VersionService struct {
	sling *sling.Sling
	info  *model.VersionInfo
}

const (
	/* KeyHub json mediatype */
	mediatype = "application/vnd.topicus.keyhub+json"
)

func newVersionService(sling *sling.Sling) *VersionService {
	return &VersionService{
		sling: sling.Path("/keyhub/rest/v1/info"),
	}
}

func (s *VersionService) Get() (v *model.VersionInfo, err error) {
	results := new(model.VersionInfo)
	errorReport := new(model.ErrorReport)

	resp, err := s.sling.New().Get("").Receive(results, errorReport)
	if err != nil {
		return
	}
	if errorReport.Code > 0 {
		return nil, errorReport.Wrap("Could not fetch acceptable contract versions")
	}
	if resp.StatusCode >= 300 {
		err = fmt.Errorf("Could not fetch acceptable contract versions. Error: %s", resp.Status)
		return
	}
	results.KeyhubVersion = strings.TrimPrefix(results.KeyhubVersion, "keyhub-")

	return results, nil
}

func (s *VersionService) CheckAndUpdateVersionedSling(version int, base *sling.Sling) (headerVersion string, err error) {

	if s.info == nil {
		s.info, err = s.Get()
		if err != nil {
			return "", err
		}
	}

	if version > 0 {

		isContractVersionSupported := false
		for _, contractVersion := range s.info.ContractVersions {
			if version == contractVersion {
				isContractVersionSupported = true
				break
			}
		}
		if !isContractVersionSupported {
			return "", fmt.Errorf("KeyHub %v does not support api contract version %v", s.info.KeyhubVersion, version)
		}

		headerVersion = fmt.Sprintf("%d", version)
	} else {
		headerVersion = "latest"
	}

	base.Set("Accept", fmt.Sprintf("%v;version=%s", mediatype, headerVersion))
	base.Set("Content-Type", fmt.Sprintf("%v;version=%s", mediatype, headerVersion))

	return
}
