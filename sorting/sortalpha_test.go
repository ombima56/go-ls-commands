package sorting

import (
	"io/fs"
	"testing"
	"time"
)

// MockFileinfo implements fs.FileInfo for testing
type MockFileinfo struct {
	name string
}

func (m MockFileinfo) Name() string       { return m.name }
func (m MockFileinfo) Size() int64        { return 0 }
func (m MockFileinfo) Mode() fs.FileMode  { return 0 }
func (m MockFileinfo) ModTime() time.Time { return time.Time{} }
func (m MockFileinfo) IsDir() bool        { return false }
func (m MockFileinfo) Sys() interface{}   { return nil }

func TestBubbleSortLowercaseFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Mixed case letters",
			input:    []string{"Zebra", "apple", "Banana", "cat"},
			expected: []string{"apple", "cat", "Banana", "Zebra"},
		},
		{
			name:     "All lowercase",
			input:    []string{"zebra", "apple", "banana", "cat"},
			expected: []string{"apple", "banana", "cat", "zebra"},
		},
		{
			name:     "All uppercase",
			input:    []string{"ZEBRA", "APPLE", "BANANA", "CAT"},
			expected: []string{"APPLE", "BANANA", "CAT", "ZEBRA"},
		},
		{
			name:     "With numbers and special characters",
			input:    []string{"1file", "!file", "File", "file"},
			expected: []string{"file", "File", "1file", "!file"},
		},
		// {
		// 	name:     "Empty strings",
		// 	input:    []string{"", "test", "", "Test"},
		// 	expected: []string{"test", "", "Test", ""},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create FileInfo slice
			files := make([]fs.FileInfo, len(tt.input))
			for i, name := range tt.input {
				files[i] = MockFileinfo{name: name}
			}

			// Sort the slice
			BubbleSortLowercaseFirst(files)

			// Check results
			for i, expected := range tt.expected {
				if files[i].Name() != expected {
					t.Errorf("Position %d: expected %s, got %s", i, expected, files[i].Name())
				}
			}
		})
	}
}

func TestShouldSwap(t *testing.T) {
	tests := []struct {
		name     string
		name1    string
		name2    string
		expected bool
	}{
		{
			name:     "Lowercase vs uppercase",
			name1:    "Zebra",
			name2:    "apple",
			expected: true,
		},
		{
			name:     "Same case different letters",
			name1:    "banana",
			name2:    "apple",
			expected: true,
		},
		{
			name:     "Special characters",
			name1:    "!file",
			name2:    "file",
			expected: true,
		},
		{
			name:     "Numbers",
			name1:    "1file",
			name2:    "afile",
			expected: true,
		},
		{
			name:     "Empty strings",
			name1:    "",
			name2:    "file",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSwap(tt.name1, tt.name2)
			if result != tt.expected {
				t.Errorf("shouldSwap(%s, %s) = %v; want %v",
					tt.name1, tt.name2, result, tt.expected)
			}
		})
	}
}

func TestTrimNotAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "With special characters",
			input:    "Hello!@#World123",
			expected: "HelloWorld123",
		},
		{
			name:     "Only special characters",
			input:    "!@#$%^",
			expected: "",
		},
		{
			name:     "Mixed alphanumeric",
			input:    "abc123DEF",
			expected: "abc123DEF",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimNotAlpha(tt.input)
			if result != tt.expected {
				t.Errorf("TrimNotAlpha(%s) = %s; want %s",
					tt.input, result, tt.expected)
			}
		})
	}
}
