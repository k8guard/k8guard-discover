package discover

import (
	lib "github.com/k8guard/k8guardlibs"
	vlib "github.com/k8guard/k8guardlibs/violations"
	"strings"
)

func isNotIgnoredViolation(name string, vt vlib.ViolationType) bool {
	for _, ignoredV := range lib.Cfg.IgnoredViolations {
		config := strings.Split(ignoredV, ":")
		if len(config) > 1 {
			if name == config[0] && string(vt) == config[1] {
				return false
			}
		} else {
			if string(vt) == ignoredV {
				return false
			}
		}

	}
	return true
}
