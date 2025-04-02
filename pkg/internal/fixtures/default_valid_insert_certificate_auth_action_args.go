package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

const (
	// Base64-encoded string (original: "exampleType")
	TypeField = "ZXhhbXBsZVR5cGU="

	// Base64-encoded string (original: "serial123")
	SerialNumber = "c2VyaWFsMTIz"

	// PubKeyHex (33-byte compressed public key)
	Certifier = "02c123eabcdeff1234567890abcdef1234567890abcdef1234567890abcdef1234"

	// PubKeyHex (33-byte compressed public key)
	SubjectPubKey = "02c123eabcdeff1234567890abcdef1234567890abcdef1234567890abcdef5678"

	// OutpointString (format: txid:vout)
	RevocationOutpoint = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890.0"

	// HexString (64-byte signature)
	Signature = "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
)

func DefaultInsertCertAuth(userID int) *wdk.TableCertificateX {
	return &wdk.TableCertificateX{
		TableCertificate: wdk.TableCertificate{
			UserID:             userID,
			Type:               TypeField,
			SerialNumber:       SerialNumber,
			Certifier:          Certifier,
			Subject:            SubjectPubKey,
			RevocationOutpoint: RevocationOutpoint,
			Signature:          Signature,
		},
		Fields: []*wdk.TableCertificateField{
			{
				UserID:     userID,
				FieldName:  "exampleField",
				FieldValue: "exampleValue",
				MasterKey:  "exampleMasterKey",
			},
		},
	}
}
