package methodtests_test

import (
	"fmt"
	"testing"

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
	userIdentityKey := "03f17660f611ce531402a2ce1e070380b6fde57aca211d707bfab27bce42d86beb"
	user, err := activeStorage.FindOrInsertUser(userIdentityKey)
	require.NoError(t, err)

	// and:
	certToInsert := &wdk.TableCertificateX{
		TableCertificate: wdk.TableCertificate{
			UserID:             user.User.UserID,
			Type:               "exampleType",
			SerialNumber:       "exampleSerialNumber",
			Certifier:          "certifierPubKey",
			Subject:            "subjectPubKey",
			RevocationOutpoint: "outpointString",
			Signature:          "signatureHex",
		},
		Fields: []*wdk.TableCertificateField{
			{
				UserID:     user.User.UserID,
				FieldName:  "exampleField",
				FieldValue: "exampleValue",
				MasterKey:  "exampleMasterKey",
			},
		},
	}
	// when:
	id, err := activeStorage.InsertCertificateAuth(wdk.AuthID{
		UserID: &user.User.UserID,
	}, certToInsert)

	// then:
	require.NoError(t, err)
	require.Positive(t, id)

	// List certs
	// check if id exists there
	certs, err := activeStorage.ListCertificates(wdk.AuthID{
		UserID: &user.User.UserID,
	}, wdk.ListCertificatesArgs{})
	require.NoError(t, err)
	fmt.Println(certs)

	// assert.Equal(t, testabilities.StorageName, tableSettings.StorageName)
	// assert.Equal(t, testabilities.StorageIdentityKey, tableSettings.StorageIdentityKey)
	// assert.Equal(t, defs.NetworkTestnet, tableSettings.Chain)
	// assert.Equal(t, 1024, tableSettings.MaxOutputScript)
}
