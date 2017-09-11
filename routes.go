package main

import (
	"net/http"
	"time"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"

	"text/template"

	"github.com/k8guard/k8guard-discover/caching"
	"github.com/k8guard/k8guard-discover/discover"
	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/k8guard/k8guard-discover/templates"
	lib "github.com/k8guard/k8guardlibs"
)

func startHttpServer() {

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// Only 5 request will be processed at a time.
	// r.Use(middleware.Throttle(5))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// When a client closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	r.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(500 * time.Second))

	//r.Mount("/metrics", promMetricsHandler())

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.PromMetricsHandler(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type myRoute struct {
			Names   []string
			Cluster string
			Version string
		}

		data := myRoute{Cluster: lib.Cfg.ClusterName, Names: []string{"namespaces", "deploys", "daemonsets", "podswo", "ingresses", "jobs", "cronjobs", "config", "metrics", "version"}, Version: Version + "_" + Build}

		t := template.New("Index Template")
		t, err := t.Parse(templates.INDEX_TEMPLATE_DISCOVER)
		if err != nil {
			lib.Log.Fatalln(err)
		}

		err = t.Execute(w, &data)
		if err != nil {
			lib.Log.Fatalln(err)
		}
	})

	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, lib.Cfg)
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		version_build := map[string]string{
			"version": Version,
			"build":   Build,
		}
		render.JSON(w, r, version_build)
	})

	r.Get("/deploys", func(w http.ResponseWriter, r *http.Request) {
		cachedBadDeploys, err := caching.GetAsJson("cached_bad_deploys")
		if err != nil {
			caching.SetAsJson("cached_bad_deploys", discover.GetBadDeploys(discover.GetAllDeployFromApi(), false), 300*time.Second)
			cachedBadDeploys, err = caching.GetAsJson("cached_bad_deploys")
			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadDeploys)
		} else {
			renderJSONString(w, r, cachedBadDeploys)
		}

	})

	r.Get("/daemonsets", func(w http.ResponseWriter, r *http.Request) {
		cachedBadDaemonSets, err := caching.GetAsJson("cached_bad_daemonsets")
		if err != nil {
			caching.SetAsJson("cached_bad_daemonsets", discover.GetBadDaemonSets(discover.GetAllDaemonSetFromApi(), false), 300*time.Second)
			cachedBadDaemonSets, err = caching.GetAsJson("cached_bad_daemonsets")
			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadDaemonSets)
		} else {
			renderJSONString(w, r, cachedBadDaemonSets)
		}

	})

	// Pods without owner
	r.Get("/podswo", func(w http.ResponseWriter, r *http.Request) {
		cachedBadPods, err := caching.GetAsJson("cached_bad_pods_wo")
		if err != nil {
			caching.SetAsJson("cached_bad_pods_wo", discover.GetBadPods(discover.GetAllPodsFromApi(), false), 300*time.Second)
			cachedBadPods, err = caching.GetAsJson("cached_bad_pods_wo")

			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadPods)
		} else {
			renderJSONString(w, r, cachedBadPods)
		}

	})

	r.Get("/ingresses", func(w http.ResponseWriter, r *http.Request) {
		cachedBadIngresses, err := caching.GetAsJson("cached_bad_ingresses")
		if err != nil {
			caching.SetAsJson("cached_bad_ingresses", discover.GetBadIngresses(discover.GetAllIngressFromApi(), false), 300*time.Second)
			cachedBadIngresses, err = caching.GetAsJson("cached_bad_ingresses")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadIngresses)
		} else {
			renderJSONString(w, r, cachedBadIngresses)
		}

	})

	r.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		cachedBadJobs, err := caching.GetAsJson("cached_bad_jobs")
		if err != nil {
			caching.SetAsJson("cached_bad_jobs", discover.GetBadJobs(discover.GetAllJobFromApi(), false), 300*time.Second)
			cachedBadJobs, err = caching.GetAsJson("cached_bad_jobs")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadJobs)
		} else {
			renderJSONString(w, r, cachedBadJobs)
		}

	})

	r.Get("/cronjobs", func(w http.ResponseWriter, r *http.Request) {
		cachedBadCronJobs, err := caching.GetAsJson("cached_bad_cronjobs")
		if err != nil {
			caching.SetAsJson("cached_bad_cronjobs", discover.GetBadCronJobs(discover.GetAllCronJobFromApi(), false), 300*time.Second)
			cachedBadCronJobs, err = caching.GetAsJson("cached_bad_cronjobs")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadCronJobs)
		} else {
			renderJSONString(w, r, cachedBadCronJobs)
		}

	})

	r.Get("/namespaces", func(w http.ResponseWriter, r *http.Request) {
		cachedBadNamespaces, err := caching.GetAsJson("cached_bad_namespaces")
		if cachedBadNamespaces == nil || err != nil {
			caching.SetAsJson("cached_bad_namespaces", discover.GetBadNamespaces(discover.GetAllNamespaceFromApi(), false), 300*time.Second)
			cachedBadNamespaces, err = caching.GetAsJson("cached_bad_namespaces")
			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadNamespaces)
		} else {
			renderJSONString(w, r, cachedBadNamespaces)
		}

	})

	http.ListenAndServe(":3000", r)
}

// func renderJSONString(w http.ResponseWriter, r *http.Request, b []byte) {
func renderJSONString(w http.ResponseWriter, r *http.Request, obj interface{}) {
	// var objmap []*json.RawMessage
	// err := json.Unmarshal(b, &objmap)
	// if err != nil {
	// 	panic(err)
	// }
	render.JSON(w, r, obj)
}
