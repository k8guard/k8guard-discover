package discover

import (
	"strings"

	"github.com/k8guard/k8guard-discover/messaging"
	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/messaging/types"
	"github.com/k8guard/k8guardlibs/violations"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func GetAllIngressFromApi() []v1beta1.Ingress {
	ingresses, err := Clientset.Ingresses(lib.Cfg.Namespace).List(metav1.ListOptions{})
	if err != nil {
		lib.Log.Error("error:", err)
		panic(err.Error())
	}
	metrics.Update(metrics.ALL_INGRESSES_COUNT, len(ingresses.Items))
	return ingresses.Items
}

func GetBadIngresses(allIngresses []v1beta1.Ingress, sendToBroker bool) []lib.Ingress {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadIngresses.Set))
	defer timer.ObserveDuration()

	allBadIngresses := []lib.Ingress{}
	for _, kin := range allIngresses {

		if isIgnoredNamespace(kin.Namespace) == true {
			continue
		}

		in := lib.Ingress{}
		in.Name = kin.Name
		in.Cluster = lib.Cfg.ClusterName
		in.Namespace = kin.Namespace

		for _, rule := range kin.Spec.Rules {
			isBadIngressRule(kin.Namespace, kin.Name, rule, &in)
		}

		if len(in.Violations) > 0 {
			allBadIngresses = append(allBadIngresses, in)
			if sendToBroker {
				messaging.SendData(types.INGRESS_MESSAGE, in.Name, in)
			}
		}
	}

	metrics.Update(metrics.BAD_INGRESSES_COUNT, len(allBadIngresses))
	return allBadIngresses

}

func isBadIngressRule(namespace string, name string, rule v1beta1.IngressRule, ingress *lib.Ingress) bool {

	if rules.IsNotIgnoredViolation(namespace, "ingress", name, violations.INGRESS_HOST_INVALID_TYPE) {
		for _, s := range lib.Cfg.IngressMustContain {
			if strings.Contains(rule.Host, s) != true {
				ingress.Violations = append(ingress.Violations, violations.Violation{Source: rule.Host, Type: violations.INGRESS_HOST_INVALID_TYPE})
				return true
			}
		}

		for _, s := range lib.Cfg.IngressMustNOTContain {
			if strings.Contains(rule.Host, s) == true {
				ingress.Violations = append(ingress.Violations, violations.Violation{Source: rule.Host, Type: violations.INGRESS_HOST_INVALID_TYPE})
				return true
			}
		}

		if len(lib.Cfg.ApprovedIngressSuffixes) > 0 {
			approvedSuffix := false
			for _, s := range lib.Cfg.ApprovedIngressSuffixes {
				if strings.HasSuffix(rule.Host, s) == true {
					approvedSuffix = true
					break
				}
			}
			if approvedSuffix == false {
				ingress.Violations = append(ingress.Violations, violations.Violation{Source: rule.Host, Type: violations.INGRESS_HOST_INVALID_TYPE})
				return true
			}
		}
	}
	return false

}
