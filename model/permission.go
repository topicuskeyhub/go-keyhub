package model

const (
	PERM_OPERATION_CREATE PermissionOperationName = "CREATE"
	PERM_OPERATION_READ   PermissionOperationName = "READ"
	PERM_OPERATION_UPDATE PermissionOperationName = "UPDATE"
	PERM_OPERATION_DELETE PermissionOperationName = "DELETE"
)

type Permission struct {
	Full       string                    `json:"full"`
	Type       string                    `json:"type"`
	Operations []PermissionOperationName `json:"operations"`
	Instances  []string                  `json:"instances"`
}

type PermissionOperationName string

type ClientPermissionsWithClient struct {
	DType string                        `json:"$type,omitempty"`
	Items []*ClientPermissionWithClient `json:"items"`
}

type ClientPermissionWithClient struct {
	DType  string                      `json:"$type,omitempty"`
	Value  Oauth2ClientPermissionValue `json:"value,omitempty"`
	Client *ClientApplication          `json:"client,omitempty"`
}

func NewClientPermissionWithClient(perm Oauth2ClientPermissionValue, client *ClientApplication) *ClientPermissionWithClient {

	cp := &ClientPermissionWithClient{}
	cp.DType = "client.OAuth2ClientPermissionWithClient"
	cp.Value = perm
	cp.Client = client.AsPrimer()

	return cp
}
