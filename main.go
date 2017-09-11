package main

import (
	"flag"
	"sync"

	"github.com/k8guard/k8guard-discover/caching"
	"github.com/k8guard/k8guard-discover/discover"
	"github.com/k8guard/k8guard-discover/messaging"
	"github.com/k8guard/k8guard-discover/metrics"
	lib "github.com/k8guard/k8guardlibs"
)

var (
	Version string
	Build   string
)

var err error

func init() {
	metrics.PromRegister()
	caching.InitCache()
}

func main() {

	messagingMode := flag.Bool("kmode", false, "messaging mode, no router")
	flag.Parse()

	if *messagingMode {
		defer messaging.CloseBroker()
		messaging.InitBroker()

		// test if broker is there before making api calls
		messaging.TestBrokerWithTestMessage()

		var waitGroup sync.WaitGroup
		waitGroup.Add(7)
		lib.Log.Infof("Starting in message mode using %s as broker", lib.Cfg.MessageBroker)
		lib.Log.Info("Version: ", Version)
		lib.Log.Info("BuildNumber: ", Build)

		if err != nil {
			panic(err)
		}
		go func() {
			defer waitGroup.Done()
			discover.GetBadNamespaces(discover.GetAllNamespaceFromApi(), true)
		}()
		go func() {
			defer waitGroup.Done()
			discover.GetBadDeploys(discover.GetAllDeployFromApi(), true)
		}()

		go func() {
			defer waitGroup.Done()
			discover.GetBadDaemonSets(discover.GetAllDaemonSetFromApi(), true)
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

		go func() {
			messaging.InitStatsHandler()
		}()

		startHttpServer()
	}

}
