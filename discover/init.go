package discover

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/k8s"
	"github.com/k8guard/k8guardlibs/messaging/kafka"
	"k8s.io/client-go/kubernetes"
)

// Clientset talks to kubernetes API
var Clientset *kubernetes.Clientset
var Memcached *memcache.Client
var KafkaProducer kafka.KafkaProducer
var err error

func init() {
	Clientset, err = k8s.LoadClientset()

	if err != nil {
		lib.Log.Error("error loading ClientSet ", err)
		panic(err)
	}

	Memcached = memcache.New(fmt.Sprintf("%s:11211", lib.Cfg.MemCachedHostname))

	KafkaProducer, err = kafka.NewProducer(kafka.DISCOVER_CLIENTID, lib.Cfg)
	if err != nil {
		lib.Log.Error("Creating Kafka Producer ", err)
		panic(err)
	}
	// defer KafkaProducer.Close()

}

func TestKafkaWithTestMessage() error {
	// Sending Test Data
	err := KafkaProducer.SendData(lib.Cfg.KafkaActionTopic, kafka.TEST_MESSAGE, "Testing")
	if err != nil {
		lib.Log.Error("Error trying to send test data to Kafka ", err)
		panic(err)
	}
	return err
}
