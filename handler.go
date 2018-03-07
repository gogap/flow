package flow

import (
	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow/cache"
)

type ctxConfigKey struct {
	Name string
}

type ctxCacheKey struct {
	Name string
}

func ValueConfig(ctx context.Context, name string) config.Configuration {
	conf, _ := ctx.Value(ctxConfigKey{name}).(config.Configuration)
	return conf
}

func ValueCache(ctx context.Context) cache.Cache {
	c, _ := ctx.Value(ctxCacheKey{}).(cache.Cache)
	return c
}

func ConfigString(str string) config.Option {
	return config.ConfigString(str)
}

func ConfigFile(filename string) config.Option {
	return config.ConfigFile(filename)
}

type SubscriberFunc func(context.Context)

type HandlerFunc func(context.Context) error

func (p HandlerFunc) Run(ctx context.Context) error {
	return p(ctx)
}

func (p HandlerFunc) Then(next HandlerFunc) HandlerFunc {

	var h HandlerFunc = func(ctx context.Context) (err error) {

		err = p.Run(ctx)

		if err != nil {
			return
		}

		err = next.Run(ctx)

		if err != nil {
			return
		}

		return
	}

	return h
}

func (p HandlerFunc) Subscribe(subscribers ...SubscriberFunc) HandlerFunc {
	var h HandlerFunc = func(ctx context.Context) (err error) {
		err = p.Run(ctx)

		if len(subscribers) > 0 {
			var newCtx context.Context
			if ctx == nil {
				newCtx = context.NewContext()
			} else {
				newCtx = ctx.Copy()
			}

			if err != nil {
				newCtx.WithError(err)
			}

			p.publish(newCtx, subscribers...)
		}

		return
	}

	return h
}

func (p HandlerFunc) publish(ctx context.Context, subscribers ...SubscriberFunc) {

	if len(subscribers) == 0 {
		return
	}

	go func(ctx context.Context, subscribers ...SubscriberFunc) {

		for _, s := range subscribers {
			s(ctx)
		}
	}(ctx, subscribers...)
}
