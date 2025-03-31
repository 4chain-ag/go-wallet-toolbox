package methodtests_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
)

func TestInsertCertificateAuth(t *testing.T) {
	// given:
	given := testabilities.Given(t)

	// and:
	activeStorage := given.GormProvider()

	t.Run("should insert a certificate for Alice", func(t *testing.T) {
		// and:
		certToInsert := &wdk.TableCertificateX{
			TableCertificate: wdk.TableCertificate{
				UserID:             testusers.Alice.ID,
				Type:               "exampleType",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				Subject:            "subjectPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
			},
			Fields: []*wdk.TableCertificateField{
				{
					UserID:     testusers.Alice.ID,
					FieldName:  "exampleField",
					FieldValue: "exampleValue",
					MasterKey:  "exampleMasterKey",
				},
			},
		}
		// when:
		_, err := activeStorage.InsertCertificateAuth(wdk.AuthID{
			UserID: &testusers.Alice.ID,
		}, certToInsert)

		// then:
		require.NoError(t, err)

		// check if id exists there
		certs, err := activeStorage.ListCertificates(wdk.AuthID{
			UserID: &testusers.Alice.ID,
		}, wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, wdk.PositiveIntegerOrZero(1), certs.TotalCertificates)
		require.Equal(t, 1, len(certs.Certificates[0].Fields))
		require.Equal(t, 1, len(certs.Certificates[0].Keyring))
		require.Equal(t, wdk.Base64String(certToInsert.Fields[0].FieldValue), certs.Certificates[0].Keyring["exampleField"])
		require.Equal(t, certToInsert.TableCertificate.Type, certs.Certificates[0].Type)
	})

	t.Run("should insert a certificate for Bob", func(t *testing.T) {
		// and:
		certToInsert := &wdk.TableCertificateX{
			TableCertificate: wdk.TableCertificate{
				UserID:             testusers.Bob.ID,
				Type:               "exampleType",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				Subject:            "subjectPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
			},
			Fields: []*wdk.TableCertificateField{
				{
					UserID:     testusers.Bob.ID,
					FieldName:  "exampleField",
					FieldValue: "exampleValue",
					MasterKey:  "exampleMasterKey",
				},
			},
		}
		// when:
		_, err := activeStorage.InsertCertificateAuth(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, certToInsert)

		// then:
		require.NoError(t, err)

		// when:
		certToInsert.Type = "exampleType2"
		_, err = activeStorage.InsertCertificateAuth(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, certToInsert)

		// then:
		require.NoError(t, err)

		certs, err := activeStorage.ListCertificates(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, wdk.PositiveIntegerOrZero(2), certs.TotalCertificates)
		require.Equal(t, 1, len(certs.Certificates[0].Fields))
		require.Equal(t, 1, len(certs.Certificates[1].Keyring))
		require.Equal(t, wdk.Base64String("exampleType"), certs.Certificates[0].Type)
		require.Equal(t, wdk.Base64String("exampleType2"), certs.Certificates[1].Type)
	})

	t.Run("should delete a certificate for Bob", func(t *testing.T) {
		// given:
		certs, err := activeStorage.ListCertificates(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, wdk.PositiveIntegerOrZero(2), certs.TotalCertificates)

		// when:
		err = activeStorage.RelinquishCertificate(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, wdk.RelinquishCertificateArgs{
			Type:         certs.Certificates[0].Type,
			SerialNumber: certs.Certificates[0].SerialNumber,
			Certifier:    certs.Certificates[0].Certifier,
		})

		// then:
		require.NoError(t, err)
		certs, err = activeStorage.ListCertificates(wdk.AuthID{
			UserID: &testusers.Bob.ID,
		}, wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, wdk.PositiveIntegerOrZero(1), certs.TotalCertificates)
	})
}

func TestInsertCertificateAuthFailure(t *testing.T) {
	// given:
	given := testabilities.Given(t)

	// and:
	activeStorage := given.GormProvider()

	t.Run("should fail to insert a certificate when no UserID provided in auth and when certificate UserID is different than authID", func(t *testing.T) {
		// when:
		_, err := activeStorage.InsertCertificateAuth(wdk.AuthID{}, &wdk.TableCertificateX{})

		// then:
		require.ErrorContains(t, err, "access is denied due to an authorization error")

		// and when:
		_, err = activeStorage.InsertCertificateAuth(wdk.AuthID{UserID: &testusers.Alice.ID}, &wdk.TableCertificateX{
			TableCertificate: wdk.TableCertificate{
				UserID: testusers.Bob.ID,
			},
		})

		// then:
		require.ErrorContains(t, err, "access is denied due to an authorization error")
	})

	t.Run("should fail to relinquish a certificate when no UserID provided in auth", func(t *testing.T) {
		// when:
		err := activeStorage.RelinquishCertificate(wdk.AuthID{}, wdk.RelinquishCertificateArgs{})

		// then:
		require.ErrorContains(t, err, "access is denied due to an authorization error")
	})

	t.Run("should fail to list certificate when no UserID provided in auth", func(t *testing.T) {
		// when:
		_, err := activeStorage.ListCertificates(wdk.AuthID{}, wdk.ListCertificatesArgs{})

		// then:
		require.ErrorContains(t, err, "access is denied due to an authorization error")
	})

	t.Run("should fail to delete certificate when no cert is found", func(t *testing.T) {
		// when:
		err := activeStorage.RelinquishCertificate(wdk.AuthID{
			UserID: &testusers.Alice.ID,
		}, wdk.RelinquishCertificateArgs{
			Type:         "not-type",
			SerialNumber: "not-serial",
			Certifier:    "not-certifier",
		})

		// then:
		require.ErrorContains(t, err, "failed to delete certificate model: certificate not found")
	})
}
