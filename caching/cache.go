package caching

import (
	"encoding/json"
	"strconv"
	"time"

	lib "github.com/k8guard/k8guardlibs"
	c "github.com/k8guard/k8guardlibs/caching"
	"github.com/k8guard/k8guardlibs/caching/types"
)

var CacheClient types.Cache

func InitCache() {
	s, err := c.CreateCache(
		types.CacheType(lib.Cfg.CacheType), lib.Cfg)
	if err != nil {
		lib.Log.Error("Creating cache client ", err)
		panic(err)
	}
	CacheClient = s
}

func Set(key string, value interface{}, expiration time.Duration) {
	lib.Log.Debugf("Setting %s: %v", key, value)
	err := CacheClient.Set(key, value, expiration)
	if err != nil {
		panic(err)
	}
}

func SetAsJson(key string, value interface{}, expiration time.Duration) {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	Set(key, v, expiration)
}

func Get(key string) (interface{}, error) {
	lib.Log.Debugf("Getting %s", key)
	return CacheClient.Get(key)
}

func GetAsJson(key string) (interface{}, error) {
	lib.Log.Debugf("Getting %s as json", key)
	c, err := Get(key)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	var objmap []*json.RawMessage
	if v, ok := c.([]byte); ok {
		err = json.Unmarshal(v, &objmap)
	} else {
		err = json.Unmarshal([]byte(c.(string)), &objmap)
	}
	if err != nil {
		panic(err)
	}
	return objmap, nil
}

func GetAsInt(key string) (int64, error) {
	val, err := Get(key)
	if err == nil {
		if v, ok := val.([]byte); ok {
			return strconv.ParseInt(string(v), 10, 64)
		}
		myVal := val.(string)
		return strconv.ParseInt(myVal, 10, 64)
	}
	return 0, err
}
