package memorystorage

import (
	"errors"
	"testing"

	memorystorage "github.com/sn5ake6/otus-final-project/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()

	subnet := "192.1.1.0/25"
	ipFromSubnet := "192.1.1.6"
	ipNotFromSubnet := "127.0.0.1"

	//nolint:dupl
	t.Run("blacklist cases", func(t *testing.T) {
		err := storage.AddToBlacklist(subnet)
		require.NoError(t, err)

		err = storage.AddToBlacklist(subnet)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, memorystorage.ErrSubnetAlreadyExists))

		result, err := storage.FindIPInBlacklist(ipFromSubnet)
		require.NoError(t, err)
		require.True(t, result)

		result, err = storage.FindIPInBlacklist(ipNotFromSubnet)
		require.Nil(t, err)
		require.False(t, result)

		err = storage.DeleteFromBlacklist(subnet)
		require.NoError(t, err)

		err = storage.DeleteFromBlacklist(subnet)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, memorystorage.ErrSubnetNotExists))

		result, err = storage.FindIPInBlacklist(ipFromSubnet)
		require.Nil(t, err)
		require.False(t, result)
	})

	//nolint:dupl
	t.Run("whitelist cases", func(t *testing.T) {
		err := storage.AddToWhitelist(subnet)
		require.NoError(t, err)

		err = storage.AddToWhitelist(subnet)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, memorystorage.ErrSubnetAlreadyExists))

		result, err := storage.FindIPInWhitelist(ipFromSubnet)
		require.NoError(t, err)
		require.True(t, result)

		result, err = storage.FindIPInWhitelist(ipNotFromSubnet)
		require.Nil(t, err)
		require.False(t, result)

		err = storage.DeleteFromWhitelist(subnet)
		require.NoError(t, err)

		err = storage.DeleteFromWhitelist(subnet)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, memorystorage.ErrSubnetNotExists))

		result, err = storage.FindIPInWhitelist(ipFromSubnet)
		require.Nil(t, err)
		require.False(t, result)
	})
}
