package wildcard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Match(t *testing.T) {
	tests := []struct {
		Name        string
		Template    string
		Str         string
		ExpectMatch bool
	}{
		{
			Name:        "Should match when strings are equal - no wildcard",
			Template:    "hello world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string end and uses it to match",
			Template:    "hello *",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string end and does not use it to match",
			Template:    "hello world*",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string start and uses it to match",
			Template:    "* world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string start and does not use it to match",
			Template:    "*hello world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string middle and uses it to match",
			Template:    "hello*world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string middle and does not use it to match (A)",
			Template:    "hello *world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should match when string fits template with wildcard at string middle and does not use it to match (B)",
			Template:    "hello* world",
			Str:         "hello world",
			ExpectMatch: true,
		},
		{
			Name:        "Should NOT match when strings are not equal - no wildcard",
			Template:    "hello world",
			Str:         "goodbye world",
			ExpectMatch: false,
		},
		{
			Name:        "Should NOT match when string does not fits template with wildcard at string start",
			Template:    "* hello world", // note the extraspace
			Str:         "hello world",
			ExpectMatch: false,
		},
		{
			Name:        "Should NOT match when string does not fits template with wildcard at string end",
			Template:    "hello world *", // note the extra space
			Str:         "hello world",
			ExpectMatch: false,
		},
		{
			Name:        "Should NOT match when string does not fits template with wildcard at string middle",
			Template:    "hello * world", // note the extra space
			Str:         "hello world",
			ExpectMatch: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectMatch, Match(test.Template, test.Str))
		})
	}
}

func Test_wildcardToRegexp(t *testing.T) {
	tests := []struct {
		Name          string
		Template      string
		ExpectedRegex string
	}{
		{
			Name:          "Empty string template with no wildcards",
			Template:      "",
			ExpectedRegex: "^$",
		},
		{
			Name:          "Non empty string template with no wildcards",
			Template:      "hello",
			ExpectedRegex: "^hello$",
		},
		{
			Name:          "Template with one wildcard",
			Template:      "hello*",
			ExpectedRegex: "^hello.*$",
		},
		{
			Name:          "Template with multiple wildcards (A)",
			Template:      "hello**",
			ExpectedRegex: "^hello.*.*$",
		},
		{
			Name:          "Template with multiple wildcards (B)",
			Template:      "**hello",
			ExpectedRegex: "^.*.*hello$",
		},
		{
			Name:          "Template with multiple wildcards (C)",
			Template:      "*hel*lo*",
			ExpectedRegex: "^.*hel.*lo.*$",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedRegex, wildcardToRegexp(test.Template))
		})
	}
}
