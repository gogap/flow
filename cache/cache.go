package cache

import (
	"fmt"

	"github.com/gogap/config"
)

type Option func(*Options)

type Options struct {
	Config config.Configuration
}

type NewCacheFunc func(...Option) (Cache, error)

func Config(conf config.Configuration) Option {
	return func(o *Options) {
		o.Config = conf
	}
}

type Cache interface {
	Set(k string, v interface{})
	Get(k string) (interface{}, bool)

	Increment(k string, delta int64) (int64, error)
	Decrement(k string, delta int64) (int64, error)

	Delete(k string)
	Flush()

	IsLocal() bool
	CanStoreInterface() bool
}

var (
	caches map[string]NewCacheFunc = make(map[string]NewCacheFunc)
)

func RegisterCache(name string, fn NewCacheFunc) {
	if len(name) == 0 {
		panic("cache name is empty")
	}

	if fn == nil {
		panic("cache fn is nil")
	}

	_, exist := caches[name]

	if exist {
		panic(fmt.Sprintf("cache already registered: %s", name))
	}

	caches[name] = fn
}

func NewCache(name string, opts ...Option) (c Cache, err error) {
	fn, exist := caches[name]

	if !exist {
		err = fmt.Errorf("cache not exist '%s'", name)
		return
	}

	c, err = fn(opts...)

	return
}
