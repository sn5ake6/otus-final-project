package bucket

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/sn5ake6/otus-final-project/internal/config"
	"github.com/sn5ake6/otus-final-project/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestBucket(t *testing.T) {
	ctx := context.Background()

	limit := config.LimitConf{
		Login:         1,
		Password:      10,
		IP:            100,
		ResetInterval: "1s",
	}

	bucket := NewLeakyBucket(ctx, limit)

	go bucket.leak()

	authorize := storage.Authorize{
		Login:    "login",
		Password: "password",
		IP:       "127.0.0.1",
	}

	authorize2 := storage.Authorize{
		Login:    "login2",
		Password: "password2",
		IP:       "127.0.0.2",
	}

	t.Run("check authorize login case", func(t *testing.T) {
		bucket.resetAll()
		require.True(t, bucket.Check(authorize))

		require.False(t, bucket.Check(authorize))
	})

	t.Run("check authorize password case", func(t *testing.T) {
		bucket.resetAll()
		for i := 0; i < int(limit.Password); i++ {
			authorize.Login = "login" + strconv.Itoa(i)
			require.True(t, bucket.Check(authorize))
		}

		overLimit := int(limit.Password + 1)

		authorize.Login = "login" + strconv.Itoa(overLimit)
		require.False(t, bucket.Check(authorize))
	})

	t.Run("check authorize ip case", func(t *testing.T) {
		bucket.resetAll()
		for i := 0; i < int(limit.IP); i++ {
			iAsString := strconv.Itoa(i)
			authorize.Login = "login" + iAsString
			authorize.Password = "password" + iAsString
			require.True(t, bucket.Check(authorize))
		}

		overLimit := strconv.Itoa(int(limit.IP + 1))

		authorize.Login = "login" + overLimit
		authorize.Password = "password" + overLimit
		require.False(t, bucket.Check(authorize))
	})

	t.Run("check reset bucket case", func(t *testing.T) {
		bucket.resetAll()

		require.True(t, bucket.Check(authorize))
		require.True(t, bucket.Check(authorize2))

		require.False(t, bucket.Check(authorize))
		require.False(t, bucket.Check(authorize2))

		bucket.Reset(authorize)

		require.True(t, bucket.Check(authorize))
		require.False(t, bucket.Check(authorize2))
	})

	t.Run("check reset interval case", func(t *testing.T) {
		bucket.resetAll()
		require.True(t, bucket.Check(authorize))
		require.True(t, bucket.Check(authorize2))

		require.False(t, bucket.Check(authorize))
		require.False(t, bucket.Check(authorize2))

		time.Sleep(time.Second * 2)

		require.True(t, bucket.Check(authorize))
		require.True(t, bucket.Check(authorize2))
	})
}
