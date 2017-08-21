package discover

import (
	"fmt"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/k8guard/k8guard-discover/metrics"
	"github.com/k8guard/k8guard-discover/rules"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// gets all the images with their sizes and puts them in memcached
func cacheAllImages(storeInMemcached bool) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNCacheAllImagesGauge.Set))
	defer timer.ObserveDuration()

	imageCount := int64(0)

	nodes, err := Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	nodesList := nodes.Items
	sem := make(chan int64, len(nodesList)) // semaphore pattern

	for _, node := range nodesList {
		go func(node v1.Node) {
			nodeImageCount := int64(0)
			for _, kImage := range node.Status.Images {
				for _, name := range kImage.Names {
					nodeImageCount += 1
					if storeInMemcached == true {
						Memcached.Set(&memcache.Item{Key: fmt.Sprintf("image_%s", name), Expiration: 600, Value: []byte(strconv.FormatInt(kImage.SizeBytes, 10))})
					}

				}
			}
			sem <- nodeImageCount
		}(node)
	}

	// wait for goroutines to finish
	for i := 0; i < len(nodesList); i++ {

		nodeImageCount := <-sem
		imageCount += nodeImageCount

	}

	metrics.Update(metrics.ALL_IMAGE_COUNT, int(imageCount))

}

func isValidImageRepo(namespace string, entityType string, entityName string, imageName string) bool {
	match, _ := rules.IsValueMatchContainsRule(namespace, entityType, entityName, imageName, lib.Cfg.ApprovedImageRepos)
	return match
}

func isValidImageSize(imageSizeByte int64) bool {
	// converting mb to bytes
	approvedImageSize := (lib.Cfg.ApprovedImageSize * 1024 * 1024)
	if imageSizeByte <= approvedImageSize {
		return true
	}
	return false
}
