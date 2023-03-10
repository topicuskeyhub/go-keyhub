package model

// Section: provisioning

import (
	"fmt"
	"github.com/gosimple/slug"
	"net/url"
	"time"
)

const (
	GOS_TYPE_POSIX                 GroupOnSystemType = "POSIX_GROUP"
	GOS_TYPE_GROUP_OF_NAMES        GroupOnSystemType = "GROUP_OF_NAMES"
	GOS_TYPE_GROUP_OF_UNIQUE_NAMES GroupOnSystemType = "GROUP_OF_UNIQUE_NAMES"
	GOS_TYPE_GROUP                 GroupOnSystemType = "GROUP"
	GOS_TYPE_AZURE_ROLE            GroupOnSystemType = "AZURE_ROLE"
	GOS_TYPE_AZURE_UNIFIED_GROUP   GroupOnSystemType = "AZURE_UNIFIED_GROUP"
	GOS_TYPE_AZURE_SECURITY_GROUP  GroupOnSystemType = "AZURE_SECURITY_GROUP"

	PSYSTEM_TYPE_LDAP                             ProvisionedSystemType = "LDAP"
	PSYSTEM_TYPE_ACTIVE_DIRECTORY                 ProvisionedSystemType = "ACTIVE_DIRECTORY"
	PSYSTEM_TYPE_AZURE_TENANT                     ProvisionedSystemType = "AZURE_TENANT"
	PSYSTEM_TYPE_SOURCE_LDAP_DIRECTORY            ProvisionedSystemType = "SOURCE_LDAP_DIRECTORY"
	PSYSTEM_TYPE_SOURCE_AZURE_OIDC_DIRECTORY      ProvisionedSystemType = "SOURCE_AZURE_OIDC_DIRECTORY"
	PSYSTEM_TYPE_SOURCE_AZURE_SYNC_LDAP_DIRECTORY ProvisionedSystemType = "SOURCE_AZURE_SYNC_LDAP_DIRECTORY"
	PSYSTEM_TYPE_INTERNAL_LDAP                    ProvisionedSystemType = "INTERNAL_LDAP"
)

// GroupOnSystemType Use constants as enum for type
type GroupOnSystemType string

// ProvisionedSystemType Use constants as enum for type
type ProvisionedSystemType string

// NewGroupOnSystemList Initialize a new GroupOnSystemList
func NewGroupOnSystemList() *GroupOnSystemList {
	return &GroupOnSystemList{}
}

// GroupOnSystemList List of GroupOnSystems
type GroupOnSystemList struct {
	Items []GroupOnSystem `json:"items"`
}

// GroupOnSystem Group On System object
type GroupOnSystem struct {
	Linkable

	Type              GroupOnSystemType  `json:"type,omitempty"`
	NameInSystem      string             `json:"nameInSystem,omitempty"`
	ShortNameInSystem string             `json:"shortNameInSystem,omitempty"` // Read Only
	DisplayName       string             `json:"displayName,omitempty"`
	System            *ProvisionedSystem `json:"system,omitempty"`
	Owner             *Group             `json:"owner,omitempty"`

	AdditionalObjects *GroupOnSystemAdditionalObject `json:"additionalObjects,omitempty"`
}

// NewProvisioningGroupList Initialize a new ProvisioningGroupList
func NewProvisioningGroupList() *ProvisioningGroupList {
	return &ProvisioningGroupList{
		DType: "LinkableWrapper",
		Items: make([]ProvisioningGroup, 0),
	}
}

// ProvisioningGroupList List of ProvisioningGroup
type ProvisioningGroupList struct {
	DType string              `json:"$type,omitempty"`
	Items []ProvisioningGroup `json:"items"`
}

// NewGroupOnSystemAdditionalObject Initialize a new GroupOnSystemAdditionalObject
func NewGroupOnSystemAdditionalObject() *GroupOnSystemAdditionalObject {
	return &GroupOnSystemAdditionalObject{
		ProvGroups: NewProvisioningGroupList(),
	}
}

