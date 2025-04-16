package satoshi_test

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("add two ints", func(t *testing.T) {
		c, err := satoshi.Add(1, 2)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(3), c)
	})

	t.Run("add int and uint64", func(t *testing.T) {
		c, err := satoshi.Add(1, uint64(2))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(3), c)
	})

	t.Run("add uint and negative int", func(t *testing.T) {
		c, err := satoshi.Add(uint(1), -2)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-1), c)
	})

	t.Run("add two uints", func(t *testing.T) {
		c, err := satoshi.Add(uint(1), uint(2))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(3), c)
	})

	t.Run("add two max satoshis", func(t *testing.T) {
		_, err := satoshi.Add(primitives.MaxSatoshis, primitives.MaxSatoshis)
		require.Error(t, err)
	})
}

func TestSubtract(t *testing.T) {
	t.Run("subtract two ints", func(t *testing.T) {
		c, err := satoshi.Subtract(5, 3)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(2), c)
	})

	t.Run("subtract int and uint64", func(t *testing.T) {
		c, err := satoshi.Subtract(5, uint64(2))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(3), c)
	})

	t.Run("subtract uint and negative int", func(t *testing.T) {
		// 1 - (-2) equals 3
		c, err := satoshi.Subtract(uint(1), -2)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(3), c)
	})

	t.Run("subtract resulting in zero", func(t *testing.T) {
		c, err := satoshi.Subtract(2, 2)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(0), c)
	})

	t.Run("subtract to obtain a negative result", func(t *testing.T) {
		c, err := satoshi.Subtract(3, 5)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-2), c)
	})

	t.Run("subtract exceeding max positive value", func(t *testing.T) {
		// primitives.MaxSatoshis - (-1) equals primitives.MaxSatoshis + 1 (overflow)
		_, err := satoshi.Subtract(primitives.MaxSatoshis, -1)
		require.Error(t, err)
	})

	t.Run("subtract exceeding max negative value", func(t *testing.T) {
		// (-primitives.MaxSatoshis) - 1 equals -(primitives.MaxSatoshis + 1) (underflow)
		_, err := satoshi.Subtract(-primitives.MaxSatoshis, 1)
		require.Error(t, err)
	})
}

type otherTypeAlias int64

func TestFrom(t *testing.T) {
	t.Run("from int", func(t *testing.T) {
		c, err := satoshi.From(int64(1))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(1), c)
	})

	t.Run("from uint", func(t *testing.T) {
		c, err := satoshi.From(uint64(1))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(1), c)
	})

	t.Run("from negative int", func(t *testing.T) {
		c, err := satoshi.From(int64(-1))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-1), c)
	})

	t.Run("from max uint32", func(t *testing.T) {
		c, err := satoshi.From(uint32(math.MaxUint32))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(4294967295), c)
	})

	t.Run("from max int32", func(t *testing.T) {
		c, err := satoshi.From(int32(math.MaxInt32))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(2147483647), c)
	})

	t.Run("from min int32", func(t *testing.T) {
		c, err := satoshi.From(int32(math.MinInt32))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-2147483648), c)
	})

	t.Run("from max int64", func(t *testing.T) {
		_, err := satoshi.From(int64(math.MaxInt64))
		require.Error(t, err)
	})

	t.Run("from min int64", func(t *testing.T) {
		_, err := satoshi.From(int64(math.MinInt64))
		require.Error(t, err)
	})

	t.Run("from max uint64", func(t *testing.T) {
		_, err := satoshi.From(uint64(math.MaxUint64))
		require.Error(t, err)
	})

	t.Run("from max satoshi", func(t *testing.T) {
		c, err := satoshi.From(primitives.MaxSatoshis)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(primitives.MaxSatoshis), c)
	})

	t.Run("from negative max satoshi", func(t *testing.T) {
		c, err := satoshi.From(-primitives.MaxSatoshis)
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-primitives.MaxSatoshis), c)
	})

	t.Run("from max satoshi + 1", func(t *testing.T) {
		_, err := satoshi.From(primitives.MaxSatoshis + 1)
		require.Error(t, err)
	})

	t.Run("from negative max satoshi - 1", func(t *testing.T) {
		_, err := satoshi.From(-primitives.MaxSatoshis - 1)
		require.Error(t, err)
	})

	t.Run("from max satoshi as uint64", func(t *testing.T) {
		c, err := satoshi.From(uint64(primitives.MaxSatoshis))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(primitives.MaxSatoshis), c)
	})

	t.Run("from max satoshi as int64", func(t *testing.T) {
		c, err := satoshi.From(int64(primitives.MaxSatoshis))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(primitives.MaxSatoshis), c)
	})

	t.Run("from negative max satoshi as int64", func(t *testing.T) {
		c, err := satoshi.From(int64(-primitives.MaxSatoshis))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(-primitives.MaxSatoshis), c)
	})

	t.Run("from other type alias", func(t *testing.T) {
		c, err := satoshi.From(otherTypeAlias(1))
		require.NoError(t, err)
		require.Equal(t, satoshi.Value(1), c)
	})

	t.Run("from other type alias equals max satoshi + 1", func(t *testing.T) {
		_, err := satoshi.From(otherTypeAlias(primitives.MaxSatoshis + 1))
		require.Error(t, err)
	})
}
