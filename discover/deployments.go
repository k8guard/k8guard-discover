package discover

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/k8guard/k8guard-discover/messaging"
	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/types"
	"github.com/k8guard/k8guardlibs/violations"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

func GetAllDeployFromApi() []v1beta1.Deployment {
	deploys, err := Clientset.AppsV1beta1().Deployments(lib.Cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		lib.Log.Error("error: ", err)
		panic(err.Error())
	}

	if lib.Cfg.OutputPodsToFile == true {
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

func GetBadDeploys(theDeploys []v1beta1.Deployment, sendToBroker bool) []lib.Deployment {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadDeploys.Set))
	defer timer.ObserveDuration()

	allBadDeploys := []lib.Deployment{}

	cacheAllImages(true)

	allBadDeploys = append(allBadDeploys, verifyRequiredDeployments(theDeploys, sendToBroker)...)

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
		getVolumesWithHostPathForAPod(kd.Name, kd.Spec.Template.Spec, &d.ViolatableEntity)

		verifyRequiredAnnotations(kd.ObjectMeta.Annotations, &d.ViolatableEntity, "deployment", violations.REQUIRED_DEPLOYMENT_ANNOTATIONS_TYPE)
		verifyRequiredLabels(kd.ObjectMeta.Labels, &d.ViolatableEntity, "deployment", violations.REQUIRED_DEPLOYMENT_LABELS_TYPE)

		verifyRequiredAnnotations(kd.Spec.Template.ObjectMeta.Annotations, &d.ViolatableEntity, "pod", violations.REQUIRED_POD_ANNOTATIONS_TYPE)
		verifyRequiredLabels(kd.Spec.Template.ObjectMeta.Labels, &d.ViolatableEntity, "pod", violations.REQUIRED_POD_LABELS_TYPE)

		GetBadContainers(kd.Namespace, "deployment", kd.Spec.Template.Spec, &d.ViolatableEntity)
		if isValidReplicaSize(kd) == false && rules.IsNotIgnoredViolation(kd.Namespace, "deployment", kd.Name, violations.SINGLE_REPLICA_TYPE) {
			d.Violations = append(d.Violations, violations.Violation{Source: kd.Name, Type: violations.SINGLE_REPLICA_TYPE})
		}

		if len(d.ViolatableEntity.Violations) > 0 {
			allBadDeploys = append(allBadDeploys, d)
			if sendToBroker {
				messaging.SendData(types.DEPLOYMENT_MESSAGE, d.Name, d)
			}
		}

	}
	metrics.Update(metrics.BAD_DEPLOYMENT_COUNT, len(allBadDeploys))
	return allBadDeploys
}

func isValidReplicaSize(deployment v1beta1.Deployment) bool {
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

func verifyRequiredDeployments(theDeployments []v1beta1.Deployment, sendToBroker bool) []lib.Deployment {
	entityType := "deployment"
	badDeployments := []lib.Deployment{}

	for _, ns := range GetAllNamespaceFromApi() {
		if rules.IsNotIgnoredViolation(ns.Name, entityType, "*", violations.REQUIRED_ENTITIES_TYPE) {
			for _, a := range lib.Cfg.RequiredEntities {
				rule := strings.Split(a, ":")

				// does the rule apply to this namespace and entity type?
				if !(rules.Exact(ns.Name, rule[0]) && rules.Exact(entityType, rule[1])) {
					continue
				}

				found := false
				for _, kd := range theDeployments {
					if rules.Exact(kd.ObjectMeta.Namespace, rule[0]) && rules.Exact(kd.ObjectMeta.Name, rule[2]) {
						found = true
						break
					}
				}

				if !found {
					d := lib.Deployment{}
					d.Name = rule[2]
					d.Cluster = lib.Cfg.ClusterName
					d.Namespace = ns.Name
					d.ViolatableEntity.Violations = append(d.ViolatableEntity.Violations, violations.Violation{Source: rule[2], Type: violations.REQUIRED_DEPLOYMENTS_TYPE})
					badDeployments = append(badDeployments, d)

					if sendToBroker {
						messaging.SendData(types.DEPLOYMENT_MESSAGE, d.Name, d)
					}
				}
			}
		}
	}

	return badDeployments
}
