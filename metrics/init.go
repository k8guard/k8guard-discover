package metrics

import (
	"strconv"
	"time"

	"github.com/k8guard/k8guard-discover/caching"
)

// to avoid conflicting metrics from sleeping workers. updates the metrics in memcached.
func Update(key string, value int) {
	expiration := time.Duration(METRIC_EXPIRE_SECONDS) * time.Second
	caching.Set(key, strconv.FormatInt(int64(value), 10), expiration)
}
