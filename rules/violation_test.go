package rules

import (
	"testing"
)

var isViolationMatchExactTests = []struct {
	namespace  string
	entityType string
	entityName string
	value      string
	rule       []string
	expected   bool
}{
	{"namespace1", "entity1", "name1", "value1", []string{"value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name2:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace2:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity2:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace1:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:!name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace2:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity2:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:!value2"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:*:value1"}, true},
}

func TestIsViolationMatchExact(t *testing.T) {
	for _, vt := range isViolationMatchExactTests {
		actual, err := IsValueMatchExactRule(vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule)
		if err != nil {
			t.Errorf("IsValueMatchExactRule(%s, %s, %s, %s, %s): FAILED! %s",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, err)
		} else if actual != vt.expected {
			t.Errorf("IsValueMatchExactRule(%s, %s, %s, %s, %s): expected %v, actual %v",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, vt.expected, actual)
		}
	}
}

var isViolationMatchLikeTests = []struct {
	namespace  string
	entityType string
	entityName string
	value      string
	rule       []string
	expected   bool
}{
	{"namespace1", "entity1", "name1", "value1", []string{"value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name2:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace2:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity2:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace1:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:!value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace2:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity2:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:!value2"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value2"}, false},

	{"namespace1", "entity1", "name1", "value1___", []string{"value1"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"namespace1:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"!namespace2:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"namespace1:!entity2:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"namespace1:entity1:name1:!value2"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"*:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1___", []string{"namespace1:*:name1:value1"}, true},

	{"namespace1", "entity1", "name1", "___value1___", []string{"value1"}, false},

	{"namespace1", "entity1", "name1", "value1___", []string{"!value1"}, false},
	{"namespace1", "entity1", "name1", "value1___", []string{"!value1"}, false},
}

func TestIsViolationMatchLike(t *testing.T) {
	for _, vt := range isViolationMatchLikeTests {
		actual, err := IsValueMatchLikeRule(vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule)
		if err != nil {
			t.Errorf("IsValueMatchLikeRule(%s, %s, %s, %s, %s): FAILED! %s",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, err)
		} else if actual != vt.expected {
			t.Errorf("IsValueMatchLikeRule(%s, %s, %s, %s, %s): expected %v, actual %v",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, vt.expected, actual)
		}
	}
}

var isViolationMatchContainsTests = []struct {
	namespace  string
	entityType string
	entityName string
	value      string
	rule       []string
	expected   bool
}{
	{"namespace1", "entity1", "name1", "value1", []string{"value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity2:name2:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace2:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity2:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace1:entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity1:name1:value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:!value1"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"!namespace2:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:!entity2:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:entity1:name1:!value2"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "value1", []string{"*:entity1:name1:value2"}, false},
	{"namespace1", "entity1", "name1", "value1", []string{"namespace1:*:name1:value2"}, false},

	{"namespace1", "entity1", "name1", "___value1___", []string{"value1"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"namespace1:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"!namespace2:entity1:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"namespace1:!entity2:name1:value1"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"namespace1:entity1:name1:!value2"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"*:entity1:name1:alue1"}, true},
	{"namespace1", "entity1", "name1", "___value1___", []string{"namespace1:*:name1:value1"}, true},
}

func TestIsViolationMatchContains(t *testing.T) {
	for _, vt := range isViolationMatchContainsTests {
		actual, err := IsValueMatchContainsRule(vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule)
		if err != nil {
			t.Errorf("IsValueMatchContainsRule(%s, %s, %s, %s, %s): FAILED! %s",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, err)
		} else if actual != vt.expected {
			t.Errorf("IsValueMatchContainsRule(%s, %s, %s, %s, %s): expected %v, actual %v",
				vt.namespace, vt.entityType, vt.entityName, vt.value, vt.rule, vt.expected, actual)
		}
	}
}

var isValuesMatchesRequiredRuleTests = []struct {
	namespace  string
	entityType string
	entityName string
	values     map[string]string
	rule       []string
	expected   bool
}{
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:name1:key1"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:name1:key2"}, false},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:name2:key1"}, true},

	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:*:key1"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:*:key2"}, false},

	{"namespace1", "entity2", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:name1:key1"}, true},
	{"namespace2", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:entity1:name1:key1"}, true},

	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"*:entity1:name1:key1"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"*:entity1:name1:key2"}, false},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:*:name1:key1"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:*:name1:key2"}, false},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"*:*:name1:key1"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"*:*:name1:ey2"}, false},

	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"!namespace1:entity1:name1:key2"}, true},
	{"namespace1", "entity1", "name1", map[string]string{"key1": "value1"}, []string{"namespace1:!entity1:name1:key2"}, true},
}

func TestIsValuesMatchesRequiredRule(t *testing.T) {
	for _, vt := range isValuesMatchesRequiredRuleTests {
		actual, source, err := IsValuesMatchesRequiredRule(vt.namespace, vt.entityType, vt.entityName, vt.values, vt.rule)
		if err != nil {
			t.Errorf("IsValuesMatchesRequiredRule(%s, %s, %s, %v, %v): FAILED! %v",
				vt.namespace, vt.entityType, vt.entityName, vt.values, vt.rule, err)
		} else if actual != vt.expected {
			t.Errorf("IsValuesMatchesRequiredRule(%s, %s, %s, %v, %v, : expected %v, actual %v-%s",
				vt.namespace, vt.entityType, vt.entityName, vt.values, vt.rule, vt.expected, actual, source)
		}
	}
}
