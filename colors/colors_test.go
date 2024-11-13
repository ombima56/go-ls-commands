package colors

import (
	"os"
	"testing"
	"time"
)

// mockFileInfo implements os.FileInfo interface for testing
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

func TestGetFileColor(t *testing.T) {
	// Save original LS_COLORS
	originalLSColors := os.Getenv("LS_COLORS")
	defer os.Setenv("LS_COLORS", originalLSColors)

	// Set default LS_COLORS for this test
	os.Setenv("LS_COLORS", "di=01;34:ln=01;36:ex=01;32:bd=40;33;01:cd=40;33;01:pi=40;33")
	colorMap = Colors() // Update package-level colorMap

	// Set up test cases
	tests := []struct {
		name     string
		fileInfo mockFileInfo
		want     string
	}{
		{
			name: "Directory",
			fileInfo: mockFileInfo{
				name:  "testdir",
				mode:  os.ModeDir,
				isDir: true,
			},
			want: "\033[01;34m",
		},
		{
			name: "Symlink",
			fileInfo: mockFileInfo{
				name: "testlink",
				mode: os.ModeSymlink,
			},
			want: "\033[01;36m",
		},
		{
			name: "Executable",
			fileInfo: mockFileInfo{
				name: "testexec",
				mode: 0755,
			},
			want: "\033[01;32m",
		},
		// ... other test cases ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileColor(tt.fileInfo)
			if got != tt.want {
				t.Errorf("%s color mismatch:\ngot:  %s\nwant: %s", tt.name, got, tt.want)
			}
		})
	}
}

func TestColors(t *testing.T) {
	// Save original LS_COLORS
	originalLSColors := os.Getenv("LS_COLORS")
	defer os.Setenv("LS_COLORS", originalLSColors)

	tests := []struct {
		name      string
		lsColors  string
		wantKeys  []string
		wantColor map[string]string
	}{
		{
			name:      "Empty LS_COLORS",
			lsColors:  "",
			wantKeys:  []string{},
			wantColor: map[string]string{},
		},
		{
			name:     "Basic colors",
			lsColors: "di=01;94:ln=01;36",
			wantKeys: []string{"di", "ln"},
			wantColor: map[string]string{
				"di": "01;94",
				"ln": "01;36",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LS_COLORS", tt.lsColors)
			got := Colors()

			// Check if all expected keys exist with correct values
			for _, key := range tt.wantKeys {
				if color, ok := got[key]; !ok {
					t.Errorf("Key %s not found in color map", key)
				} else if color != tt.wantColor[key] {
					t.Errorf("For key %s, got color %s, want %s", key, color, tt.wantColor[key])
				}
			}

			// Check if there are any unexpected keys
			for key := range got {
				if _, ok := tt.wantColor[key]; !ok {
					t.Errorf("Unexpected key %s in color map", key)
				}
			}
		})
	}
}
