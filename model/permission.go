package model

const (
	PERM_ACCOUNTS_QUERY                        = "ACCOUNTS_QUERY"
	PERM_ACCOUNTS_REMOVE                       = "ACCOUNTS_REMOVE"
	PERM_GROUPONSYSTEM_CREATE                  = "GROUPONSYSTEM_CREATE"
	PERM_GROUPS_CREATE                         = "GROUPS_CREATE"
	PERM_GROUPS_VAULT_ACCESS_AFTER_CREATE      = "GROUPS_VAULT_ACCESS_AFTER_CREATE"
	PERM_GROUPS_GRANT_PERMISSIONS_AFTER_CREATE = "GROUPS_GRANT_PERMISSIONS_AFTER_CREATE"
	PERM_GROUPS_QUERY                          = "GROUPS_QUERY"
	PERM_GROUP_FULL_VAULT_ACCESS               = "GROUP_FULL_VAULT_ACCESS"
	PERM_GROUP_READ_CONTENTS                   = "GROUP_READ_CONTENTS"
	PERM_GROUP_SET_AUTHORIZATION               = "GROUP_SET_AUTHORIZATION"
	PERM_CLIENTS_CREATE                        = "CLIENTS_CREATE"
	PERM_CLIENTS_QUERY                         = "CLIENTS_QUERY"

	PERM_OPERATION_READ   = "READ"
	PERM_OPERATION_UPDATE = "UPDATE"
	PERM_OPERATION_DELETE = "DELETE"
)

type Permission struct {
	Full       string   `json:"full"`
	Type       string   `json:"type"`
	Operations []string `json:"operations"`
	Instances  []string `json:"instances"`
}
