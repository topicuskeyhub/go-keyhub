package model

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

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

type ServiceAccountQueryParams struct {
	UUID          string                               `url:"uuid,omitempty"`
	CreatedAfter  time.Time                            `url:"createdAfter,omitempty" layout:"2006-01-02T15:04:05Z"`
	CreatedBefore time.Time                            `url:"createdBefore,omitempty" layout:"2006-01-02T15:04:05Z"`
	ModifiedSince time.Time                            `url:"createdBefore,omitempty" layout:"2006-01-02T15:04:05Z"`
	Additional    *ServiceAccountAdditionalQueryParams `url:"additional,omitempty"`
	Exclude       []int64                              `url:"exclude,omitempty"`
	id            []int64                              `url:"id,omitempty"`
}

type ServiceAccountAdditionalQueryParams struct {
	Audit   bool `url:"audit"`
	Groups  bool `url:"groups"`
	Secrets bool `url:"secrets"`
}

// EncodeValues Custom url encoder to convert bools to list
func (p ServiceAccountAdditionalQueryParams) EncodeValues(key string, v *url.Values) error {
	return additionalQueryParamsUrlEncoder(p, key, v)
}

func NewServiceAccount() *ServiceAccount {

	return &ServiceAccount{
		ServiceAccountPrimer: ServiceAccountPrimer{
			Linkable: Linkable{DType: "serviceaccount.ServiceAccount"},
		},
		PasswordRotation: SA_PASSWORD_ROTATION_MANUAL,
	}
}
