package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	LAUNCHPAD_TYPE_MANUAL = LaunchpadTileType("MANUAL")
	LAUNCHPAD_TYPE_SSO    = LaunchpadTileType("SSO_APPLICATION")
	LAUNCHPAD_TYPE_RECORD = LaunchpadTileType("VAULT_RECORD")
)

type LaunchpadTileType string

func NewLaunchPadTile(tiletype LaunchpadTileType) *LaunchPadTile {

	t := new(LaunchPadTile)
	t.DType = "launchpad.LaunchPadTile"
	t.Type = tiletype

	switch tiletype {
	case LAUNCHPAD_TYPE_SSO:
		t.DType = "launchpad.SsoApplicationLaunchpadTile"
	case LAUNCHPAD_TYPE_RECORD:
		t.DType = "launchpad.VaultRecordLaunchpadTile"
	case LAUNCHPAD_TYPE_MANUAL:
		t.DType = "launchpad.ManualLaunchpadTile"
	}

	return t

}

func NewManualLaunchPadTile(title string, uri string, group *Group) *LaunchPadTile {

	tile := NewLaunchPadTile(LAUNCHPAD_TYPE_MANUAL)
	tile.Title = title
	tile.Uri = uri
	tile.Group = group

	return tile

}

func NewVaultRecordLaunchPadTile(record *VaultRecord) *LaunchPadTile {

	tile := NewLaunchPadTile(LAUNCHPAD_TYPE_RECORD)
	tile.VaultRecord = record

	return tile

}

func NewApplicationLaunchPadTile(uri string, application *ClientApplication) *LaunchPadTile {

	tile := NewLaunchPadTile(LAUNCHPAD_TYPE_SSO)
	tile.Application = application
	tile.Uri = uri

	return tile

}

type LaunchPadTileList struct {
	DType string          `json:"$type,omitempty"`
	Items []LaunchPadTile `json:"items"`
}

type LaunchPadTilePrimer struct {
	Linkable
	DType string `json:"$type,omitempty"`
}

type LaunchPadTile struct {
	LaunchPadTilePrimer
	DType string `json:"$type,omitempty"`

	Type          LaunchpadTileType  `json:"type"`
	IdenticonCode int32              `json:"identiconCode"`
	Logo          []byte             `json:"logo,omitempty"`
	Group         *Group             `json:"group,omitempty"`
	Application   *ClientApplication `json:"application,omitempty"`
	VaultRecord   *VaultRecord       `json:"vaultRecord,omitempty"`
	Uri           string             `json:"uri,omitempty"`
	Title         string             `json:"title,omitempty"`
}

// AsPrimer Return LaunchPadTile with only Primer data
func (t *LaunchPadTile) AsPrimer() *LaunchPadTile {
	tile := &LaunchPadTile{}
	tile.LaunchPadTilePrimer = t.LaunchPadTilePrimer
	return tile
}

// ToPrimer Convert to LaunchPadTilePrimer
func (t *LaunchPadTile) ToPrimer() *LaunchPadTilePrimer {
	tile := t.LaunchPadTilePrimer
	return &tile
}

// Custom marshal function to format time.Time enddate to "Y-m-d" string
func (t LaunchPadTile) MarshalJSON() ([]byte, error) {

	if t.Title != "" && t.Type != LAUNCHPAD_TYPE_MANUAL {
		return nil, fmt.Errorf("Title can not be set for LaunchPadTile type %s", t.Type)
	}

	if t.Uri != "" && t.Type == LAUNCHPAD_TYPE_RECORD {
		return nil, fmt.Errorf("Uri can not be set for LaunchPadTile type %s", t.Type)
	}
	type Alias LaunchPadTile
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(&t),
	}

	return json.Marshal(aux)
}

/*
launchpad_LaunchpadTile	{…}
launchpad_LaunchpadTilePrimer	{…}
launchpad_VaultRecordLaunchpadTile	{…}
launchpad_SsoApplicationLaunchpadTile
LinkableWrapper.launchpad_LaunchpadTile
launchpad_DisplayedLaunchpadTile
launchpad_DisplayedLaunchpadTiles
launchpad_ManualLaunchpadTile
*/

type LaunchPadTileVaultRecord struct {
	Uri   string `json:"-"`
	Title string `json:"-"`
	LaunchPadTile
}

type LaunchPadTileSsoApplication struct {
	Title string `json:"-"`
	LaunchPadTile
}

type LaunchPadTileManual struct {
	LaunchPadTile
}

type LaunchPadTileQueryParams struct {
	Additional    interface{} `url:"additional,omitempty"`
	Any           bool        `url:"any,omitempty"`
	CreatedAfter  time.Time   `url:"createdAfter,omitempty"`
	CreatedBefore time.Time   `url:"createdBefore,omitempty"`
	exclude       []int64     `url:"exclude,omitempty"`
	Id            int64       `url:"id,omitempty"`
	ModifiedSince time.Time   `url:"modifiedSince,omitempty"`
	Q             string      `url:"q,omitempty"`
	Application   int64       `url:"application,omitempty"`
	Group         int64       `url:"group,omitempty"`
	Title         string      `url:"title,omitempty"`
	VaultRecords  []int64     `url:"vaultRecord,omitempty"`
}

type LaunchPadTileAdditionalQueryParams struct {
	Audit bool `url:"audit"`
}

// EncodeValues Custom url encoder to convert bools to list
func (p LaunchPadTileAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}
