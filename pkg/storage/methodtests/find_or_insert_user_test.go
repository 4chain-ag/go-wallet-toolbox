package methodtests_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOrInsertUser(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	userIdentityKey := "03f17660f611ce531402a2ce1e070380b6fde57aca211d707bfab27bce42d86beb"

	// and:
	activeStorage := given.GormProviderWithCleanDatabase()

	// when:
	tableUser, err := activeStorage.FindOrInsertUser(userIdentityKey)

	// then:
	require.NoError(t, err)

	assert.Equal(t, true, tableUser.IsNew)
	assert.Equal(t, userIdentityKey, tableUser.User.IdentityKey)

	// and when:
	tableUser, err = activeStorage.FindOrInsertUser(userIdentityKey)

	// then:
	require.NoError(t, err)

	assert.Equal(t, false, tableUser.IsNew)
	assert.Equal(t, userIdentityKey, tableUser.User.IdentityKey)
}
