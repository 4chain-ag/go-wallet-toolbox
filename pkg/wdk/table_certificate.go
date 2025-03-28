package wdk

import (
	"time"
)

// TableCertificate represents a certificate with JSON tags
type TableCertificate struct {
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CertificateID      uint           `json:"certificateId"`
	UserID             int            `json:"userId"`
	Type               Base64String   `json:"type"`
	SerialNumber       Base64String   `json:"serialNumber"`
	Certifier          PubKeyHex      `json:"certifier"`
	Subject            PubKeyHex      `json:"subject"`
	Verifier           *PubKeyHex     `json:"verifier,omitempty"`
	RevocationOutpoint OutpointString `json:"revocationOutpoint"`
	Signature          HexString      `json:"signature"`
	IsDeleted          bool           `json:"isDeleted"`
}

// TableCertificateField represents a field related to a certificate
type TableCertificateField struct {
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	UserID        int          `json:"userId"`
	CertificateID uint         `json:"certificateId"`
	FieldName     string       `json:"fieldName"`
	FieldValue    string       `json:"fieldValue"`
	MasterKey     Base64String `json:"masterKey"`
}

// TableCertificateX extends TableCertificate with optional fields
type TableCertificateX struct {
	TableCertificate
	Fields []*TableCertificateField `json:"fields,omitempty"`
}

// WalletCertificate is a wallet certificate object
type WalletCertificate struct {
	Type               Base64String                                `json:"type"`
	Subject            PubKeyHex                                   `json:"subject"`
	SerialNumber       Base64String                                `json:"serialNumber"`
	Certifier          PubKeyHex                                   `json:"certifier"`
	RevocationOutpoint OutpointString                              `json:"revocationOutpoint"`
	Signature          HexString                                   `json:"signature"`
	Fields             map[CertificateFieldNameUnder50Bytes]string `json:"fields"`
}

// ListCertificatesResult is a response for ListCertificates action
type ListCertificatesResult struct {
	TotalCertificates PositiveIntegerOrZero `json:"totalCertificates"`
	Certificates      []*CertificateResult  `json:"certificates"`
}

// CertificateResult is a response with WalletCertificate
// extended with keyring and verifier
type CertificateResult struct {
	WalletCertificate
	Keyring  map[CertificateFieldNameUnder50Bytes]Base64String `json:"keyring"`
	Verifier string                                            `json:"verifier"`
}