// GroupOnSystemAdditionalObject Additional objects for GroupOnSystem
type GroupOnSystemAdditionalObject struct {
	ProvGroups *ProvisioningGroupList `json:"provgroups,omitempty"`
}

// AddProvGroup Add Additional ProvisioningGroup
func (g *GroupOnSystem) AddProvGroup(group ...ProvisioningGroup) {

	if g.AdditionalObjects == nil {
		g.AdditionalObjects = NewGroupOnSystemAdditionalObject()
	}
	g.AdditionalObjects.ProvGroups.Items = append(
		g.AdditionalObjects.ProvGroups.Items,
		group...,
	)
}

// NoProvGroups Set empty list for ProvisioningGroups
func (g *GroupOnSystem) NoProvGroups() {
	g.AdditionalObjects = NewGroupOnSystemAdditionalObject()
}

// SetType Set type of GroupOnSystem, use one of the GOS_TYPE_* constants
func (g *GroupOnSystem) SetType(typename GroupOnSystemType) {
	g.Type = typename
}

// SetType Set type of GroupOnSystem, use one of the GOS_TYPE_* constants
func (g *GroupOnSystem) SetTypeString(typename string) error {

	switch typename {
	case string(GOS_TYPE_POSIX):
		g.Type = GOS_TYPE_POSIX
	case string(GOS_TYPE_GROUP_OF_NAMES):
		g.Type = GOS_TYPE_GROUP_OF_NAMES
	case string(GOS_TYPE_GROUP_OF_UNIQUE_NAMES):
		g.Type = GOS_TYPE_GROUP_OF_UNIQUE_NAMES
	case string(GOS_TYPE_GROUP):
		g.Type = GOS_TYPE_GROUP
	case string(GOS_TYPE_AZURE_ROLE):
		g.Type = GOS_TYPE_AZURE_ROLE
	case string(GOS_TYPE_AZURE_UNIFIED_GROUP):
		g.Type = GOS_TYPE_AZURE_UNIFIED_GROUP
	case string(GOS_TYPE_AZURE_SECURITY_GROUP):
		g.Type = GOS_TYPE_AZURE_SECURITY_GROUP
	default:
		return fmt.Errorf("value '%s' is not valid for type", typename)
	}

	return nil
}

// SetName Set DisplayName and NameInSystem from name given, NameInSystem is slugified
func (g *GroupOnSystem) SetName(name string) {
	g.DisplayName = name
	g.NameInSystem = slug.Make(name)
}

// NewGroupOnSystem Initialize a new GroupOnSystem
func NewGroupOnSystem() *GroupOnSystem {
	return &GroupOnSystem{
		Linkable: Linkable{DType: "provisioning.GroupOnSystem"},
	}
}

// NewGroupOnSystemWithEmptyProvGroup Initialize a new GroupOnSystem with an empty provgroup list, disabling owner as default provgroup
func NewGroupOnSystemWithEmptyProvGroup() *GroupOnSystem {
	return &GroupOnSystem{
		Linkable:          Linkable{DType: "provisioning.GroupOnSystem"},
		AdditionalObjects: NewGroupOnSystemAdditionalObject(),
	}
}

// NewProvisionedSystem Initialize a new ProvisionedSystem
func NewProvisionedSystem() *ProvisionedSystem {
	ps := &ProvisionedSystem{}
	ps.Linkable = Linkable{DType: "provisioning.ProvisionedSystem"}
	return ps
}

// ProvisionedSystemPrimer Primer Class for ProvisionedSystem
type ProvisionedSystemPrimer struct {
	Linkable
	Active bool                  `json:"active,omitempty"`
	UUID   string                `json:"uuid,omitempty"`
	Name   string                `json:"name,omitempty"`
	Type   ProvisionedSystemType `json:"type,omitempty"`
}

