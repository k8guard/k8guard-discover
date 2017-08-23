package discover

import (
	"encoding/json"
	"strings"

	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
	"github.com/k8guard/k8guardlibs/violations"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func isIgnoredNamespace(namespace string) bool {
	for _, n := range lib.Cfg.IgnoredNamespaces {
		if n == namespace {
			return true
		}
	}
	return false
}

func GetAllNamespaceFromApi() []v1.Namespace {
	namespaces := Clientset.Namespaces()

	namespaceList, err := namespaces.List(metav1.ListOptions{})

	if err != nil {
		lib.Log.Error("error: ", err)
		panic(err.Error())
	}

	metrics.Update(metrics.ALL_NAMESPACE_COUNT, len(namespaceList.Items))

	return namespaceList.Items
}

func GetBadNamespaces(theNamespaces []v1.Namespace, sendToKafka bool) []lib.Namespace {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadNamespaces.Set))
	defer timer.ObserveDuration()

	allBadNamespaces := []lib.Namespace{}

	allBadNamespaces = append(allBadNamespaces, verifyRequiredNamespaces(theNamespaces, sendToKafka)...)

	for _, kn := range theNamespaces {
		if isIgnoredNamespace(kn.Namespace) == true {
			continue
		}
		n := lib.Namespace{}
		n.Name = kn.Name
		n.Namespace = kn.Name
		n.Cluster = lib.Cfg.ClusterName
		// this one feels weird but to be consistent

		if hasOwnerAnnotation(kn, lib.Cfg.AnnotationFormatForEmails) == false &&
			hasOwnerAnnotation(kn, lib.Cfg.AnnotationFormatForChatIds) == false &&
			rules.IsNotIgnoredViolation(kn.Name, "namespace", kn.Name, violations.NO_OWNER_ANNOTATION_TYPE) {
			jsonString, err := json.Marshal(kn.Annotations)
			if err != nil {
				lib.Log.Error("Can not convert annotation to a valid json ", err)

			}
			n.Violations = append(n.Violations, violations.Violation{Source: string(jsonString), Type: violations.NO_OWNER_ANNOTATION_TYPE})
		}

		verifyRequiredAnnotations(kn.ObjectMeta.Annotations, &n.ViolatableEntity, "namespace", violations.REQUIRED_NAMESPACE_ANNOTATIONS_TYPE)
		verifyRequiredLabels(kn.ObjectMeta.Labels, &n.ViolatableEntity, "namespace", violations.REQUIRED_NAMESPACE_LABELS_TYPE)
		verifyRequiredResourceQuota(&n.ViolatableEntity)

		if len(n.ViolatableEntity.Violations) > 0 {
			allBadNamespaces = append(allBadNamespaces, n)
			if sendToKafka {
				lib.Log.Debug("Sending ", n.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.NAMESPACE_MESSAGE, n)
				if err != nil {
					panic(err)
				}
			}
		}

	}
	metrics.Update(metrics.BAD_NAMESPACE_COUNT, len(allBadNamespaces))
	return allBadNamespaces
}

func hasOwnerAnnotation(namespace v1.Namespace, annotationKind string) bool {
	teamString, ok := namespace.Annotations[annotationKind]
	if ok {
		team := strings.Split(teamString, ",")
		if len(team) > 0 {
			return true
		}
	}
	return false
}

func verifyRequiredNamespaces(theNamespaces []v1.Namespace, sendToKafka bool) []lib.Namespace {
	badNamespaces := []lib.Namespace{}

	for _, a := range lib.Cfg.RequiredEntities {
		rule := strings.Split(a, ":")

		// does the rule apply to this entity type?
		if !rules.Exact("namespace", rule[1]) {
			continue
		}

		found := false
		for _, kn := range theNamespaces {
			if rules.Exact(kn.ObjectMeta.Name, rule[3]) {
				found = true
				break
			}
		}

		if !found {
			ns := lib.Namespace{}
			ns.Name = rule[3]
			ns.Cluster = lib.Cfg.ClusterName
			ns.Namespace = ns.Name
			ns.ViolatableEntity.Violations = append(ns.ViolatableEntity.Violations, violations.Violation{Source: rule[3], Type: violations.REQUIRED_NAMESPACES_TYPE})
			badNamespaces = append(badNamespaces, ns)

			if sendToKafka {
				lib.Log.Debug("Sending ", ns.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.NAMESPACE_MESSAGE, ns)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return badNamespaces
}

func verifyRequiredResourceQuota(entity *lib.ViolatableEntity) {

	entityType := "resourcequota"
	if !rules.IsNotIgnoredViolation(entity.Namespace, entityType, "*", violations.REQUIRED_ENTITIES_TYPE) {
		return
	}

	resourcequotas, err := Clientset.CoreV1Client.ResourceQuotas(entity.Namespace).List(metav1.ListOptions{})
	if err != nil {
		lib.Log.Error("error: ", err)
		panic(err.Error())
	}

	for _, a := range lib.Cfg.RequiredEntities {
		rule := strings.Split(a, ":")

		if len(rule) != 3 {
			continue
		}

		// does the rule apply to this namespace and entity type (resourcequota)?
		if !(rules.Exact(entity.Namespace, rule[0]) && rules.Exact(entityType, rule[1])) {
			continue
		}

		found := false
		for _, krq := range resourcequotas.Items {
			if rules.Exact(krq.ObjectMeta.Name, rule[2]) {
				found = true
				break
			}
		}

		if !found {
			entity.Violations = append(entity.Violations, violations.Violation{Source: entity.Namespace, Type: violations.REQUIRED_RESOURCEQUOTA_TYPE})
		}
	}
}
