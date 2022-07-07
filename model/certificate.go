package model

import "time"

type Certificate struct {
	Linkable

	UUID              string    `json:"uuid,omitempty"`
	Global            bool      `json:"global"`
	Alias             string    `json:"alias,omitempty"`
	SubjectDN         string    `json:"subjectDN,omitempty"`
	FingerprintSha1   string    `json:"fingerprintSha1,omitempty"`
	FingerprintSha256 string    `json:"fingerprintSha256,omitempty"`
	CertificateData   []byte    `json:"certificateData"`
	Expiration        time.Time `json:"expiration,omitempty"`
}