// ProvisionedAbstract Base parameters for ProvisionedSystems
type ProvisionedAbstract struct {
	Host                       string       `json:"host,omitempty"`
	Port                       int          `json:"port,omitempty"`
	FailoverHost               string       `json:"failoverHost,omitempty"`
	TLS                        string       `json:"tls,omitempty"`
	BaseDN                     string       `json:"baseDN,omitempty"`
	BindDN                     string       `json:"bindDN,omitempty"`
	BindPassword               string       `json:"bindPassword,omitempty"`
	UserDN                     string       `json:"userDN,omitempty"`
	GroupDN                    string       `json:"groupDN,omitempty"`
	SshPublicKeySupported      bool         `json:"sshPublicKeySupported,omitempty"`
	ObjectClasses              string       `json:"objectClasses,omitempty"`
	TrustedCertificate         *Certificate `json:"trustedCertificate,omitempty"`
	FailoverTrustedCertificate *Certificate `json:"failoverTrustedCertificate,omitempty"`
	ClientCertificate          *Certificate `json:"clientCertificate,omitempty"`
	Attributes                 interface{}  `json:"attributes,omitempty"`
}

// ProvisionedSystem Generic ProvisionedSystem Class
type ProvisionedSystem struct {
	ProvisionedSystemPrimer

	ProvisionedAbstract

	AccountCount           int    `json:"accountCount,omitempty"`
	UsernamePrefix         string `json:"usernamePrefix,omitempty"`
	TechnicalAdministrator *Group `json:"technicalAdministrator,omitempty"`
	ExternalUUID           string `json:"externalUUID,omitempty"`
}

// AsPrimer Return ProvisionedSystem with only Primer data
func (p *ProvisionedSystem) AsPrimer() *ProvisionedSystem {
	system := &ProvisionedSystem{}
	system.ProvisionedSystemPrimer = p.ProvisionedSystemPrimer
	return system
}

// AsPrimer Convert to ProvisionedSystemPrimer
func (p *ProvisionedSystem) ToPrimer() *ProvisionedSystemPrimer {
	primer := p.ProvisionedSystemPrimer
	return &primer
}

// NewProvisionedSystemList Initialize a new ProvisionedSystemList
func NewProvisionedSystemList() *ProvisionedSystemList {
	return &ProvisionedSystemList{
		DType: "LinkableWrapper",
	}
}

// ProvisionedSystemList list of ProvisionedSystem
type ProvisionedSystemList struct {
	DType string              `json:"$type,omitempty"`
	Items []ProvisionedSystem `json:"items,omitempty"`
}

// ProvisionedAD Additional Parameters for specific ProvisionedSystem type
type ProvisionedAD struct {
	//provisioning.ProvisionedAD
	SamAccountNameScheme string `json:"samAccountNameScheme,omitempty"`
	//"OMIT"
	//"TRUNCATE"
	//"TRANSFER"
	//"TRANSFER_TRUNCATE"
	//"USERNAME"
}

// ProvisionedLDAPDirectory Additional Parameters for specific ProvisionedSystem type
type ProvisionedLDAPDirectory struct {
	//provisioning.ProvisionedLDAPDirectory
	Directory interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
	GroupDN   string      `json:"groupDN,omitempty"`
}

// ProvisionedAzureTenant Additional Parameters for specific ProvisionedSystem type
type ProvisionedAzureTenant struct {
	//provisioning.ProvisionedAzureTenant
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Tenant       string `json:"tenant,omitempty"`
	IdpDomain    string `json:"idpDomain,omitempty"`
}

// ProvisionedInternalLDAP Additional Parameters for specific ProvisionedSystem type
type ProvisionedInternalLDAP struct {
	//provisioning.ProvisionedInternalLDAP
	Client interface{} `json:"client,omitempty"` // client_ldapClient
}

// ProvisionedLDAP Additional Parameters for specific ProvisionedSystem type
type ProvisionedLDAP struct {
	//provisioning.ProvisionedLDAP
	GID           int         `json:"gid,omitempty"`
	HashingScheme string      `json:"hashingScheme"`
	Numbering     interface{} `json:"numbering,omitempty"` // provisioning_ProvisionNumberSequence
}

