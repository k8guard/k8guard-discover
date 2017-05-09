package discover

import (
	lib "github.com/k8guard/k8guardlibs"
	vlib "github.com/k8guard/k8guardlibs/violations"
)

func isNotIgnoredViloation(vt vlib.ViolationType) bool {
	for _, ignoredV := range (lib.Cfg.IgnoredViolations) {
		if string(vt) == ignoredV {
			return false
		}

	}
	return true
}
