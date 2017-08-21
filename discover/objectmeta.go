package discover

import (
	"fmt"

	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/violations"
)

// verify whether a specific annotation(s) exist
func verifyRequiredAnnotations(annotations map[string]string, entity *lib.ViolatableEntity, entityType string, violation violations.ViolationType) {
	lib.Log.Info(fmt.Sprintf("Verifying namespace: %s, entity: %s, against: %v", entity.Namespace, entityType, lib.Cfg.RequiredAnnotations))
	if rules.IsNotIgnoredViolation(entity.Namespace, entityType, entity.Name, violations.REQUIRED_ANNOTATIONS_TYPE) {
		lib.Log.Info("   not ignored")
		if found, source, _ := rules.IsValuesMatchesRequiredRule(entity.Namespace, entityType, entity.Name, annotations, lib.Cfg.RequiredAnnotations); !found {
			entity.Violations = append(entity.Violations, violations.Violation{Source: source, Type: violation})
		}
	}
}

// verify whether a specific label(s) exists
func verifyRequiredLabels(labels map[string]string, entity *lib.ViolatableEntity, entityType string, violation violations.ViolationType) {
	if rules.IsNotIgnoredViolation(entity.Namespace, entityType, entity.Name, violations.REQUIRED_LABELS_TYPE) {
		if found, source, _ := rules.IsValuesMatchesRequiredRule(entity.Namespace, entityType, entity.Name, labels, lib.Cfg.RequiredLabels); !found {
			entity.Violations = append(entity.Violations, violations.Violation{Source: source, Type: violation})
		}
	}
}
