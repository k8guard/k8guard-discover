package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/k8guard/k8guard-discover/discover"
	"github.com/k8guard/k8guard-discover/metrics"
	lib "github.com/k8guard/k8guardlibs"
)

var (
	Version string
	Build   string
)

var Memcached *memcache.Client

func init() {
	Memcached = memcache.New(fmt.Sprintf("%s:11211", lib.Cfg.MemCachedHostname))

}

var err error

func init() {
	metrics.PromRegister()

}

func main() {
	defer discover.KafkaProducer.Close()

	kafkaMode := flag.Bool("kmode", false, "messaging mode, no router")
	flag.Parse()

	if *kafkaMode {
		// test if kafka is there before making api calls
		discover.TestKafkaWithTestMessage()

		var waitGroup sync.WaitGroup
		waitGroup.Add(5)
		lib.Log.Info("Starting in Kafka Mode")

		if err != nil {
			panic(err)
		}
		go func() {
			defer waitGroup.Done()
			discover.GetBadDeploys(discover.GetAllDeployFromApi(), true)
		}()

		go func() {
			defer waitGroup.Done()
			discover.GetBadIngresses(discover.GetAllIngressFromApi(), true)
		}()

		go func() {
			defer waitGroup.Done()
			discover.GetBadPods(discover.GetAllPodsFromApi(), true)
		}()

		go func() {
			defer waitGroup.Done()
			discover.GetBadJobs(discover.GetAllJobFromApi(), true)
		}()

		go func() {
			defer waitGroup.Done()
			discover.GetBadCronJobs(discover.GetAllCronJobFromApi(), true)
		}()

		waitGroup.Wait()

	} else {
		startHttpServer()

	}

}
