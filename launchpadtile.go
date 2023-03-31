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
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
	"github.com/topicuskeyhub/go-keyhub/model"
)

// LaunchPadTileService Service to manage launch pad tiles in keyhub
type LaunchPadTileService struct {
	sling *sling.Sling
}

func newLaunchPadTileService(sling *sling.Sling) *LaunchPadTileService {
	return &LaunchPadTileService{
		sling: sling.Path("/keyhub/rest/v1/launchpadtile/"),
	}
}

// Create a new launch pad tile in Keyhub
func (s *LaunchPadTileService) Create(tile *model.LaunchPadTile) (result *model.LaunchPadTile, err error) {

	tiles := new(model.LaunchPadTileList)
	results := new(model.LaunchPadTileList)
	errorReport := new(model.ErrorReport)
	tiles.Items = append(tiles.Items, *tile)

	_, err = s.sling.New().Post("").BodyJSON(tiles).Receive(results, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not create LaunchPadTile.")
	}
	if err == nil {
		if len(results.Items) > 0 {
			result = &results.Items[0]
		} else {
			err = fmt.Errorf("created LaunchPadTile not found")
		}
	}

	return
}

// List all available launch pad tiles.
func (s *LaunchPadTileService) List(queryParams *model.LaunchPadTileQueryParams) (tiles []model.LaunchPadTile, err error) {
	searchRange := model.NewRange()

	if queryParams == nil {
		queryParams = new(model.LaunchPadTileQueryParams)
	}

	var response *http.Response

	for ok := true; ok; ok = searchRange.NextPage() {

		errorReport := new(model.ErrorReport)
		results := new(model.LaunchPadTileList)
		response, err = s.sling.New().Get("").QueryStruct(*queryParams).Add(searchRange.GetRequestRangeHeader()).Add(searchRange.GetRequestModeHeader()).Receive(results, errorReport)
		searchRange.ParseResponse(response)

		if errorReport.Code > 0 {
			err = errorReport.Wrap("Could not get LaunchPadTiles")
		}
		if err == nil {
			if len(results.Items) > 0 {
				tiles = append(tiles, results.Items...)
			}
		}
	}
	return
}

// GetById Retrieve a launch pad tile by keyhub id
func (s *LaunchPadTileService) GetById(id int64) (result *model.LaunchPadTile, err error) {
	al := new(model.LaunchPadTile)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Get(idString).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not get LaunchPadTile %q. Error: %s", idString)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("LaunchPadTile %q not found", idString)
		return
	}

	return al, nil
}

// GetById Retrieve a launch pad tile by keyhub id
func (s *LaunchPadTileService) DeleteById(id int64) (err error) {
	al := new(model.LaunchPadTile)
	errorReport := new(model.ErrorReport)
	idString := strconv.FormatInt(id, 10)

	_, err = s.sling.New().Delete(idString).Receive(al, errorReport)
	if errorReport.Code > 0 {
		err = errorReport.Wrap("Could not delete LaunchPadTile %q. Error: %s", idString)
		return
	}
	if err == nil && al == nil {
		err = fmt.Errorf("LaunchPadTile %q not found", idString)
		return
	}

	return nil
}
