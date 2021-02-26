package redmine

import (
	"errors"
	"github.com/stretchr/testify/require"
	"reflect"
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
	}
	tests := []struct {
		name     string
		authType AuthType
		want     error
	}{
		{"AuthTypeBasicAuth should validate", AuthTypeBasicAuth, nil},
		{"AuthTypeBasicAuth should validate", AuthTypeBasicAuthWithTokenPassword, nil},
		{"AuthTypeBasicAuth should validate", AuthTypeTokenQueryParam, nil},
		{"AuthTypeBasicAuth should validate", AuthTypeNoAuth, nil},
		{"negative int should return error", -1, errors.New("invalid AuthType -1 found")},
		{"invalid int should return error", 4, errors.New("invalid AuthType 4 found")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := APIAuth{
				AuthType: tt.authType,
			}
			if got := auth.validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
