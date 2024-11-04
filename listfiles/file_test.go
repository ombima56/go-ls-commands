package listfiles_test

import (
	"testing"

	"github.com/ombima56/go-ls-commands/listfiles"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		args        []string
		expectedErr bool
		expected    [5]bool // Represents longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag
	}{
		{[]string{"-l"}, false, [5]bool{true, false, false, false, false}},
		{[]string{"-a"}, false, [5]bool{false, true, false, false, false}},
		{[]string{"-R"}, false, [5]bool{false, false, true, false, false}},
		{[]string{"-t"}, false, [5]bool{false, false, false, true, false}},
		{[]string{"-r"}, false, [5]bool{false, false, false, false, true}},
		{[]string{"--long"}, false, [5]bool{true, false, false, false, false}},
		{[]string{"--invalid"}, true, [5]bool{false, false, false, false, false}},
	}
}