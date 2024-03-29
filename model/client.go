package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	CLIENT_TYPE_OAUTH2 ClientApplicationType = "OAUTH2"
	CLIENT_TYPE_SAML2  ClientApplicationType = "SAML2"
	CLIENT_TYPE_LDAP   ClientApplicationType = "LDAP"

	CLIENT_SCOPE_PROFILE        = "profile"
	CLIENT_SCOPE_MANAGE_ACCOUNT = "manage_account"
	CLIENT_SCOPE_PROVISIONING   = "provisioning"
	CLIENT_SCOPE_ACCESS_VAULT   = "access_vault"
	CLIENT_SCOPE_GROUP_ADMIN    = "group_admin"
	CLIENT_SCOPE_GLOBAL_ADMIN   = "global_admin"

	CLIENT_SUBJECT_FORMAT_ID       = "ID"
	CLIENT_SUBJECT_FORMAT_UPN      = "UPN"
	CLIENT_SUBJECT_FORMAT_USERNAME = "USERNAME"
	CLIENT_SUBJECT_FORMAT_EMAIL    = "EMAIL"

	CLIENT_PERM_ACCOUNTS_QUERY                        Oauth2ClientPermissionValue = "ACCOUNTS_QUERY"
	CLIENT_PERM_ACCOUNTS_REMOVE                       Oauth2ClientPermissionValue = "ACCOUNTS_REMOVE"
	CLIENT_PERM_GROUPONSYSTEM_CREATE                  Oauth2ClientPermissionValue = "GROUPONSYSTEM_CREATE"
	CLIENT_PERM_GROUPS_CREATE                         Oauth2ClientPermissionValue = "GROUPS_CREATE"
	CLIENT_PERM_GROUPS_VAULT_ACCESS_AFTER_CREATE      Oauth2ClientPermissionValue = "GROUPS_VAULT_ACCESS_AFTER_CREATE"
	CLIENT_PERM_GROUPS_GRANT_PERMISSIONS_AFTER_CREATE Oauth2ClientPermissionValue = "GROUPS_GRANT_PERMISSIONS_AFTER_CREATE"
	CLIENT_PERM_GROUPS_QUERY                          Oauth2ClientPermissionValue = "GROUPS_QUERY"
	CLIENT_PERM_GROUP_FULL_VAULT_ACCESS               Oauth2ClientPermissionValue = "GROUP_FULL_VAULT_ACCESS"
	CLIENT_PERM_GROUP_READ_CONTENTS                   Oauth2ClientPermissionValue = "GROUP_READ_CONTENTS"
	CLIENT_PERM_GROUP_SET_AUTHORIZATION               Oauth2ClientPermissionValue = "GROUP_SET_AUTHORIZATION"
	CLIENT_PERM_CLIENTS_CREATE                        Oauth2ClientPermissionValue = "CLIENTS_CREATE"
	CLIENT_PERM_CLIENTS_QUERY                         Oauth2ClientPermissionValue = "CLIENTS_QUERY"
)

func NewOAuth2ClientApplication(name string, managerGroup *Group) *ClientApplication {
	return NewClientApplication(name, managerGroup, CLIENT_TYPE_OAUTH2)
}

func NewLdapClientApplication(name string, managerGroup *Group) *ClientApplication {
	return NewClientApplication(name, managerGroup, CLIENT_TYPE_LDAP)
}

func NewSaml2ClientApplication(name string, managerGroup *Group) *ClientApplication {
	return NewClientApplication(name, managerGroup, CLIENT_TYPE_SAML2)
}

func NewClientApplication(name string, managerGroup *Group, clienttype ClientApplicationType) *ClientApplication {

	cl := &ClientApplication{}
	cl.Name = name
	cl.Type = clienttype
	cl.Owner = managerGroup
	cl.TechnicalAdministrator = managerGroup

	cl.Attributes = make(map[string]string)

	switch clienttype {
	case CLIENT_TYPE_OAUTH2:
		cl.DType = "client.OAuth2Client"
		cl.Confidential = true
		cl.SetOAuth2SSO()
	case CLIENT_TYPE_LDAP:
		cl.DType = "client.LdapClient"
		cl.ClientCertificate = nil
	case CLIENT_TYPE_SAML2:
		cl.DType = "client.Saml2Client"
		cl.SSOApplication = true
		cl.Scopes = append(cl.Scopes, CLIENT_SCOPE_PROFILE)
	}

	return cl

}

