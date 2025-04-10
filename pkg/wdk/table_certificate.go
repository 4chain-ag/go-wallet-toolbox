package wdk

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
)

// TableCertificate represents a certificate with JSON tags
type TableCertificate struct {
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
	CertificateID      uint                      `json:"certificateId"`
	UserID             int                       `json:"userId"`
	Type               primitives.Base64String   `json:"type"`
	SerialNumber       primitives.Base64String   `json:"serialNumber"`
	Certifier          primitives.PubKeyHex      `json:"certifier"`
	Subject            primitives.PubKeyHex      `json:"subject"`
	Verifier           *primitives.PubKeyHex     `json:"verifier,omitempty"`
	RevocationOutpoint primitives.OutpointString `json:"revocationOutpoint"`
	Signature          primitives.HexString      `json:"signature"`
	IsDeleted          bool                      `json:"isDeleted"`
}

// TableCertificateField represents a field related to a certificate
type TableCertificateField struct {
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	UserID        int                     `json:"userId"`
	CertificateID uint                    `json:"certificateId"`
	FieldName     string                  `json:"fieldName"`
	FieldValue    string                  `json:"fieldValue"`
	MasterKey     primitives.Base64String `json:"masterKey"`
}

// TableCertificateX extends TableCertificate with optional fields
type TableCertificateX struct {
	TableCertificate
	Fields []*TableCertificateField `json:"fields,omitempty"`
}

// WalletCertificate is a wallet certificate object
type WalletCertificate struct {
	Type               primitives.Base64String                                `json:"type"`
	Subject            primitives.PubKeyHex                                   `json:"subject"`
	SerialNumber       primitives.Base64String                                `json:"serialNumber"`
	Certifier          primitives.PubKeyHex                                   `json:"certifier"`
	RevocationOutpoint primitives.OutpointString                              `json:"revocationOutpoint"`
	Signature          primitives.HexString                                   `json:"signature"`
	Fields             map[primitives.CertificateFieldNameUnder50Bytes]string `json:"fields"`
}

// ListCertificatesResult is a response for ListCertificates action
type ListCertificatesResult struct {
	TotalCertificates primitives.PositiveIntegerOrZero `json:"totalCertificates"`
	Certificates      []*CertificateResult             `json:"certificates"`
}

// CertificateResult is a response with WalletCertificate
// extended with keyring and verifier
type CertificateResult struct {
	WalletCertificate
	Keyring  map[primitives.CertificateFieldNameUnder50Bytes]primitives.Base64String `json:"keyring"`
	Verifier string                                                                  `json:"verifier"`
}
