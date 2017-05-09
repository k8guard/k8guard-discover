package main



import (
	"net/http"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"time"

	lib "github.com/k8guard/k8guardlibs"
	"encoding/json"
	"text/template"
	"k8guard-discover/templates"
	"k8guard-discover/metrics"
	"k8guard-discover/discover"
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
			Names []string
			Cluster string
			Version string
		}

		data := myRoute{Cluster: lib.Cfg.ClusterName, Names: []string{"deploys", "podswo", "ingresses", "jobs", "cronjobs", "config", "metrics", "version"},Version:Version+"_"+Build}

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
		cachedBadDeploys, err := Memcached.Get("cached_bad_deploys")
		if err != nil {
			response, _ := json.Marshal(discover.GetBadDeploys(discover.GetAllDeployFromApi(), false))
			Memcached.Set(&memcache.Item{Key: "cached_bad_deploys", Expiration: 300, Value: []byte(response)})
			cachedBadDeploys, err = Memcached.Get("cached_bad_deploys")
			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadDeploys.Value)
		} else {
			renderJSONString(w, r, cachedBadDeploys.Value)
		}

	})

	// Pods without owner
	r.Get("/podswo", func(w http.ResponseWriter, r *http.Request) {
		cachedBadPods, err := Memcached.Get("cached_bad_pods_wo")
		if err != nil {
			response, _ := json.Marshal(discover.GetBadPods(discover.GetAllPodsFromApi(), false))
			Memcached.Set(&memcache.Item{Key: "cached_bad_pods_wo", Expiration: 300, Value: []byte(response)})
			cachedBadPods, err = Memcached.Get("cached_bad_pods_wo")

			if err != nil {
				panic(err)
			}

			renderJSONString(w, r, cachedBadPods.Value)
		} else {
			renderJSONString(w, r, cachedBadPods.Value)
		}

	})

	r.Get("/ingresses", func(w http.ResponseWriter, r *http.Request) {
		cachedBadIngresses, err := Memcached.Get("cached_bad_ingresses")
		if err != nil {
			response, _ := json.Marshal(discover.GetBadIngresses(discover.GetAllIngressFromApi(), false))
			Memcached.Set(&memcache.Item{Key: "cached_bad_ingresses", Expiration: 300, Value: []byte(response)})
			cachedBadIngresses, err = Memcached.Get("cached_bad_ingresses")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadIngresses.Value)
		} else {
			renderJSONString(w, r, cachedBadIngresses.Value)
		}

	})

	r.Get("/jobs", func(w http.ResponseWriter, r *http.Request) {
		cachedBadJobs, err := Memcached.Get("cached_bad_jobs")
		if err != nil {
			response, _ := json.Marshal(discover.GetBadJobs(discover.GetAllJobFromApi(), false))
			Memcached.Set(&memcache.Item{Key: "cached_bad_jobs", Expiration: 300, Value: []byte(response)})
			cachedBadJobs, err = Memcached.Get("cached_bad_jobs")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadJobs.Value)
		} else {
			renderJSONString(w, r, cachedBadJobs.Value)
		}

	})

	r.Get("/cronjobs", func(w http.ResponseWriter, r *http.Request) {
		cachedBadCronJobs, err := Memcached.Get("cached_bad_cronjobs")
		if err != nil {
			response, _ := json.Marshal(discover.GetBadCronJobs(discover.GetAllCronJobFromApi(), false))
			Memcached.Set(&memcache.Item{Key: "cached_bad_cronjobs", Expiration: 300, Value: []byte(response)})
			cachedBadCronJobs, err = Memcached.Get("cached_bad_cronjobs")
			if err != nil {
				panic(err)
			}
			renderJSONString(w, r, cachedBadCronJobs.Value)
		} else {
			renderJSONString(w, r, cachedBadCronJobs.Value)
		}

	})

	http.ListenAndServe(":3000", r)
}
func renderJSONString(w http.ResponseWriter, r *http.Request, b []byte) {
	var objmap []*json.RawMessage
	err := json.Unmarshal(b, &objmap)
	if err != nil {
		panic(err)
	}
	render.JSON(w, r, objmap)
}