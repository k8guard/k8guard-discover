package discover

import (
	"fmt"

	"github.com/k8guard/k8guard-discover/caching"
	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/violations"
	"k8s.io/client-go/pkg/api/v1"
)

func getContainerImageSize(imageName string) int64 {
	size, _ := caching.GetAsInt(fmt.Sprintf("image_%s", imageName))
	return size
}

//GetBadContainers Gets Containers with invalid size, invalid repo, or cababilities/privileged
func GetBadContainers(namespace string, entityType string, spec v1.PodSpec, entity *lib.ViolatableEntity) {
	for _, c := range spec.Containers {
		cImageSize := getContainerImageSize(c.Image)
		if isValidImageRepo(namespace, entityType, c.Name, c.Image) == false && rules.IsNotIgnoredViolation(namespace, entityType, c.Name, violations.IMAGE_REPO_TYPE) {
			entity.Violations = append(entity.Violations, violations.Violation{Source: c.Image, Type: violations.IMAGE_REPO_TYPE})
		}

		if isValidImageSize(cImageSize) == false && rules.IsNotIgnoredViolation(namespace, entityType, c.Name, violations.IMAGE_SIZE_TYPE) {
			entity.Violations = append(entity.Violations, violations.Violation{Source: c.Image, Type: violations.IMAGE_SIZE_TYPE})
		}

		if c.SecurityContext != nil {
			if c.SecurityContext.Privileged != nil && *c.SecurityContext.Privileged && rules.IsNotIgnoredViolation(namespace, entityType, c.Name, violations.PRIVILEGED_TYPE) {
				entity.Violations = append(entity.Violations, violations.Violation{Source: c.Image, Type: violations.PRIVILEGED_TYPE})
			}

			//Check if containers have extra capabilities set like NET_ADMIN...
			if c.SecurityContext.Capabilities != nil && len(c.SecurityContext.Capabilities.Add) > 0 && rules.IsNotIgnoredViolation(namespace, entityType, c.Name, violations.CAPABILITIES_TYPE) {
				entity.Violations = append(entity.Violations, violations.Violation{Source: c.Image, Type: violations.CAPABILITIES_TYPE})
			}
		}
	}
}
