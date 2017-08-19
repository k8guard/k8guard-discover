package rules

import (
	"errors"
	"fmt"
	"strings"

	lib "github.com/k8guard/k8guardlibs"
	vlib "github.com/k8guard/k8guardlibs/violations"
)

const NEGATED_RULE = "!"
const ALLOW_ALL_RULE = "*"

func IsNotIgnoredViolation(namespace string, entity string, vt vlib.ViolationType) bool {
	ignored, _ := IsValueMatchExactRule(namespace, entity, string(vt), lib.Cfg.IgnoredViolations)
	return !ignored
}

func IsValuesMatchesRequiredRule(namespace string, entity string, values map[string]string, violationConfig []string) (bool, string, error) {
	for _, r := range violationConfig {
		rule := strings.Split(r, ":")

		// if len(rule) != 3 {
		// 	return false, "", errors.New(fmt.Sprintf("Incorrect format for violation rule: %v", violationConfig))
		// }
		if len(rule) == 1 {
			// a length of 1 is for backwards compatibility, and is a global config (across all namespaces and entity types
			if _, ok := values[rule[0]]; !ok {
				return false, rule[0], nil
			}
		} else if len(rule) == 3 {
			// does the rule apply to this namespace and entity type?
			if !(Exact(namespace, rule[0]) && Exact(entity, rule[1])) {
				continue
			}
			if _, ok := values[rule[2]]; !ok {
				return false, rule[2], nil
			}
		} else {
			return false, "", fmt.Errorf("Incorrect format for violation rule: %v", violationConfig)
		}
	}
	return true, "", nil
}

func IsValueMatchExactRule(namespace string, entity string, value string, violationConfig []string) (bool, error) {
	for _, ignoredV := range violationConfig {
		rule := strings.Split(ignoredV, ":")

		if len(rule) == 1 {
			// a length of 1 is for backwards compatibility, and is a global config (across all namespaces and entity types
			return Exact(value, rule[0]), nil
		} else if len(rule) == 3 {
			// denotes namespace, entity and value
			return Exact(namespace, rule[0]) && Exact(entity, rule[1]) && Exact(value, rule[2]), nil
		} else {
			return false, fmt.Errorf("Incorrect format for violation rule: %v", violationConfig)
		}
	}
	return false, nil
}

func IsValueMatchLikeRule(namespace string, entity string, value string, violationConfig []string) (bool, error) {
	for _, ignoredV := range violationConfig {
		rule := strings.Split(ignoredV, ":")

		if len(rule) == 1 {
			// a length of 1 is for backwards compatibility, and is a global config (across all namespaces and entity types
			return Like(value, rule[0]), nil
		} else if len(rule) == 3 {
			// denotes namespace, entity and value
			return Exact(namespace, rule[0]) && Exact(entity, rule[1]) && Like(value, rule[2]), nil
		} else {
			return false, errors.New(fmt.Sprintf("Incorrect format for violation rule: %v", violationConfig))
		}
	}
	return false, nil
}

func IsValueMatchContainsRule(namespace string, entity string, value string, violationConfig []string) (bool, error) {
	for _, ignoredV := range violationConfig {
		rule := strings.Split(ignoredV, ":")

		if len(rule) == 1 {
			// a length of 1 is for backwards compatibility, and is a global config (across all namespaces and entity types
			return Contains(value, rule[0]), nil
		} else if len(rule) == 3 {
			// denotes namespace, entity and value
			return Exact(namespace, rule[0]) && Exact(entity, rule[1]) && Contains(value, rule[2]), nil
		} else {
			return false, fmt.Errorf("Incorrect format for violation rule: %v", violationConfig)
		}
	}
	return false, nil
}

func Exact(value string, rule string) bool {
	if strings.HasPrefix(rule, NEGATED_RULE) {
		return value != strings.Replace(rule, NEGATED_RULE, "", -1)
	} else if rule == ALLOW_ALL_RULE {
		return true
	} else {
		return value == rule
	}
}

func Like(value string, rule string) bool {
	if strings.HasPrefix(rule, NEGATED_RULE) {
		return !strings.HasPrefix(value, strings.Replace(rule, NEGATED_RULE, "", -1))
	} else if rule == ALLOW_ALL_RULE {
		return true
	} else {
		return strings.HasPrefix(value, rule)
	}
}

func Contains(value string, rule string) bool {
	if strings.HasPrefix(rule, NEGATED_RULE) {
		return value != strings.Replace(rule, NEGATED_RULE, "", -1)
	} else if rule == ALLOW_ALL_RULE {
		return true
	} else {
		return strings.Contains(value, rule)
	}
}
