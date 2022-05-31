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
}

func newVersionService(sling *sling.Sling) *VersionService {
	return &VersionService{
		sling: sling.Path("/keyhub/rest/v1/info"),
	}
}

func (s *VersionService) Get() (v *model.VersionInfo, err error) {
	results := new(model.VersionInfo)
	errorReport := new(model.ErrorReport)

	resp, err := s.sling.New().Get("").Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = fmt.Errorf("Could not get acceptable contract versions. Error: %s", errorReport.Message)
	}
	if resp.StatusCode >= 300 {
		err = fmt.Errorf("Could not fetch acceptable contract versions. Error: %s", resp.Status)
	}
	if err == nil {
		results.KeyhubVersion = strings.TrimPrefix(results.KeyhubVersion, "keyhub-")
	}

	v = results
	return
}