/** Enums **/
type ClientApplicationType string
type Oauth2ClientPermissionValue string

type oauth2 struct {
	Confidential         bool   `json:"confidential,omitempty"`         // Oauth
	UseClientCredentials bool   `json:"useClientCredentials,omitempty"` // OAuth Server2Server
	CallbackURI          string `json:"callbackURI,omitempty"`          // OAuth SSO
	InitiateLoginURI     string `json:"initiateLoginURI,omitempty"`     // OAuth SSO
	IdTokenClaims        string `json:"idTokenClaims,omitempty"`        // OAuth SSO
	ShowLandingPage      bool   `json:"showLandingPage,omitempty"`      // OAuth SSO

}

type ldap struct {
	BindDN              string       `json:"bindDN,omitempty"`              // LDAP
	UsedForProvisioning bool         `json:"usedForProvisioning,omitempty"` // LDAP
	ClientCertificate   *Certificate `json:"clientCertificate,omitempty"`   // LDAP
}

type saml2 struct {
	MetadataUrl   string `json:"metadataUrl,omitempty"`   // Saml
	Metadata      string `json:"metadata,omitempty"`      // Saml
	SubjectFormat string `json:"subjectFormat,omitempty"` // Saml
	Segments      string `json:"segments,omitempty"`      // Saml
}

type ClientApplicationPrimer struct {
	Linkable

	UUID string                `json:"uuid,omitempty"`
	Name string                `json:"name"`
	Type ClientApplicationType `json:"type,omitempty"`

	// Shared
	ClientId       string   `json:"clientId,omitempty"`       // OAuth: the clientId, Saml: The ClientApplication identifier
	Scopes         []string `json:"scopes,omitempty"`         // Oauth SSO: required, Saml/Ldap: fixed to profile
	SSOApplication bool     `json:"ssoApplication,omitempty"` // Oauth SSO + Saml

}

type ClientApplication struct {
	ClientApplicationPrimer

	oauth2
	ldap
	saml2

	URL                    string                   `json:"url,omitempty"`
	Permissions            []Permission             `json:"permissions,omitempty"`
	AdditionalObjects      *ClientAdditionalObjects `json:"additionalObjects,omitempty"`
	LastModifiedAt         time.Time                `json:"lastModifiedAt,omitempty"`
	Owner                  *Group                   `json:"owner,omitempty"`
	TechnicalAdministrator *Group                   `json:"technicalAdministrator,omitempty"`
	DebugMode              bool                     `json:"debugMode,omitempty"`
	AccountPermissions     []Permission             `json:"accountPermissions,omitempty"`
	Attributes             map[string]string        `json:"attributes,omitempty"` // OAuth SSO + Saml

}

func (c *ClientApplication) IsOAuth2Server2Server() bool {
	return c.Type == CLIENT_TYPE_OAUTH2 && c.UseClientCredentials == true
}

func (c *ClientApplication) IsOAuth2SSO() bool {
	return c.Type == CLIENT_TYPE_OAUTH2 && c.SSOApplication == true
}

func (c *ClientApplication) IsOAuth2() bool {
	return c.Type == CLIENT_TYPE_OAUTH2
}

func (c *ClientApplication) IsSAML2() bool {
	return c.Type == CLIENT_TYPE_SAML2
}

func (c *ClientApplication) IsLDAP() bool {
	return c.Type == CLIENT_TYPE_LDAP
}

func (c *ClientApplication) AsPrimer() *ClientApplication {
	x := ClientApplication{}
	x.ClientApplicationPrimer = c.ClientApplicationPrimer
	return &x
}

func (c *ClientApplication) SetOAuth2Server2Server() error {
	if !c.IsOAuth2() {
		return fmt.Errorf("client type is not an oauth2 type, can not set to server2server")
	}
	c.UseClientCredentials = true
	c.SSOApplication = false
	return nil
}

func (c *ClientApplication) SetOAuth2SSO() error {
	if !c.IsOAuth2() {
		return fmt.Errorf("client type is not an oauth2 type, can not set to server2server")
	}
	c.UseClientCredentials = false
	c.SSOApplication = true
	c.AddScope(CLIENT_SCOPE_PROFILE)
	return nil
}

