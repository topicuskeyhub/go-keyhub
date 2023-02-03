package model

import "github.com/google/uuid"

const (
	SA_PASSWORD_ROTATION_MANUAL     SAPasswordRotation = "MANUAL"
	SA_PASSWORD_ROTATION_MANUAL_SIV SAPasswordRotation = "MANUAL_STORED_IN_VAULT"
	SA_PASSWORD_ROTATION_DAILY      SAPasswordRotation = "DAILY"
)

type SAPasswordRotation string

type ServiceAccountList struct {
	Items []ServiceAccount `json:"items"`
}

type ServiceAccountPrimer struct {
	//#/components/schemas/serviceaccount_ServiceAccountPrimer
	//serviceaccount_ServiceAccount

	Linkable

	Active   bool               `json:"active"`
	UUID     uuid.UUID          `json:"uuid,omitempty"`
	Name     string             `json:"name"`
	Username string             `json:"username"`
	System   *ProvisionedSystem `json:"system"`
}

type ServiceAccount struct {
	ServiceAccountPrimer

	TechnicalAdministrator *Group             `json:"technicalAdministrator,omitempty"`
	PasswordRotation       SAPasswordRotation `json:"passwordRotation"`
	Description            string             `json:"description,omitempty"`
	Password               *VaultRecord       `json:"password,omitempty"`
}

// AsPrimer Return ServiceAccount with only Primer data
func (s *ServiceAccount) AsPrimer() *ServiceAccount {
	serviceAccount := &ServiceAccount{}
	serviceAccount.ServiceAccountPrimer = s.ServiceAccountPrimer
	return serviceAccount
}

// ToPrimer Convert to serviceAccountPrimer
func (s *ServiceAccount) ToPrimer() *ServiceAccountPrimer {
	serviceAccountPrimer := s.ServiceAccountPrimer
	return &serviceAccountPrimer
}

type ServiceAccountGroupList struct {
	Items []ServiceAccountGroup `json:"items"`
}

type ServiceAccountGroup struct {
	GroupOnSystem
}

func NewServiceAccount() *ServiceAccount {

	return &ServiceAccount{
		ServiceAccountPrimer: ServiceAccountPrimer{
			Linkable: Linkable{DType: "serviceaccount.ServiceAccount"},
		},
		PasswordRotation: SA_PASSWORD_ROTATION_MANUAL,
	}
}
