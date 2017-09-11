package metrics

import (
	"net/http"
	"strconv"

	"github.com/k8guard/k8guard-discover/caching"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//  To add a new prometheus metric, you need to add it to 4 places, all of them in this file.
//  1- create const.
//  2- create a gauge var.
//  3- register it.
//  4- and add it to the handler.

// step 1
const (
	ALL_NAMESPACE_COUNT  string = "k8guard_all_namespace_count"
	ALL_DEPLOYMENT_COUNT string = "k8guard_all_deployment_count"
	ALL_DAEMONSET_COUNT  string = "k8guard_all_daemonset_count"
	ALL_POD_COUNT        string = "k8guard_all_pod_count"
	ALL_IMAGE_COUNT      string = "k8guard_all_image_count"
	ALL_INGRESSES_COUNT  string = "k8guard_all_ingresses_count"
	ALL_JOB_COUNT        string = "k8guard_all_job_count"
	ALL_CRONJOB_COUNT    string = "k8guard_all_cronjob_count"

	BAD_NAMESPACE_COUNT    string = "k8guard_bad_namespace_count"
	BAD_POD_COUNT          string = "k8guard_bad_pod_count"
	BAD_POD_WO_OWNER_COUNT string = "k8guard_bad_pod_wo_owner_count"
	BAD_DEPLOYMENT_COUNT   string = "k8guard_bad_deployment_count"
	BAD_DAEMONSET_COUNT    string = "k8guard_bad_daemonset_count"
	BAD_INGRESSES_COUNT    string = "k8guard_bad_ingresses_count"
	BAD_JOB_COUNT          string = "k8guard_bad_job_count"
	BAD_JOB_WO_OWNER_COUNT string = "k8guard_bad_job_wo_owner_count"
	BAD_CRONJOB_COUNT      string = "k8guard_bad_cronjob_count"

	METRIC_EXPIRE_SECONDS int32 = 43200
)

// step 2
var (
	AllNamespaceCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_NAMESPACE_COUNT,
			Help: "the number of all namespaces",
		},
	)

	AllDeploymentCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_DEPLOYMENT_COUNT,
			Help: "the number of all deployments",
		},
	)

	AllDaemonSetCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_DAEMONSET_COUNT,
			Help: "the number of all daemonsets",
		},
	)

	AllPodCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_POD_COUNT,
			Help: "the number of all pods",
		},
	)

	AllImageCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_IMAGE_COUNT,
			Help: "the number of all pods",
		},
	)

	AllIngressesGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_INGRESSES_COUNT,
			Help: "the number of all ingresses",
		},
	)

	AllJobsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_JOB_COUNT,
			Help: "the number of all jobs",
		},
	)

	AllCronJobsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_CRONJOB_COUNT,
			Help: "the number of all cron jobs",
		},
	)

	BadNamespaceCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_NAMESPACE_COUNT,
			Help: "the number of namespaces without correct annotation",
		},
	)

	BadPodCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_POD_COUNT,
			Help: "the number of bad pods",
		},
	)

	BadPodWoOwnerGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_POD_WO_OWNER_COUNT,
			Help: "the number of bad pods without an owner",
		},
	)

	BadDeploymentCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_DEPLOYMENT_COUNT,
			Help: "the number of bad deployments",
		},
	)

	BadDaemonSetCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_DAEMONSET_COUNT,
			Help: "the number of bad daemonsets",
		},
	)
	BadIngressesCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_INGRESSES_COUNT,
			Help: "the number of bad ingresses",
		},
	)

	BadJobCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_JOB_COUNT,
			Help: "the number of bad jobs",
		},
	)

	BadJobWoOwnerGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_JOB_WO_OWNER_COUNT,
			Help: "the number of bad jobs without an owner",
		},
	)

	BadCronJobCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: BAD_CRONJOB_COUNT,
			Help: "the number of bad cronjobs",
		},
	)

	// Function Metrics
	FNGetBadNamespaces = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_namespaces_duration",
			Help: "time took for GetBadNamespaces",
		},
	)

	FNCacheAllImagesGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_cache_all_images_duration",
			Help: "time took for cacheAllImages",
		},
	)
	FNGetBadDeploys = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_deploys_duration",
			Help: "time took for GetBadDeploys",
		},
	)
	FNGetBadDaemonSets = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_daemonsets_duration",
			Help: "time took for GetBadDaemonSets",
		},
	)
	FNGetBadPods = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_pods_duration",
			Help: "time took for GetBadPods",
		},
	)
	FNGetBadIngresses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_ingresses_duration",
			Help: "time took for GetBadIngresss",
		},
	)
	FNGetBadJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_jobs_duration",
			Help: "time took for GetBadJobs",
		},
	)
	FNGetBadCronJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "fn_get_bad_cronjobs_duration",
			Help: "time took for GetBadCronJobs",
		},
	)
)

