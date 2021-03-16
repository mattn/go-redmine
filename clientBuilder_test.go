package redmine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	testEndpoint = "http://example"
	testAPIToken = "1q2w3e4r5t6z"
)

func TestNewClientBuilder(t *testing.T) {
	sut := NewClientBuilder()
	assert.Equal(t, -1, sut.resultLimit)
	assert.Equal(t, -1, sut.resultOffset)
}

func TestClientBuilder_Build(t *testing.T) {
	t.Run("should successfully creates a client with defaults", func(t *testing.T) {
		sut, err := NewClientBuilder().
			Endpoint(testEndpoint).
			AuthAPIToken(testAPIToken).
			Build()

		require.NoError(t, err)
		require.NotNil(t, sut)
		assert.Equal(t, testEndpoint, sut.endpoint)
		assert.Equal(t, APIAuth{AuthType: AuthTypeTokenQueryParam, Token: testAPIToken}, sut.auth)
		assert.Equal(t, -1, sut.Limit)
		assert.Equal(t, -1, sut.Offset)
		assert.Equal(t, sut.Client, http.DefaultClient)
	})

	t.Run("should successfully creates a client without default values", func(t *testing.T) {
		sut, err := NewClientBuilder().
			Endpoint(testEndpoint).
			AuthAPIToken(testAPIToken).
			ResultLimit(12).
			ResultOffset(5).
			SkipSSLVerify(true).
			Build()

		require.NoError(t, err)
		require.NotNil(t, sut)
		assert.Equal(t, testEndpoint, sut.endpoint)
		assert.Equal(t, APIAuth{AuthType: AuthTypeTokenQueryParam, Token: testAPIToken}, sut.auth)
		assert.Equal(t, 12, sut.Limit)
		assert.Equal(t, 5, sut.Offset)
		assert.NotEqual(t, sut.Client, http.DefaultClient)
	})

	t.Run("should overwrite authentication configuration", func(t *testing.T) {
		sut, err := NewClientBuilder().
			Endpoint(testEndpoint).
			AuthAPIToken(testAPIToken).
			AuthBasicAuth("user", "1234").
			Build()

		require.NoError(t, err)
		require.NotNil(t, sut)
		assert.Equal(t, testEndpoint, sut.endpoint)
		assert.Equal(t, APIAuth{AuthType: AuthTypeBasicAuth, User: "user", Password: "1234"}, sut.auth)
		assert.Equal(t, -1, sut.Limit)
		assert.Equal(t, -1, sut.Offset)
	})

	t.Run("should fail create a client on no endpoint", func(t *testing.T) {
		sut, err := NewClientBuilder().Build()

		require.Error(t, err)
		require.Nil(t, sut)
		assert.Contains(t, err.Error(), "endpoint must not be empty")
	})

	t.Run("should fail create a client on missing authentication but pass on basic auth", func(t *testing.T) {
		sut := NewClientBuilder()

		// when
		client, err := sut.
			Endpoint(testEndpoint).
			Build()

		// then
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "authentication must not be empty")

		// when again
		client, err = sut.
			AuthBasicAuth("user", "pass").
			Build()

		// then again
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("should fail create a client on missing authentication but pass on API key", func(t *testing.T) {
		sut := NewClientBuilder()

		// when
		client, err := sut.
			Endpoint(testEndpoint).
			Build()

		// then
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "authentication must not be empty")

		// when again
		client, err = sut.
			AuthAPIToken(testAPIToken).
			Build()

		// then again
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}
