package metrics

import (
	"fmt"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	lib "github.com/k8guard/k8guardlibs"
)

var Memcached *memcache.Client

func init() {
	Memcached = memcache.New(fmt.Sprintf("%s:11211", lib.Cfg.MemCachedHostname))

}

// to avoid conflicting metrics from sleeping workers. updates the metrics in memcached.
func Update(key string, value int) {
	Memcached.Set(&memcache.Item{Key: key, Expiration: METRIC_EXPIRE_SECONDS, Value: []byte(strconv.FormatInt(int64(value), 10))})
}
