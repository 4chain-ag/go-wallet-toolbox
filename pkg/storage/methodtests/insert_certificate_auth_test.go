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

	// and:
	expectedResult := &wdk.ListCertificatesResult{
		Certificates:      make([]*wdk.CertificateResult, 0),
		TotalCertificates: wdk.PositiveIntegerOrZero(0),
	}

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

		// and:
		expectedResult.TotalCertificates = wdk.PositiveIntegerOrZero(1)
		expectedResult.Certificates = []*wdk.CertificateResult{{
			Verifier: "",
			WalletCertificate: wdk.WalletCertificate{
				Type:               "exampleType",
				Subject:            "subjectPubKey",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
				Fields: map[wdk.CertificateFieldNameUnder50Bytes]string{
					"exampleField": "exampleValue",
				},
			},
			Keyring: map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String{
				"exampleField": "exampleValue",
			},
		}}

		// when: insert certificate for Alice
		_, err := activeStorage.InsertCertificateAuth(testusers.Alice.AuthID(), certToInsert)

		// then:
		require.NoError(t, err)

		// when: listing certificates
		certs, err := activeStorage.ListCertificates(testusers.Alice.AuthID(), wdk.ListCertificatesArgs{})

		// then:
		require.NoError(t, err)
		require.Equal(t, expectedResult, certs)
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

		// and:
		expectedResult.TotalCertificates = wdk.PositiveIntegerOrZero(2)
		expectedResult.Certificates = []*wdk.CertificateResult{{
			Verifier: "",
			WalletCertificate: wdk.WalletCertificate{
				Type:               "exampleType2",
				Subject:            "subjectPubKey",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
				Fields: map[wdk.CertificateFieldNameUnder50Bytes]string{
					"exampleField": "exampleValue",
				},
			},
			Keyring: map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String{
				"exampleField": "exampleValue",
			},
		}, {
			Verifier: "",
			WalletCertificate: wdk.WalletCertificate{
				Type:               "exampleType",
				Subject:            "subjectPubKey",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
				Fields: map[wdk.CertificateFieldNameUnder50Bytes]string{
					"exampleField": "exampleValue",
				},
			},
			Keyring: map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String{
				"exampleField": "exampleValue",
			},
		}}

		// when: insert certificate for Bob
		_, err := activeStorage.InsertCertificateAuth(testusers.Bob.AuthID(), certToInsert)

		// then:
		require.NoError(t, err)

		// when:
		certToInsert.Type = "exampleType2"
		_, err = activeStorage.InsertCertificateAuth(testusers.Bob.AuthID(), certToInsert)

		// then:
		require.NoError(t, err)

		// when: listing certificates
		certs, err := activeStorage.ListCertificates(testusers.Bob.AuthID(), wdk.ListCertificatesArgs{})

		// then:
		require.NoError(t, err)
		require.Equal(t, expectedResult, certs)
	})

	t.Run("should delete a certificate for Bob", func(t *testing.T) {
		// given:
		certs, err := activeStorage.ListCertificates(testusers.Bob.AuthID(), wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, wdk.PositiveIntegerOrZero(2), certs.TotalCertificates)

		// and:
		expectedResult.TotalCertificates = wdk.PositiveIntegerOrZero(1)
		expectedResult.Certificates = []*wdk.CertificateResult{{
			Verifier: "",
			WalletCertificate: wdk.WalletCertificate{
				Type:               "exampleType",
				Subject:            "subjectPubKey",
				SerialNumber:       "exampleSerialNumber",
				Certifier:          "certifierPubKey",
				RevocationOutpoint: "outpointString",
				Signature:          "signatureHex",
				Fields: map[wdk.CertificateFieldNameUnder50Bytes]string{
					"exampleField": "exampleValue",
				},
			},
			Keyring: map[wdk.CertificateFieldNameUnder50Bytes]wdk.Base64String{
				"exampleField": "exampleValue",
			},
		}}

		// when:
		err = activeStorage.RelinquishCertificate(testusers.Bob.AuthID(), wdk.RelinquishCertificateArgs{
			Type:         "exampleType2",
			SerialNumber: "exampleSerialNumber",
			Certifier:    "certifierPubKey",
		})

		// then:
		require.NoError(t, err)
		certs, err = activeStorage.ListCertificates(testusers.Bob.AuthID(), wdk.ListCertificatesArgs{})
		require.NoError(t, err)
		require.Equal(t, expectedResult, certs)
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
		_, err = activeStorage.InsertCertificateAuth(testusers.Alice.AuthID(), &wdk.TableCertificateX{
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
		err := activeStorage.RelinquishCertificate(
			testusers.Alice.AuthID(),
			wdk.RelinquishCertificateArgs{
				Type:         "not-type",
				SerialNumber: "not-serial",
				Certifier:    "not-certifier",
			})

		// then:
		require.ErrorContains(t, err, "failed to delete certificate model: certificate not found")
	})
}
