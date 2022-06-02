package model

type Client struct {
	Linkable

	UUID           string   `json:"uuid,omitempty"`
	ClientId       string   `json:"clientId,omitempty"`
	SSOApplication bool     `json:"ssoApplication,omitempty"`
	Name           string   `json:"name"`
	Scopes         []string `json:"scopes,omitempty"`
	URL            string   `json:"url,omitempty"`
}
