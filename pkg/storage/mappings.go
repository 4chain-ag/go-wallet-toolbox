package storage

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func tableCertificateXFieldsToModelFields(userID int) func(*wdk.TableCertificateField) *models.CertificateField {
	return func(t *wdk.TableCertificateField) *models.CertificateField {
		return &models.CertificateField{
			FieldName:  t.FieldName,
			FieldValue: t.FieldValue,
			MasterKey:  string(t.MasterKey),
			UserID:     userID,
		}
	}
}

func listCertificatesArgsToOptions(args wdk.ListCertificatesArgs) repo.ListCertificatesOptions {
	opts := repo.ListCertificatesOptions{
		Limit:  args.Limit,
		Offset: args.Offset,
	}
	types := args.Types
	certifiers := args.Certifiers

	if args.Partial != nil {
		opts.SerialNumber = args.Partial.SerialNumber
		opts.Subject = args.Partial.Subject
		opts.RevocationOutpoint = args.Partial.RevocationOutpoint
		opts.Signature = args.Partial.Signature

		if args.Partial.Type != nil {
			types = append(types, *args.Partial.Type)
		}

		if args.Partial.Certifier != nil {
			certifiers = append(certifiers, *args.Partial.Certifier)
		}
	}

	opts.Types = types
	opts.Certifiers = certifiers

	return opts
}

func certModelToResult(model *models.Certificate) *wdk.CertificateResult {
	return &wdk.CertificateResult{
		Verifier: model.Verifier,
		Keyring:  certificateModelFieldsToKeyringResult(model.CertificateFields),
		WalletCertificate: wdk.WalletCertificate{
			Type:               wdk.Base64String(model.Type),
			Subject:            wdk.PubKeyHex(model.Subject),
			SerialNumber:       wdk.Base64String(model.SerialNumber),
			Certifier:          wdk.PubKeyHex(model.Certifier),
			RevocationOutpoint: wdk.OutpointString(model.RevocationOutpoint),
			Signature:          wdk.HexString(model.Signature),
			Fields:             certificateModelFieldsToFieldsResult(model.CertificateFields),
		},
	}
}

func certificateModelFieldsToKeyringResult(fields []*models.CertificateField) map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String {
	result := make(map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String, len(fields))
	for _, field := range fields {
		result[wdk.CertificateFieldNameUnder50Bytes(field.FieldName)] = wdk.Base64String(field.FieldValue)
	}

	return result
}

func certificateModelFieldsToFieldsResult(fields []*models.CertificateField) map[wdk.CertificateFieldNameUnder50Bytes]string {
	result := make(map[wdk.CertificateFieldNameUnder50Bytes]string, len(fields))
	for _, field := range fields {
		result[wdk.CertificateFieldNameUnder50Bytes(field.FieldName)] = field.FieldValue
	}

	return result
}
