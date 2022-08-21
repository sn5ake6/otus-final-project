package internalhttp

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sn5ake6/otus-final-project/internal/app"
	"github.com/sn5ake6/otus-final-project/internal/bucket"
	"github.com/sn5ake6/otus-final-project/internal/config"
	"github.com/sn5ake6/otus-final-project/internal/logger"
	memorystorage "github.com/sn5ake6/otus-final-project/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerOperations(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	limit := config.LimitConf{
		Login:         1,
		Password:      10,
		IP:            100,
		ResetInterval: "1s",
	}

	ctx := context.Background()
	bucket := bucket.NewLeakyBucket(ctx, limit)

	storage := memorystorage.New()
	router := NewRouter(
		logg,
		app.New(logg, storage, bucket),
	)

	subnet := `{
		"subnet": "192.1.1.0/25"
	}`

	authorize := `{
		"login": "login",
		"password": "password",
		"ip": "192.1.1.1"
	}`

	t.Run("authorize success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		checkSuccessAuthorizeResponse(t, w)
	})

	t.Run("add to blacklist success case", func(t *testing.T) {
		checkSuccessAddSubnet(t, router, blacklistURI, subnet)
	})

	t.Run("authorize from blacklist fail case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		checkFailAuthorizeResponse(t, w)
	})

	t.Run("add to blacklist fail case", func(t *testing.T) {
		checkFailAddSubnet(t, router, blacklistURI, subnet)
	})

	t.Run("delete from blacklist success cases", func(t *testing.T) {
		checkSuccessDeleteSubnet(t, router, blacklistURI, subnet)
	})

	t.Run("delete from blacklist fail cases", func(t *testing.T) {
		checkFailDeleteSubnet(t, router, blacklistURI, subnet)
	})

	t.Run("add to whitelist success case", func(t *testing.T) {
		checkSuccessAddSubnet(t, router, whitelistURI, subnet)
	})

	t.Run("authorize from whitelist success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		checkSuccessAuthorizeResponse(t, w)
	})

	t.Run("add to whitelist fail case", func(t *testing.T) {
		checkFailAddSubnet(t, router, whitelistURI, subnet)
	})

	t.Run("delete from whitelist success cases", func(t *testing.T) {
		checkSuccessDeleteSubnet(t, router, whitelistURI, subnet)
	})

	t.Run("delete from whitelist fail cases", func(t *testing.T) {
		checkFailDeleteSubnet(t, router, whitelistURI, subnet)
	})
}

func TestHTTPServerAuthorizeAndResetOperations(t *testing.T) {
	logg, err := logger.New("info")
	require.NoError(t, err)

	limit := config.LimitConf{
		Login:         1,
		Password:      10,
		IP:            100,
		ResetInterval: "1s",
	}

	ctx := context.Background()
	bucket := bucket.NewLeakyBucket(ctx, limit)

	storage := memorystorage.New()
	router := NewRouter(
		logg,
		app.New(logg, storage, bucket),
	)

	authorize := `{
		"login": "login",
		"password": "password",
		"ip": "192.1.1.1"
	}`

	t.Run("authorize success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		checkSuccessAuthorizeResponse(t, w)
	})

	t.Run("authorize fail case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		checkFailAuthorizeResponse(t, w)
	})

	t.Run("authorize after reset interval success case", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		checkSuccessAuthorizeResponse(t, w)
	})

	t.Run("authorize after reset success case", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		reqFailAuthorize := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		wFailAuthorize := httptest.NewRecorder()

		router.ServeHTTP(wFailAuthorize, reqFailAuthorize)

		checkFailAuthorizeResponse(t, wFailAuthorize)

		reqReset := httptest.NewRequest(http.MethodPost, resetURI, bytes.NewBufferString(authorize))
		wReset := httptest.NewRecorder()

		router.ServeHTTP(wReset, reqReset)

		respReset := wReset.Result()
		defer respReset.Body.Close()
		require.Equal(t, http.StatusOK, respReset.StatusCode)

		reqSuccessAuthorize := httptest.NewRequest(http.MethodPost, authorizeURI, bytes.NewBufferString(authorize))
		wSuccessAuthorize := httptest.NewRecorder()

		router.ServeHTTP(wSuccessAuthorize, reqSuccessAuthorize)

		checkSuccessAuthorizeResponse(t, wSuccessAuthorize)
	})
}

func checkSuccessAddSubnet(t *testing.T, router http.Handler, uri, subnet string) {
	t.Helper()
	statusCode := getSubnetResponseStatusCode(uri, http.MethodPost, subnet, router)
	require.Equal(t, http.StatusCreated, statusCode)
}

func checkFailAddSubnet(t *testing.T, router http.Handler, uri, subnet string) {
	t.Helper()
	statusCode := getSubnetResponseStatusCode(uri, http.MethodPost, subnet, router)
	require.NotEqual(t, http.StatusCreated, statusCode)
}

func checkSuccessDeleteSubnet(t *testing.T, router http.Handler, uri, subnet string) {
	t.Helper()
	statusCode := getSubnetResponseStatusCode(uri, http.MethodDelete, subnet, router)
	require.Equal(t, http.StatusOK, statusCode)
}

func checkFailDeleteSubnet(t *testing.T, router http.Handler, uri, subnet string) {
	t.Helper()
	statusCode := getSubnetResponseStatusCode(uri, http.MethodDelete, subnet, router)
	require.NotEqual(t, http.StatusOK, statusCode)
}

func getSubnetResponseStatusCode(uri, method, subnet string, router http.Handler) int {
	req := httptest.NewRequest(method, uri, bytes.NewBufferString(subnet))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	return resp.StatusCode
}

func checkFailAuthorizeResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	checkResponse(t, w, http.StatusTooManyRequests, `{"ok": false}`)
}

func checkSuccessAuthorizeResponse(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	checkResponse(t, w, http.StatusOK, `{"ok": true}`)
}

func checkResponse(
	t *testing.T,
	w *httptest.ResponseRecorder,
	expectedStatus int,
	expectedBody string,
) {
	t.Helper()
	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, expectedStatus, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expectedBody, string(body))
}
