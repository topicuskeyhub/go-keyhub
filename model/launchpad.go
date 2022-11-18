package model

const (
	LAUNCHPAD_TYPE_MANUAL = LaunchpadTileType("MANUAL")
	LAUNCHPAD_TYPE_SSO    = LaunchpadTileType("SSO_APPLICATION")
	LAUNCHPAD_TYPE_RECORD = LaunchpadTileType("VAULT_RECORD")
)

type LaunchpadTileType string

func NewLaunchPadTile(tiletype LaunchpadTileType) *LaunchPadTile {

	t := new(LaunchPadTile)
	t.DType = "launchpad.LaunchPadTile"

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
	LaunchPadTile
	DType string `json:"$type,omitempty"`
}

type LaunchPadTileSsoApplication struct {
	LaunchPadTile
	DType string `json:"$type,omitempty"`
	Uri   string `json:"uri"`
}

type LaunchPadTileManual struct {
	LaunchPadTile
	DType string `json:"$type,omitempty"`
	Uri   string `json:"uri"`
	Title string `json:"title"`
}
