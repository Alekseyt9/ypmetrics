package version

import (
	"fmt"
	"io"
)

// Info holds build information.
type Info struct {
	Version string
	Date    string
	Commit  string
}

func (i Info) Print(w io.Writer) {
	fmt.Fprintf(w, "Build version: %s\n", getBuildInfo(i.Version))
	fmt.Fprintf(w, "Build date: %s\n", getBuildInfo(i.Date))
	fmt.Fprintf(w, "Build commit: %s\n", getBuildInfo(i.Commit))
}

// getBuildInfo returns the build information or "N/A" if it is empty.
func getBuildInfo(info string) string {
	if info == "" {
		return "N/A"
	}
	return info
}
