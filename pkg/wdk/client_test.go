package wdk_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/integrationtests/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("test client call to actual server 2", func(t *testing.T) {
		// given:
		given := testabilities.Given(t)

		// and:
		activeStorage := given.GormProvider()

		// and:
		cleanupSrv := given.StartedRPCServerFor(activeStorage)
		defer cleanupSrv()

		// and:
		mux := http.NewServeMux()
		httpClient := httptest.NewServer(mux).Client()
		userIdentityKey := "03f17660f611ce531402a2ce1e070380b6fde57aca211d707bfab27bce42d86beb"

		client, cleanup, err := wdk.NewClient(
			given.ServerURL(),
			wdk.WithHttpClient(httpClient),
		)
		defer cleanup()

		require.NoError(t, err)
		require.NotNil(t, client)

		tableUser, err := client.FindOrInsertUser(userIdentityKey)
		require.NoError(t, err)
		require.Equal(t, userIdentityKey, tableUser.User.IdentityKey)
	})
}
