package redmine

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const (
	urlSimple                  = "http://localhost/endpoint"
	urlWithPort                = "http://localhost:3000/endpoint"
	urlWithPortContextPath     = "http://localhost:3000/endpoint"
	urlWithPortPathQueryParams = "http://localhost:3000/endpoint?key=value&key=doublevalue&important_id=2&specialCharacter=äöüß+àÀ%20."
)

func TestClient_concatParameters(t *testing.T) {
	type args struct {
		requestParameters []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should return empty string for zero parameters", args{}, ""},
		{"should return same string for one parameters", args{[]string{"key=value"}}, "key=value"},
		{"should return &-delimited string for two parameters", args{[]string{"key=value", "hello=world"}},
			"key=value&hello=world"},
		{"should remove empty parameter at the start and return two parameters", args{[]string{"", "key=value", "hello=world"}},
			"key=value&hello=world"},
		{"should remove empty parameter in the middle and return two parameters", args{[]string{"key=value", "", "hello=world"}},
			"key=value&hello=world"},
		{"should remove empty parameter in the end and return two parameters", args{[]string{"key=value", "hello=world", ""}},
			"key=value&hello=world"},
		{"should remove multiple empty parameter at the start and return two parameters", args{[]string{"", "", "key=value", "hello=world"}},
			"key=value&hello=world"},
		{"should remove multiple empty parameter in the middle and return two parameters", args{[]string{"key=value", "", "", "hello=world"}},
			"key=value&hello=world"},
		{"should remove multiple empty parameter in the end and return two parameters", args{[]string{"key=value", "hello=world", "", ""}},
			"key=value&hello=world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if got := c.concatParameters(tt.args.requestParameters...); got != tt.want {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAPIAuth_validate(t *testing.T) {
	type fields struct {
		AuthType AuthType
		Token    string
		User     string
		Password string
	}
	tests := []struct {
		name string
		args fields
		want error
	}{
		{"should validate basic auth", fields{AuthType: AuthTypeBasicAuth, User: "leUser"}, nil},
		{"should fail basic auth with empty user", fields{AuthType: AuthTypeBasicAuth}, errors.New("invalid auth configuration for type 0: user must not be empty")},

		{"should validate basic auth/token", fields{AuthType: AuthTypeBasicAuthWithTokenPassword, User: "leUser", Token: "leToken"}, nil},
		{"should fail basic auth/token with empty user", fields{AuthType: AuthTypeBasicAuthWithTokenPassword, Token: "leToken"}, errors.New("invalid auth configuration for type 2: user must not be empty")},
		{"should fail basic auth/token with empty token", fields{AuthType: AuthTypeBasicAuthWithTokenPassword, User: "leUser"}, errors.New("invalid auth configuration for type 2: API token must not be empty")},

		{"should validate token query auth", fields{AuthType: AuthTypeTokenQueryParam, Token: "leToken"}, nil},
		{"should fail token query auth with empty token", fields{AuthType: AuthTypeTokenQueryParam}, errors.New("invalid auth configuration for type 1: API token must not be empty")},

		{"should validate (no auth)", fields{AuthType: AuthTypeNoAuth}, nil},

		{"negative int should return error", fields{AuthType: -1}, errors.New("invalid auth configuration: AuthType -1 found")},
		{"invalid int should return error", fields{AuthType: 4}, errors.New("invalid auth configuration: AuthType 4 found")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := APIAuth{
				AuthType: tt.args.AuthType,
				Token:    tt.args.Token,
				User:     tt.args.User,
				Password: tt.args.Password,
			}
			if got := auth.validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Run("should return Client without error", func(t *testing.T) {
		sut, err := NewClient("http://localhost:3000/", APIAuth{AuthType: AuthTypeNoAuth})

		require.NoError(t, err)
		require.NotNil(t, sut)
	})
	t.Run("should return error on APIAuth misconfiguration", func(t *testing.T) {
		sut, err := NewClient("http://localhost:3000/", APIAuth{AuthType: 5})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not create redmine client:")
		assert.Contains(t, err.Error(), "AuthType 5 found")
		require.Nil(t, sut)
	})
}

func Test_safelyAddQueryParameter(t *testing.T) {
	const endpoint = "http://1.2.3.4:3030/endpoint"

	type args struct {
		req   *http.Request
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"should not alter URL on empty parameter key", args{reqFromURL(t, ""), "", "asdfasd"}, endpoint, false},
		{"should add delimiter ? after endpoint without query parameters", args{reqFromURL(t, ""), "key", "value"}, endpoint + "?key=value", false},
		{"should add another k/v pair delimited by &", args{reqFromURL(t, "?anotherKey=1"), "key", "value"}, endpoint + "?anotherKey=1&key=value", false},
		{"should encode spaces", args{reqFromURL(t, "?anotherKey=1"), "key", "space & ampersands"}, endpoint + "?anotherKey=1&key=space+%26+ampersands", false},
		{"should sort query existing and added params", args{reqFromURL(t, "?a=afterZ&z=afteraAndAlsoTheEnd&Z=afterA"), "A", "start"}, endpoint + "?A=start&Z=afterA&a=afterZ&z=afteraAndAlsoTheEnd", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := safelyAddQueryParameter(tt.args.req, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("safelyAddQueryParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := tt.args.req.URL.String()
			if got != tt.want {
				t.Errorf("safelyAddQueryParameter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSafelyAddQueryParameter_requestSafety(t *testing.T) {
	t.Run("should not change request besides URL", func(t *testing.T) {
		var ir projectRequest
		ir.Project = Project{
			Name:       "name",
			Identifier: "ident",
		}
		s, err := json.Marshal(ir)
		assert.NoError(t, err)
		postReader := strings.NewReader(string(s))
		expectedReaderLen := postReader.Len()

		sutGET, _ := http.NewRequest("GET", "http://1.2.3.4:3030/endpoint?getAddInfo=true", nil)
		sutPOST, _ := http.NewRequest("POST", "http://1.2.3.4:3030/endpoint?doTheThingDifferently", postReader)

		// when
		getErr := safelyAddQueryParameter(sutGET, "key", "value")
		postErr := safelyAddQueryParameter(sutPOST, "key", "value")

		// then
		require.NoError(t, getErr)
		assert.Equal(t, "GET", sutGET.Method)
		assert.Nil(t, sutGET.Body)
		assert.Contains(t, sutGET.URL.String(), "http://1.2.3.4:3030/endpoint?")

		require.NoError(t, postErr)
		assert.Equal(t, "POST", sutPOST.Method)
		assert.NotNil(t, postReader)
		assert.Equal(t, expectedReaderLen, postReader.Len())
		assert.Contains(t, sutPOST.URL.String(), "http://1.2.3.4:3030/endpoint?")
	})
}

func reqFromURL(t *testing.T, params string) *http.Request {
	t.Helper()

	req, err := http.NewRequest("GET", "http://1.2.3.4:3030/endpoint"+params, nil)
	assert.NoError(t, err)

	return req
}
