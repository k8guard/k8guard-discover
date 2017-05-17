package discover

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"k8s.io/client-go/pkg/api/v1"
	lib "github.com/k8guard/k8guardlibs"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
	"github.com/k8guard/k8guardlibs/violations"
	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func GetAllDeployFromApi() []v1beta1.Deployment {
	deploys, err := Clientset.Deployments(lib.Cfg.Namespace).List(v1.ListOptions{})
	if err != nil {
		lib.Log.Error("error:", err)
		panic(err.Error())
	}

	if (lib.Cfg.OutputPodsToFile == true) {
		r, _ := json.Marshal(deploys.Items)
		err = ioutil.WriteFile("deployments.txt", r, 0644)
		if err != nil {
			lib.Log.Error("error:", err)
			panic(err)
		}
	}
	metrics.Update(metrics.ALL_DEPLOYMENT_COUNT, len(deploys.Items))
	return deploys.Items
}

func GetBadDeploys(theDeploys []v1beta1.Deployment, sendToKafka bool) []lib.Deployment {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadDeploys.Set))
	defer timer.ObserveDuration()

	allBadDeploys := []lib.Deployment{}

	cacheAllImages(true)

	for _, kd := range theDeploys {
		if isIgnoredNamespace(kd.Namespace) == true || isIgnoredDeployment(kd.ObjectMeta.Name) == true {
			continue
		}
		if kd.Status.Replicas == 0 {
			continue
		}
		d := lib.Deployment{}
		d.Name = kd.Name
		d.Cluster = lib.Cfg.ClusterName
		d.Namespace = kd.Namespace
		getVolumesWithHostPathForAPod(kd.Spec.Template.Spec, &d.ViolatableEntity)
		GetBadContainers(kd.Spec.Template.Spec, &d.ViolatableEntity)
		if isValidReplicaSize(kd) == false && isNotIgnoredViloation(violations.SINGLE_REPLICA_TYPE){
			d.Violations = append(d.Violations, violations.Violation{Source: kd.Name, Type: violations.SINGLE_REPLICA_TYPE})
		}

		if len(d.ViolatableEntity.Violations) > 0 {
			allBadDeploys = append(allBadDeploys, d)
			if sendToKafka {
				lib.Log.Debug("Sending ", d.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.DEPLOYMENT_MESSAGE, d)
				if err != nil {
					panic(err)
				}
			}
		}

	}
	metrics.Update(metrics.BAD_DEPLOYMENT_COUNT, len(allBadDeploys))
	return allBadDeploys
}

func isValidReplicaSize(deployment v1beta1.Deployment) (bool) {
	if *deployment.Spec.Replicas == 1 {
		return false
	}
	return true
}

func isIgnoredDeployment(deploymentName string) bool {
	for _, d := range lib.Cfg.IgnoredDeployments {
		if strings.HasPrefix(deploymentName, d) == true {
			return true
		}
	}
	return false
}