func (c *ClientApplication) AddAttribute(name, script string) error {
	if c.IsOAuth2SSO() == false && c.IsSAML2() == false {
		return fmt.Errorf("current client type does not support Attributes. Type: %s", c.Type)
	}

	name = strings.TrimSpace(name)
	if _, exists := c.Attributes[name]; exists {
		return fmt.Errorf("an attribute named '%s' allready exists", name)
	}

	c.Attributes[name] = script

	return nil
}

func (c *ClientApplication) RemoveAttribute(name, script string) error {
	if c.IsOAuth2SSO() == false && c.IsSAML2() == false {
		return fmt.Errorf("current client type does not support Attributes")
	}

	name = strings.TrimSpace(name)
	if _, exists := c.Attributes[name]; exists == false {
		return fmt.Errorf("an attribute named '%s' does not exist", name)
	}

	delete(c.Attributes, name)

	return nil
}

func (c *ClientApplication) RemoveScope(scope string) error {
	if !c.IsOAuth2SSO() {
		return fmt.Errorf("current client type does not support Scopes")
	}

	// Scopes should contain max 6 items, so this loop should not be that awfull.
	var newScopes []string
	for _, curScope := range c.Scopes {
		if curScope != scope {
			newScopes = append(newScopes, curScope)
		}
	}
	c.Scopes = newScopes
	return nil
}

func (c *ClientApplication) AddScope(scope string) error {
	if !c.IsOAuth2SSO() {
		return fmt.Errorf("current client type '%s' does not support Scopes: %v", c.Type, c.IsOAuth2SSO())
	}

	switch strings.ToLower(scope) {
	case CLIENT_SCOPE_PROFILE:
	case CLIENT_SCOPE_MANAGE_ACCOUNT:
	case CLIENT_SCOPE_PROVISIONING:
	case CLIENT_SCOPE_ACCESS_VAULT:
	case CLIENT_SCOPE_GROUP_ADMIN:
	case CLIENT_SCOPE_GLOBAL_ADMIN:
		// All above are correct
	default:
		return fmt.Errorf("'%s' is not a valid scope", scope)
	}

	for _, existingScope := range c.Scopes {
		if scope == existingScope {
			// Scope all ready is enabled
			return nil
		}
	}

	c.Scopes = append(c.Scopes, scope)
	return nil
}

func (c *ClientApplication) GetSecret() (string, error) {

	if c.AdditionalObjects == nil || c.AdditionalObjects.Secret == nil {
		return "", fmt.Errorf("secret is not available")
	}

	return c.AdditionalObjects.Secret.GeneratedSecret, nil
}
func (c *ClientApplication) GetSecretOrNil() *string {
	if c.AdditionalObjects == nil || c.AdditionalObjects.Secret == nil {
		return nil
	}
	return &c.AdditionalObjects.Secret.GeneratedSecret
}

type ClientList struct {
	Items []ClientApplication `json:"items"`
}

type generatedSecret struct {
	DType           string `json:"$type"`
	GeneratedSecret string `json:"generatedSecret"`
}

type ClientQueryParams struct {
	UUID       string   `url:"uuid,omitempty"`
	Additional []string `url:"additional,omitempty"`
}

type ClientAdditionalObjects struct {
	Audit  *AuditAdditionalObject `json:"audit,omitempty"`
	Secret *generatedSecret       `json:"secret,omitempty"`
}

type OAuth2ClientPermissionList struct {
	DType string                    `json:"$type,omitempty"`
	Items []*OAuth2ClientPermission `json:"items"`
}

type OAuth2ClientPermission struct {
	DType     string                      `json:"$type,omitempty"`
	Value     Oauth2ClientPermissionValue `json:"value,omitempty"`
	ForSystem *ClientApplication          `json:"forSystem,omitempty"`
	ForGroup  *Group                      `json:"forGroup,omitempty"`
}

func NewOAuth2ClientPermission(perm Oauth2ClientPermissionValue, group *Group) *OAuth2ClientPermission {

	cp := &OAuth2ClientPermission{}
	cp.DType = "client.OAuth2ClientPermission"
	cp.Value = perm
	cp.ForGroup = group

	return cp
}
