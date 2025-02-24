package listfiles_test

import (
	"go-ls-commands/listfiles"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		args        []string
		expectedErr bool
		expected    listfiles.Options
	}{
		// Single short flags
		{[]string{"-l"}, false, listfiles.Options{LongFormat: true}},
		{[]string{"-a"}, false, listfiles.Options{AllFiles: true}},
		{[]string{"-R"}, false, listfiles.Options{Recursive: true}},
		{[]string{"-t"}, false, listfiles.Options{SortByTime: true}},
		{[]string{"-r"}, false, listfiles.Options{ReverseSort: true}},

		// Single long flags
		{[]string{"--long"}, false, listfiles.Options{LongFormat: true}},
		{[]string{"--all"}, false, listfiles.Options{AllFiles: true}},
		{[]string{"--recursive"}, false, listfiles.Options{Recursive: true}},
		{[]string{"--time"}, false, listfiles.Options{SortByTime: true}},
		{[]string{"--reverse"}, false, listfiles.Options{ReverseSort: true}},

		// Combined short flags
		{[]string{"-la"}, false, listfiles.Options{LongFormat: true, AllFiles: true}},
		{[]string{"-rt"}, false, listfiles.Options{ReverseSort: true, SortByTime: true}},
		{[]string{"-lR"}, false, listfiles.Options{LongFormat: true, Recursive: true}},

		// Multiple separate flags
		{[]string{"-l", "-a"}, false, listfiles.Options{LongFormat: true, AllFiles: true}},
		{[]string{"-r", "--long"}, false, listfiles.Options{LongFormat: true, ReverseSort: true}},

		// Duplicate flags should still work
		{[]string{"-ll"}, false, listfiles.Options{LongFormat: true}},
		{[]string{"-aa"}, false, listfiles.Options{AllFiles: true}},

		// Invalid flags
		{[]string{"--invalid"}, true, listfiles.Options{}},
		{[]string{"-x"}, true, listfiles.Options{}},
		{[]string{"-l", "-x"}, true, listfiles.Options{}},
		{[]string{}, false, listfiles.Options{}},
	}

	for _, test := range tests {
		opts, err := listfiles.ValidateFlags(test.args)

		// Check if an error occurred when it shouldn't or vice versa
		if (err != nil) != test.expectedErr {
			t.Errorf("ValidateFlags(%v) error = %v, expectedErr = %v", test.args, err, test.expectedErr)
		}

		// Compare the returned options struct with the expected struct
		if opts != test.expected {
			t.Errorf("ValidateFlags(%v) = %+v, want %+v", test.args, opts, test.expected)
		}
	}
}
