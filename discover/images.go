package discover

import (
	"fmt"
	"strconv"
	"github.com/bradfitz/gomemcache/memcache"
	"k8s.io/client-go/pkg/api/v1"
	lib "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guard-discover/metrics"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

// gets all the images with their sizes and puts them in memcached
func cacheAllImages(storeInMemcached bool) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(metrics.FNCacheAllImagesGauge.Set))
	defer timer.ObserveDuration()

	imageCount := int64(0)

	nodes, err := Clientset.CoreV1().Nodes().List(v1.ListOptions{})
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


func isValidImageRepo(imageName string) bool {
	for _, i := range lib.Cfg.ApprovedImageRepos {
		if strings.Contains(imageName, i) {
			return true
		}
	}
	return false
}

func isValidImageSize(imageSizeByte int64) bool {
	// converting mb to bytes
	approvedImageSize := (lib.Cfg.ApprovedImageSize * 1024 * 1024)
	if imageSizeByte <= approvedImageSize {
		return true
	}
	return false
}
