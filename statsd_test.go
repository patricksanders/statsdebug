package main

import (
	"reflect"
	"testing"
)

func TestParseStat(t *testing.T) {
	var cases = []struct {
		description string
		input       string
		shouldError bool
		expected    *Stat
	}{
		{
			description: "Boring stat",
			input:       "bar.foo.baz:5|c",
			shouldError: false,
			expected: &Stat{
				Name:  "bar.foo.baz",
				Type:  "c",
				Value: "5",
			},
		},
		{
			description: "Stat with sample rate",
			input:       "bar.foo.baz:5|c|@0.5",
			shouldError: false,
			expected: &Stat{
				Name:       "bar.foo.baz",
				Type:       "c",
				Value:      "5",
				SampleRate: "@0.5",
			},
		},
		{
			description: "Stat with tags",
			input:       "bar.foo.baz:5|c#foo:bar,baz:bang",
			shouldError: false,
			expected: &Stat{
				Name:     "bar.foo.baz",
				Type:     "c",
				Value:    "5",
				Tags:     []string{"baz:bang", "foo:bar"},
			},
		},
		{
			description: "Stat with sample rate and tags",
			input:       "bar.foo.baz:5|c|@0.5#foo:bar,baz:bang",
			shouldError: false,
			expected: &Stat{
				Name:       "bar.foo.baz",
				Type:       "c",
				Value:      "5",
				SampleRate: "@0.5",
				Tags:       []string{"baz:bang", "foo:bar"},
			},
		},
		{
			description: "Not a stat at all",
			input:       "Hello there!",
			shouldError: true,
			expected:    nil,
		},
	}

	for _, tc := range cases {
		t.Logf("Running test case '%s' ", tc.description)
		res, err := parseStat(tc.input)
		if err != nil && !tc.shouldError {
			t.Errorf("Test failed with error %s", err)
		} else if !reflect.DeepEqual(tc.expected, res) {
			t.Errorf("Expected %v but got %v", tc.expected, res)
		}
	}
}