// step 3
func PromRegister() {
	prometheus.MustRegister(AllNamespaceCountGauge)
	prometheus.MustRegister(AllPodCountGauge)
	prometheus.MustRegister(AllImageCountGauge)
	prometheus.MustRegister(AllDeploymentCountGauge)
	prometheus.MustRegister(AllDaemonSetCountGauge)
	prometheus.MustRegister(AllIngressesGauge)
	prometheus.MustRegister(AllJobsGauge)
	prometheus.MustRegister(AllCronJobsGauge)

	prometheus.MustRegister(BadNamespaceCountGauge)
	prometheus.MustRegister(BadDeploymentCountGauge)
	prometheus.MustRegister(BadDaemonSetCountGauge)
	prometheus.MustRegister(BadPodCountGauge)
	prometheus.MustRegister(BadPodWoOwnerGauge)
	prometheus.MustRegister(BadIngressesCountGauge)
	prometheus.MustRegister(BadJobCountGauge)
	prometheus.MustRegister(BadJobWoOwnerGauge)
	prometheus.MustRegister(BadCronJobCountGauge)

	// Function Metrics
	prometheus.MustRegister(FNCacheAllImagesGauge)
	prometheus.MustRegister(FNGetBadDeploys)
	prometheus.MustRegister(FNGetBadDaemonSets)
	prometheus.MustRegister(FNGetBadIngresses)
	prometheus.MustRegister(FNGetBadPods)
	prometheus.MustRegister(FNGetBadJobs)
	prometheus.MustRegister(FNGetBadCronJobs)

}

// step 4
func PromMetricsHandler(w http.ResponseWriter, r *http.Request) {
	myList := []string{
		ALL_NAMESPACE_COUNT, ALL_DEPLOYMENT_COUNT, ALL_DAEMONSET_COUNT, ALL_POD_COUNT, ALL_IMAGE_COUNT, ALL_INGRESSES_COUNT, ALL_JOB_COUNT,
		ALL_CRONJOB_COUNT, BAD_NAMESPACE_COUNT, BAD_POD_COUNT, BAD_DEPLOYMENT_COUNT, BAD_DAEMONSET_COUNT, BAD_POD_WO_OWNER_COUNT, BAD_INGRESSES_COUNT,
		BAD_JOB_COUNT, BAD_JOB_WO_OWNER_COUNT, BAD_CRONJOB_COUNT,
	}
	for _, i := range myList {
		count, err := caching.GetAsInt(i)
		if err != nil {
			count = 0
		}
		parsedCount, _ := strconv.ParseFloat(string(count), 64)

		switch i {
		case ALL_NAMESPACE_COUNT:
			AllNamespaceCountGauge.Set(parsedCount)
			break

		case ALL_POD_COUNT:
			AllPodCountGauge.Set(parsedCount)
			break
		case ALL_DEPLOYMENT_COUNT:
			AllDeploymentCountGauge.Set(parsedCount)
			break
		case ALL_DAEMONSET_COUNT:
			AllDaemonSetCountGauge.Set(parsedCount)
			break
		case ALL_IMAGE_COUNT:
			AllImageCountGauge.Set(parsedCount)
			break
		case ALL_INGRESSES_COUNT:
			AllIngressesGauge.Set(parsedCount)
			break
		case ALL_JOB_COUNT:
			AllJobsGauge.Set(parsedCount)
			break
		case ALL_CRONJOB_COUNT:
			AllCronJobsGauge.Set(parsedCount)
			break

		case BAD_NAMESPACE_COUNT:
			BadNamespaceCountGauge.Set(parsedCount)
			break
		case BAD_DEPLOYMENT_COUNT:
			BadDeploymentCountGauge.Set(parsedCount)
			break
		case BAD_DAEMONSET_COUNT:
			BadDaemonSetCountGauge.Set(parsedCount)
			break
		case BAD_POD_COUNT:
			BadPodCountGauge.Set(parsedCount)
			break
		case BAD_POD_WO_OWNER_COUNT:
			BadPodWoOwnerGauge.Set(parsedCount)
			break
		case BAD_INGRESSES_COUNT:
			BadIngressesCountGauge.Set(parsedCount)
			break
		case BAD_JOB_COUNT:
			BadJobCountGauge.Set(parsedCount)
			break
		case BAD_JOB_WO_OWNER_COUNT:
			BadJobWoOwnerGauge.Set(parsedCount)
			break
		case BAD_CRONJOB_COUNT:
			BadCronJobCountGauge.Set(parsedCount)
			break
		}
	}
	promhttp.Handler().ServeHTTP(w, r)
}
