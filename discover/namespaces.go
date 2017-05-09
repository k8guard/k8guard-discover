package discover

import (
	lib "github.com/k8guard/k8guardlibs"
)


func isIgnoredNamespace(namespace string) bool {
	for _, n := range lib.Cfg.IgnoredNamespaces {
		if n == namespace {
			return true
		}
	}
	return false
}

