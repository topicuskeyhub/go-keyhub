package model

import (
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

type GroupOnSystemType string
type ProvisionedSystemType string

// Section: provisioning
type GroupOnSystemList struct {
	Items []GroupOnSystem `json:"items"`
}
type GroupOnSystem struct {
	Linkable

	Type              GroupOnSystemType  `json:"type,omitempty"`
	NameInSystem      string             `json:"nameInSystem,omitempty"`
	ShortNameInSystem string             `json:"shortNameInSystem,omitempty"` // Read Only
	DisplayName       string             `json:"displayName,omitempty"`
	System            *ProvisionedSystem `json:"system,omitempty"`
	Owner             *Group             `json:"owner,omitempty"`

	AdditionalObject *GroupOnSystemAdditionalObject `json:"additionalObject,omitempty"`
}

type GroupOnSystemAdditionalObject struct {
	ProvGroups struct {
		Items []ProvisioningGroup `json:"items,omitempty"`
	} `json:"provgroups,omitempty"`
}

func (g *GroupOnSystem) AddProvGroup(group ...ProvisioningGroup) {

	if g.AdditionalObject == nil {
		g.AdditionalObject = &GroupOnSystemAdditionalObject{}
	}
	g.AdditionalObject.ProvGroups.Items = append(g.AdditionalObject.ProvGroups.Items, group...)
}

func (g *GroupOnSystem) SetType(typename GroupOnSystemType) {
	g.Type = typename
}

func (g *GroupOnSystem) SetName(name string) {
	g.DisplayName = name
	g.NameInSystem = slug.Make(name)
}

func NewGroupOnSystem() *GroupOnSystem {
	return &GroupOnSystem{
		Linkable: Linkable{DType: "provisioning.GroupOnSystem"},
	}
}

// Section: provisioning

func NewProvisionedSystem() *ProvisionedSystem {

	ps := &ProvisionedSystem{}
	ps.Linkable = Linkable{DType: "provisioning.ProvisionedSystem"}

	return ps
}

type ProvisionedSystemPrimer struct {
	Linkable
	Active bool                  `json:"active,omitempty"`
	UUID   string                `json:"UUID,omitempty"`
	Name   string                `json:"name,omitempty"`
	Type   ProvisionedSystemType `json:"type,omitempty"`
}

type ProvisionedSystem struct {
	ProvisionedSystemPrimer

	ProvisionedAbstract

	AccountCount           int
	UsernamePrefix         string `json:"usernamePrefix,omitempty"`
	TechnicalAdministrator *Group `json:"technicalAdministrator,omitempty"`
	ExternalUUID           string `json:"externalUUID,omitempty"`
}

func (p *ProvisionedSystem) AsPrimer() *ProvisionedSystem {
	system := &ProvisionedSystem{}
	system.ProvisionedSystemPrimer = p.ProvisionedSystemPrimer
	return system
}

type ProvisionedSystemList struct {
	Items []ProvisionedSystem `json:"items,omitempty"`
}

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

type ProvisionedAD struct {
	//provisioning.ProvisionedAD
	SamAccountNameScheme string `json:"samAccountNameScheme,omitempty"`
	//"OMIT"
	//"TRUNCATE"
	//"TRANSFER"
	//"TRANSFER_TRUNCATE"
	//"USERNAME"
}

type ProvisionedLDAPDirectory struct {
	//provisioning.ProvisionedLDAPDirectory
	Directory interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
	GroupDN   string      `json:"groupDN,omitempty"`
}
type ProvisionedAzureTenant struct {
	//provisioning.ProvisionedAzureTenant
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	Tenant       string `json:"tenant,omitempty"`
	IdpDomain    string `json:"idpDomain,omitempty"`
}

type ProvisionedInternalLDAP struct {
	//provisioning.ProvisionedInternalLDAP
	Client interface{} `json:"client,omitempty"` // client_ldapClient
}

type ProvisionedLDAP struct {
	//provisioning.ProvisionedLDAP
	GID           int         `json:"gid,omitempty"`
	HashingScheme string      `json:"hashingScheme"`
	Numbering     interface{} `json:"numbering,omitempty"` // provisioning_ProvisionNumberSequence
}

type ProvisionedAzureSyncLDAPDirectory struct {
	//provisioning.ProvisionedAzureSyncLDAPDirectory
	ClientId     string      `json:"clientId,omitempty"`
	ClientSecret string      `json:"clientSecret,omitempty"`
	Tenant       string      `json:"tenant,omitempty"`
	Directory    interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
}

type ProvisionedAzureOIDCDirectory struct {
	//provisioning.ProvisionedAzureOIDCDirectory
	Tenant    string      `json:"tenant,omitempty"`
	Directory interface{} `json:"directory,omitempty"` // directory_AccountDirectoryPrimer
}

type ProvisionedAccount struct {
	Account

	UId int `json:"uid,omitempty"` //ReadOnly
}

type Certificate struct {
}

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

func (g *GroupOnSystemQueryParams) AddGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.GroupIds = append(g.GroupIds, group.Self().ID)
		}
	}
}

func (g *GroupOnSystemQueryParams) AddId(gos *GroupOnSystem) {
	if gos != nil {
		if gos.Self().ID > 0 {
			g.Id = append(g.Id, gos.Self().ID)
		}
	}
}
func (g *GroupOnSystemQueryParams) AddExcludeId(gos *GroupOnSystem) {
	if gos != nil {
		if gos.Self().ID > 0 {
			g.ExcludeId = append(g.ExcludeId, gos.Self().ID)
		}
	}
}
func (g *GroupOnSystemQueryParams) AddNotLinkedToGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.NotLinkedToGroupId = append(g.NotLinkedToGroupId, group.Self().ID)
		}
	}
}
func (g *GroupOnSystemQueryParams) SetAdminnedByGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.AdminnedById = group.Self().ID
		}
	} else {
		g.AdminnedById = 0
	}
}

func (g *GroupOnSystemQueryParams) SetOwnedByGroup(group *Group) {
	if group != nil {
		if group.Self().ID > 0 {
			g.OwnedById = group.Self().ID
		}
	} else {
		g.OwnedById = 0
	}
}

type GroupOnSystemAdditionalQueryParams struct {
	Audit      bool `url:"audit"`
	ProvGroups bool `url:"provgroups"`
}

func (p GroupOnSystemAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}
