package config

import (
	"os"
	"testing"
)

func TestGetEscapedEnv(t *testing.T) {
	var tests = []struct {
		env    string
		result string
	}{
		{"", ""},
		{"plain", "plain"},
		{`Line1
Line2`, `"Line1\nLine2"`},
		{`Line1"
Line2`, `"Line1\"\nLine2"`},
		{`Line1\nLine2`, `"Line1\nLine2"`},
		{`Line1"\nLine2`, `"Line1\"\nLine2"`},
		{"Line1\r\nLine2", `"Line1\nLine2"`},
	}

	for _, test := range tests {
		os.Setenv("TestGetEscapedEnv", test.env)
		result := getEscapedEnv("TestGetEscapedEnv")
		if result != test.result {
			t.Errorf("got %s, wanted %s", result, test.result)
		}
	}
}