// ProvisionedAzureSyncLDAPDirectory Additional Parameters for specific ProvisionedSystem type
type ProvisionedAzureSyncLDAPDirectory struct {
	//provisioning.ProvisionedAzureSyncLDAPDirectory
	ClientId     string      `json:"clientId,omitempty"`
	ClientSecret string      `json:"clientSecret,omitempty"`
	Tenant       string      `json:"tenant,omitempty"`
	Directory    interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
}

// ProvisionedAzureOIDCDirectory Additional Parameters for specific ProvisionedSystem type
type ProvisionedAzureOIDCDirectory struct {
	//provisioning.ProvisionedAzureOIDCDirectory
	Tenant    string      `json:"tenant,omitempty"`
	Directory interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
}

// ProvisionedAccount Placeholder
type ProvisionedAccount struct {
	Account

	UId int `json:"uid,omitempty"` //ReadOnly
}

// GroupOnSystemQueryParams Query Parameters for Search GroupOnSystem
type GroupOnSystemQueryParams struct {
	Id                 []int64   `url:"id,omitempty"`
	GroupIds           []int64   `url:"group,omitempty"`
	ExcludeId          []int64   `url:"exclude,omitempty"`
	NotLinkedToGroupId []int64   `url:"notLinkedToGroup,omitempty"`
	NameContains       string    `url:"nameContains,omitempty"`
	CreatedAfter       time.Time `url:"createdAfter,omitempty"`
	CreatedBefore      time.Time `url:"createdBefore,omitempty"`
	ModifiedSince      time.Time `url:"modifiedSince,omitempty"`
	Q                  string    `url:"q,omitempty"`
	AdminnedById       int64     `url:"adminnedBy,omitempty"`
	OwnedById          int64     `url:"ownedBy,omitempty"`

	Additional *GroupOnSystemAdditionalQueryParams `url:"additional"`
}

// AddGroup extract id from given group and add to group ids to search for
func (g *GroupOnSystemQueryParams) AddGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.GroupIds = append(g.GroupIds, group.Self().ID)
		}
	}
}

// AddId extract id from given grouponsystem and add to id's to search for
func (g *GroupOnSystemQueryParams) AddId(gos *GroupOnSystem) {
	if gos != nil {
		if gos.Self().ID > 0 {
			g.Id = append(g.Id, gos.Self().ID)
		}
	}
}

// AddExcludeId extract id from given grouponsystem and add as group id to exclude in search
func (g *GroupOnSystemQueryParams) AddExcludeId(gos *GroupOnSystem) {
	if gos != nil {
		if gos.Self().ID > 0 {
			g.ExcludeId = append(g.ExcludeId, gos.Self().ID)
		}
	}
}

// AddNotLinkedToGroup extract id from given group and add as group id to exclude if linked to
func (g *GroupOnSystemQueryParams) AddNotLinkedToGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.NotLinkedToGroupId = append(g.NotLinkedToGroupId, group.Self().ID)
		}
	}
}

// SetAdminnedByGroup extract id from given group and add as AdminnedBy filter
func (g *GroupOnSystemQueryParams) SetAdminnedByGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.AdminnedById = group.Self().ID
		}
	} else {
		g.AdminnedById = 0
	}
}

// SetOwnedByGroup extract id from given group and add ad OwnedByGroup filter
func (g *GroupOnSystemQueryParams) SetOwnedByGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.OwnedById = group.Self().ID
		}
	} else {
		g.OwnedById = 0
	}
}

// GroupOnSystemAdditionalQueryParams AdditionalQueryParameters
type GroupOnSystemAdditionalQueryParams struct {
	Audit      bool `url:"audit"`
	ProvGroups bool `url:"provgroups"`
}

// EncodeValues Custom url encoder to convert bools to list
func (p GroupOnSystemAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}
