package discover

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/k8guard/k8guard-discover/metrics"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
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

func GetBadDaemonSets(theDaemonSets []v1beta1.DaemonSet, sendToKafka bool) []lib.DaemonSet {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadDaemonSets.Set))
	defer timer.ObserveDuration()

	allBadDaemonSets := []lib.DaemonSet{}

	cacheAllImages(true)

	verifyRequiredDaemonSets(theDaemonSets)

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
		verifyRequiredAnnotations(kd.ObjectMeta, &d.ViolatableEntity, violations.REQUIRED_POD_ANNOTATIONS_TYPE, lib.Cfg.RequiredPodAnnotations)
		verifyRequiredLabels(kd.ObjectMeta, &d.ViolatableEntity, violations.REQUIRED_POD_LABELS_TYPE, lib.Cfg.RequiredPodLabels)
		GetBadContainers(kd.Name, kd.Spec.Template.Spec, &d.ViolatableEntity)

		if len(d.ViolatableEntity.Violations) > 0 {
			allBadDaemonSets = append(allBadDaemonSets, d)
			if sendToKafka {
				lib.Log.Debug("Sending ", d.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.DAEMONSET_MESSAGE, d)
				if err != nil {
					panic(err)
				}
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

func verifyRequiredDaemonSets(theDaemonSets []v1beta1.DaemonSet) {
	if isNotIgnoredViolation("", violations.REQUIRED_DAEMONSETS_TYPE) {
		for _, a := range lib.Cfg.RequiredDaemonSets {
			required := strings.Split(a, ":")
			ds := lib.DaemonSet{}
			found := false
			for _, kd := range theDaemonSets {
				if (required[0] == kd.Namespace) && (required[1] == kd.ObjectMeta.Name) {
					found = true
					break
				}
			}
			if !found {
				ds.Name = required[1]
				ds.Cluster = lib.Cfg.ClusterName
				ds.Namespace = required[0]
				ds.ViolatableEntity.Violations = append(ds.ViolatableEntity.Violations, violations.Violation{Source: a, Type: violations.REQUIRED_DAEMONSETS_TYPE})
			}
		}
	}
}
