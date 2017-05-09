package metrics

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/bradfitz/gomemcache/memcache"
)

const (
	ALL_DEPLOYMENT_COUNT string = "k8guard_all_deployment_count"
	ALL_POD_COUNT        string = "k8guard_all_pod_count"
	ALL_IMAGE_COUNT      string = "k8guard_all_image_count"
	ALL_INGRESSES_COUNT  string = "k8guard_all_ingresses_count"
	ALL_JOB_COUNT        string = "k8guard_all_job_count"
	ALL_CRONJOB_COUNT    string = "k8guard_all_cronjob_count"

	BAD_POD_COUNT          string = "k8guard_bad_pod_count"
	BAD_POD_WO_OWNER_COUNT string = "k8guard_bad_pod_wo_owner_count"
	BAD_DEPLOYMENT_COUNT   string = "k8guard_bad_deployment_count"
	BAD_INGRESSES_COUNT    string = "k8guard_bad_ingresses_count"
	BAD_JOB_COUNT          string = "k8guard_bad_job_count"
	BAD_JOB_WO_OWNER_COUNT string = "k8guard_bad_job_wo_owner_count"
	BAD_CRONJOB_COUNT      string = "k8guard_bad_cronjob_count"

	METRIC_EXPIRE_SECONDS int32 = 43200
)

var (
	AllDeploymentCountGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: ALL_DEPLOYMENT_COUNT,
			Help: "the number of all deployments",
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

func PromRegister() {
	prometheus.MustRegister(AllPodCountGauge)
	prometheus.MustRegister(AllImageCountGauge)
	prometheus.MustRegister(AllDeploymentCountGauge)
	prometheus.MustRegister(AllIngressesGauge)
	prometheus.MustRegister(AllJobsGauge)
	prometheus.MustRegister(AllCronJobsGauge)
	prometheus.MustRegister(BadDeploymentCountGauge)
	prometheus.MustRegister(BadPodCountGauge)
	prometheus.MustRegister(BadPodWoOwnerGauge)
	prometheus.MustRegister(BadIngressesCountGauge)
	prometheus.MustRegister(BadJobCountGauge)
	prometheus.MustRegister(BadJobWoOwnerGauge)
	prometheus.MustRegister(BadCronJobCountGauge)

	// Function Metrics
	prometheus.MustRegister(FNCacheAllImagesGauge)
	prometheus.MustRegister(FNGetBadDeploys)
	prometheus.MustRegister(FNGetBadIngresses)
	prometheus.MustRegister(FNGetBadPods)
	prometheus.MustRegister(FNGetBadJobs)
	prometheus.MustRegister(FNGetBadCronJobs)

}

func PromMetricsHandler(w http.ResponseWriter, r *http.Request) {
	myList := []string{
		ALL_DEPLOYMENT_COUNT, ALL_POD_COUNT, ALL_IMAGE_COUNT, ALL_INGRESSES_COUNT, ALL_JOB_COUNT,
		ALL_CRONJOB_COUNT, BAD_POD_COUNT, BAD_DEPLOYMENT_COUNT, BAD_POD_WO_OWNER_COUNT, BAD_INGRESSES_COUNT,
		BAD_JOB_COUNT, BAD_JOB_WO_OWNER_COUNT, BAD_CRONJOB_COUNT,
	}
	for _, i := range myList {
		count, err := Memcached.Get(i)
		if err != nil {
			count = &memcache.Item{Key: i, Value: []byte(strconv.Itoa(0))}
		}
		parsedCount, _ := strconv.ParseFloat(string(count.Value), 64)

		switch i {
		case ALL_POD_COUNT:
			AllPodCountGauge.Set(parsedCount)
			break
		case ALL_DEPLOYMENT_COUNT:
			AllDeploymentCountGauge.Set(parsedCount)
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
		case BAD_DEPLOYMENT_COUNT:
			BadDeploymentCountGauge.Set(parsedCount)
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

