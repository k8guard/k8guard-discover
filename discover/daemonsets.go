package discover

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/k8guard/k8guard-discover/messaging"
	"github.com/k8guard/k8guard-discover/rules"

	"github.com/k8guard/k8guard-discover/metrics"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/types"
	"github.com/k8guard/k8guardlibs/violations"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func GetAllDaemonSetFromApi() []v1beta1.DaemonSet {
	daemonsets, err := Clientset.ExtensionsV1beta1().DaemonSets(lib.Cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		lib.Log.Error("error: ", err)
		panic(err.Error())
	}

	if lib.Cfg.OutputPodsToFile == true {
		r, _ := json.Marshal(daemonsets.Items)
		err = ioutil.WriteFile("allpodsfromapi.txt", r, 0644)
		if err != nil {
			lib.Log.Error("error:", err)
			panic(err)
		}
	}

	metrics.Update(metrics.ALL_DAEMONSET_COUNT, len(daemonsets.Items))
	return daemonsets.Items
}

func GetBadDaemonSets(theDaemonSets []v1beta1.DaemonSet, sendToBroker bool) []lib.DaemonSet {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadDaemonSets.Set))
	defer timer.ObserveDuration()

	allBadDaemonSets := []lib.DaemonSet{}

	cacheAllImages(true)

	allBadDaemonSets = append(allBadDaemonSets, verifyRequiredDaemonSets(theDaemonSets, sendToBroker)...)

	for _, kd := range theDaemonSets {

		if isIgnoredNamespace(kd.Namespace) == true || isIgnoredDaemonSet(kd.ObjectMeta.Name) == true {
			continue
		}

		if kd.Status.NumberReady == 0 {
			continue
		}

		d := lib.DaemonSet{}
		d.Name = kd.Name
		d.Cluster = lib.Cfg.ClusterName
		d.Namespace = kd.Namespace
		getVolumesWithHostPathForAPod(kd.Name, kd.Spec.Template.Spec, &d.ViolatableEntity)

		verifyRequiredAnnotations(kd.ObjectMeta.Annotations, &d.ViolatableEntity, "daemonset", violations.REQUIRED_DAEMONSET_ANNOTATIONS_TYPE)
		verifyRequiredLabels(kd.ObjectMeta.Labels, &d.ViolatableEntity, "daemonset", violations.REQUIRED_DAEMONSET_LABELS_TYPE)

		GetBadContainers(kd.Namespace, "daemonset", kd.Spec.Template.Spec, &d.ViolatableEntity)

		if len(d.ViolatableEntity.Violations) > 0 {
			allBadDaemonSets = append(allBadDaemonSets, d)
			if sendToBroker {
				messaging.SendData(types.DAEMONSET_MESSAGE, d.Name, d)
			}
		}
	}
	metrics.Update(metrics.BAD_DAEMONSET_COUNT, len(allBadDaemonSets))
	return allBadDaemonSets
}

func isIgnoredDaemonSet(daemonSetName string) bool {
	for _, d := range lib.Cfg.IgnoredDaemonSets {
		if strings.HasPrefix(daemonSetName, d) == true {
			return true
		}
	}
	return false
}

func verifyRequiredDaemonSets(theDaemonSets []v1beta1.DaemonSet, sendToBroker bool) []lib.DaemonSet {
	entityType := "daemonset"
	badDaemonSets := []lib.DaemonSet{}

	for _, ns := range GetAllNamespaceFromApi() {
		if rules.IsNotIgnoredViolation(ns.Name, entityType, "*", violations.REQUIRED_ENTITIES_TYPE) {
			for _, a := range lib.Cfg.RequiredEntities {
				rule := strings.Split(a, ":")

				// does the rule apply to this namespace and entity type?
				if !(rules.Exact(ns.Name, rule[0]) && rules.Exact(entityType, rule[1])) {
					continue
				}

				found := false
				for _, kd := range theDaemonSets {
					if rules.Exact(kd.ObjectMeta.Namespace, rule[0]) && rules.Exact(kd.ObjectMeta.Name, rule[2]) {
						found = true
						break
					}
				}

				if !found {
					ds := lib.DaemonSet{}
					ds.Name = rule[2]
					ds.Cluster = lib.Cfg.ClusterName
					ds.Namespace = ns.Name
					ds.ViolatableEntity.Violations = append(ds.ViolatableEntity.Violations, violations.Violation{Source: rule[2], Type: violations.REQUIRED_DAEMONSETS_TYPE})
					badDaemonSets = append(badDaemonSets, ds)

					if sendToBroker {
						messaging.SendData(types.DAEMONSET_MESSAGE, ds.Name, ds)
					}
				}
			}
		}
	}

	return badDaemonSets
}
