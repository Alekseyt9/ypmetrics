package version

import (
	"bytes"
	"fmt"
	"testing"
)

// TestPrint tests the Print method of the Info struct.
func TestPrint(t *testing.T) {
	tests := []struct {
		version  string
		date     string
		commit   string
		expected string
	}{
		{
			version:  "",
			date:     "",
			commit:   "",
			expected: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
		{
			version:  "1.0.0",
			date:     "2024-08-01",
			commit:   "abc123",
			expected: "Build version: 1.0.0\nBuild date: 2024-08-01\nBuild commit: abc123\n",
		},
		{
			version:  "2.0.0",
			date:     "",
			commit:   "def456",
			expected: "Build version: 2.0.0\nBuild date: N/A\nBuild commit: def456\n",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("version=%s date=%s commit=%s", tt.version, tt.date, tt.commit), func(t *testing.T) {
			var buf bytes.Buffer

			info := Info{
				Version: tt.version,
				Date:    tt.date,
				Commit:  tt.commit,
			}
			info.Print(&buf)

			if got := buf.String(); got != tt.expected {
				t.Errorf("Print() = %q, want %q", got, tt.expected)
			}
		})
	}
}
