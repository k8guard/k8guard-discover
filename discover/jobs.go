package discover

import (
	"k8s.io/client-go/pkg/apis/batch/v2alpha1"
	lib "github.com/k8guard/k8guardlibs"
	"k8s.io/client-go/pkg/api/v1"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
	"strings"
	"k8guard-discover/metrics"
	"github.com/prometheus/client_golang/prometheus"
)


func GetAllJobFromApi() []v2alpha1.Job {
	jobs, err := Clientset.BatchV2alpha1().Jobs(lib.Cfg.Namespace).List(v1.ListOptions{})

	if err != nil {
		lib.Log.Error("error:", err)
		panic(err.Error())
	}
	metrics.Update(metrics.ALL_JOB_COUNT, len(jobs.Items))

	return jobs.Items
}

func GetAllCronJobFromApi() []v2alpha1.CronJob {
	cronjobs, err := Clientset.BatchV2alpha1().CronJobs(lib.Cfg.Namespace).List(v1.ListOptions{})

	if err != nil {
		lib.Log.Error("error:", err)
		panic(err.Error())
	}
	metrics.Update(metrics.ALL_CRONJOB_COUNT, len(cronjobs.Items))

	return cronjobs.Items
}

func GetBadCronJobs(allCronJobs []v2alpha1.CronJob, sendToKafka bool) []lib.CronJob {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadCronJobs.Set))
	defer timer.ObserveDuration()

	allBadCronJobs := []lib.CronJob{}

	cacheAllImages(true)

	for _, kcj := range allCronJobs {
		if isIgnoredNamespace(kcj.Namespace) == true || isIgnoredCronJob(kcj.ObjectMeta.Name) == true {
			continue
		}
		if *kcj.Spec.Suspend {
			continue
		}
		cj := lib.CronJob{}
		cj.Name = kcj.Name
		cj.Cluster = lib.Cfg.ClusterName
		cj.Namespace = kcj.Namespace
		getVolumesWithHostPathForAPod(kcj.Spec.JobTemplate.Spec.Template.Spec, &cj.ViolatableEntity)
		GetBadContainers(kcj.Spec.JobTemplate.Spec.Template.Spec, &cj.ViolatableEntity)

		if len(cj.Violations) > 0 {
			allBadCronJobs = append(allBadCronJobs, cj)
			if sendToKafka {
				lib.Log.Debug("Sending ", cj.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.CRONJOB_MESSAGE, cj)
				if err != nil {
					panic(err)
				}
			}

		}

	}
	metrics.Update(metrics.BAD_CRONJOB_COUNT, len(allBadCronJobs))
	return allBadCronJobs

}

func GetBadJobs(allJobs []v2alpha1.Job, sendToKafka bool) []lib.Job {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNGetBadJobs.Set))
	defer timer.ObserveDuration()

	allBadJobsWitoutOwner := []lib.Job{}
	badJobsCounter := int64(0)

	cacheAllImages(true)

	for _, kj := range allJobs {

		_, createdByAnnotation := kj.Annotations["kubernetes.io/created-by"]
		if createdByAnnotation == true {
			continue
		}

		if isIgnoredNamespace(kj.Namespace) == true || isIgnoredJob(kj.ObjectMeta.Name) == true {
			continue
		}
		if kj.Status.Active == 0 {
			continue
		}
		j := lib.Job{}
		j.Name = kj.Name
		j.Cluster = lib.Cfg.ClusterName
		j.Namespace = kj.Namespace
		getVolumesWithHostPathForAPod(kj.Spec.Template.Spec, &j.ViolatableEntity)
		GetBadContainers(kj.Spec.Template.Spec, &j.ViolatableEntity)

		if (len(j.Violations) > 0) {

			badJobsCounter += 1
			allBadJobsWitoutOwner = append(allBadJobsWitoutOwner, j)
			if sendToKafka {
				lib.Log.Debug("Sending ", j.Name, " to kafka")
				err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.JOB_MESSAGE, j)
				if err != nil {
					panic(err)
				}
			}

		}

	}

	metrics.Update(metrics.BAD_JOB_COUNT, int(badJobsCounter))
	metrics.Update(metrics.BAD_JOB_WO_OWNER_COUNT, len(allBadJobsWitoutOwner))
	return allBadJobsWitoutOwner

}

func isIgnoredJob(jobName string) bool {
	for _, d := range lib.Cfg.IgnoredJobs {
		if strings.HasPrefix(jobName, d) == true {
			return true
		}
	}
	return false
}

func isIgnoredCronJob(cronJobName string) bool {
	for _, d := range lib.Cfg.IgnoredCronJobs {
		if strings.HasPrefix(cronJobName, d) == true {
			return true
		}
	}
	return false
}
