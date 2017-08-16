package discover

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/k8guard/k8guard-discover/metrics"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
	"github.com/k8guard/k8guardlibs/violations"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/prometheus/client_golang/prometheus"
)

func GetAllPodsFromApi() []v1.Pod {
	pods, err := Clientset.CoreV1().Pods(lib.Cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		lib.Log.Error("error:", err)
		panic(err.Error())
	}

	if lib.Cfg.OutputPodsToFile == true {
		r, _ := json.Marshal(pods.Items)
		err = ioutil.WriteFile("allpodsfromapi.txt", r, 0644)
		if err != nil {
			lib.Log.Error("error:", err)
			panic(err)
		}
	}
	metrics.Update(metrics.ALL_POD_COUNT, len(pods.Items))
	return pods.Items
}

func GetBadPods(allPods []v1.Pod, sendToKafka bool) []lib.Pod {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadPods.Set))
	defer timer.ObserveDuration()

	allBadPodsWitoutOwner := []lib.Pod{}
	badPodsCounter := int64(0)

	cacheAllImages(true)

	for _, kp := range allPods {

		_, createdByAnnotation := kp.Annotations["kubernetes.io/created-by"]
		if createdByAnnotation == true {
			continue
		}

		if isIgnoredNamespace(kp.Namespace) == true || isIgnoredPodPrefix(kp.ObjectMeta.Name) == true {
			continue
		}
		if kp.Status.Phase != "Running" {
			continue
		}
		p := lib.Pod{}
		p.Name = kp.Name
		p.Cluster = lib.Cfg.ClusterName
		p.Namespace = kp.Namespace
		getVolumesWithHostPathForAPod(kp.Spec, &p.ViolatableEntity)
		verifyPodAnnotations(kp.ObjectMeta, &p.ViolatableEntity)
		GetBadContainers(kp.Spec, &p.ViolatableEntity)

		if len(p.Violations) > 0 {

			badPodsCounter += 1
			allBadPodsWitoutOwner = append(allBadPodsWitoutOwner, p)
			if sendToKafka {
				lib.Log.Debug("Sending ", p.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.POD_MESSAGE, p)
				if err != nil {
					panic(err)
				}
			}

		}

	}

	metrics.Update(metrics.BAD_POD_COUNT, int(badPodsCounter))
	metrics.Update(metrics.BAD_POD_WO_OWNER_COUNT, len(allBadPodsWitoutOwner))

	return allBadPodsWitoutOwner

}

// gets a list of entity and fills the host type violations for them
func getVolumesWithHostPathForAPod(spec v1.PodSpec, entity *lib.ViolatableEntity) {
	if isNotIgnoredViolation(violations.HOST_VOLUMES_TYPE) {
		for _, v := range spec.Volumes {
			if v.HostPath != nil {
				entity.Violations = append(entity.Violations, violations.Violation{Source: v.Name, Type: violations.HOST_VOLUMES_TYPE})
			}
		}
	}
}

func isIgnoredPodPrefix(podname string) bool {
	for _, p := range lib.Cfg.IgnoredPodsPrefix {
		if strings.HasPrefix(podname, p) == true {
			return true
		}
	}
	return false
}

// verify whether a specific annotation(s) exists
func verifyPodAnnotations(objectMeta metav1.ObjectMeta, entity *lib.ViolatableEntity) {
	lib.Log.Info("hit verifyPodAnnotations.  not ignored: ", isNotIgnoredViolation(violations.REQUIRED_POD_ANNOTATIONS_TYPE))
	if isNotIgnoredViolation(violations.REQUIRED_POD_ANNOTATIONS_TYPE) {
		lib.Log.Info("   RequiredPodAnnotations: ", fmt.Sprintf("%v", lib.Cfg.RequiredPodAnnotations))
		lib.Log.Info("   objectMeta.Annotations: ", fmt.Sprintf("%v", objectMeta.Annotations))
		for _, a := range lib.Cfg.RequiredPodAnnotations {
			if _, ok := objectMeta.Annotations[a]; !ok {
				entity.Violations = append(entity.Violations, violations.Violation{Source: a, Type: violations.REQUIRED_POD_ANNOTATIONS_TYPE})
			}
		}
	}
}
