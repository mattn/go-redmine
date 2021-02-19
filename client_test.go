package redmine

import (
	"github.com/stretchr/testify/require"
	"testing"
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
