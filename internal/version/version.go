package version

import (
	"runtime/debug"
	"time"

	"golang.org/x/mod/module"
)

var (
	version = "v0.0.0"
)

func init() {
	if !module.IsPseudoVersion(version) {

		bi, ok := debug.ReadBuildInfo()
		if ok {
			var rev string
			var t time.Time

			for _, s := range bi.Settings {
				if s.Key == "vcs.revision" {
					if len(s.Value) >= 12 {
						rev = s.Value[0:12]
					}
				}
				if s.Key == "vcs.time" {
					t, _ = time.Parse(time.RFC3339, s.Value)
				}
			}

			version = module.PseudoVersion("", "0.0.0", t, rev)
		}
	}

}

func Version() string {
	return version
}
