package flow

import (
	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow/cache"
)

type Options struct {
	Config config.Configuration
	Cache  cache.Cache
}

func ParseOptions(opts ...Option) *Options {
	o := &Options{}
	o.Init(opts...)
	return o
}

func (p *Options) Init(opts ...Option) *Options {
	for _, o := range opts {
		o(p)
	}
	return p
}

type Option func(*Options)

func ConfigString(str string) Option {
	return func(o *Options) {
		conf := config.NewConfig(config.ConfigString(str))
		o.Config = conf
	}
}

func ConfigFile(filename string) Option {
	return func(o *Options) {
		conf := config.NewConfig(config.ConfigFile(filename))
		o.Config = conf
	}
}

func Cache(cache cache.Cache) Option {
	return func(o *Options) {
		o.Cache = cache
	}
}

type ctxOptionsKey struct {
}

type SubscriberFunc func(context.Context, ...Option)

type HandlerFunc func(context.Context, ...Option) error

func (p HandlerFunc) Run(ctx context.Context, opts ...Option) error {

	ctx.WithValue(ctxOptionsKey{}, opts)

	return p(ctx, opts...)
}

func (p HandlerFunc) Then(next HandlerFunc, opts ...Option) HandlerFunc {

	nextOpts := opts

	var h HandlerFunc = func(ctx context.Context, opts ...Option) (err error) {

		err = p.Run(ctx, opts...)

		if err != nil {
			return
		}

		if len(nextOpts) == 0 {
			err = next.Run(ctx, opts...)
		} else {
			err = next.Run(ctx, nextOpts...)
		}

		if err != nil {
			return
		}

		return
	}

	return h
}

func (p HandlerFunc) Subscribe(subscribers ...SubscriberFunc) HandlerFunc {
	var h HandlerFunc = func(ctx context.Context, opts ...Option) (err error) {
		err = p.Run(ctx, opts...)

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

			p.publish(newCtx, opts, subscribers...)
		}

		return
	}

	return h
}

func (p HandlerFunc) publish(ctx context.Context, opts []Option, subscribers ...SubscriberFunc) {

	if len(subscribers) == 0 {
		return
	}

	go func(ctx context.Context, opts []Option, subscribers ...SubscriberFunc) {

		for _, s := range subscribers {
			s(ctx, opts...)
		}
	}(ctx, opts, subscribers...)
}
