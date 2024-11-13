package sorting

import (
	"reflect"
	"testing"
)

func TestSortFiles(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []string{"file.txt"},
			expected: []string{"file.txt"},
		},
		{
			name:     "All lowercase",
			input:    []string{"zebra.txt", "apple.txt", "banana.txt"},
			expected: []string{"apple.txt", "banana.txt", "zebra.txt"},
		},
		{
			name:     "All uppercase",
			input:    []string{"ZEBRA.txt", "APPLE.txt", "BANANA.txt"},
			expected: []string{"APPLE.txt", "BANANA.txt", "ZEBRA.txt"},
		},
		{
			name:     "Mixed case",
			input:    []string{"Zebra.txt", "apple.txt", "BANANA.txt"},
			expected: []string{"apple.txt", "BANANA.txt", "Zebra.txt"},
		},
		{
			name:     "With numbers",
			input:    []string{"file2.txt", "file10.txt", "file1.txt"},
			expected: []string{"file1.txt", "file10.txt", "file2.txt"},
		},
		{
			name:     "With special characters",
			input:    []string{"_file.txt", ".hidden", "file.txt"},
			expected: []string{"file.txt", ".hidden", "_file.txt"},
		},
		{
			name:     "Different extensions",
			input:    []string{"test.go", "test.txt", "test.md"},
			expected: []string{"test.go", "test.md", "test.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			SortFiles(input)

			if !reflect.DeepEqual(input, tt.expected) {
				t.Errorf("SortFiles() = %v, want %v", input, tt.expected)
			}
		})
	}
}
