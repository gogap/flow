package flow

import (
	"github.com/gogap/config"
	"github.com/gogap/context"
)

type SubscriberFunc func(context.Context)

type HandlerFunc func(context.Context, config.Configuration) error

func (p HandlerFunc) Run(ctx context.Context, conf config.Configuration) error {
	return p(ctx, conf)
}

func (p HandlerFunc) Then(next HandlerFunc, conf config.Configuration) HandlerFunc {

	nextConfig := conf

	var h HandlerFunc = func(ctx context.Context, conf config.Configuration) (err error) {

		err = p.Run(ctx, conf)

		if err != nil {
			return
		}

		err = next.Run(ctx, nextConfig)

		if err != nil {
			return
		}

		return
	}

	return h
}

func (p HandlerFunc) Subscribe(subscribers ...SubscriberFunc) HandlerFunc {
	var h HandlerFunc = func(ctx context.Context, conf config.Configuration) (err error) {
		err = p.Run(ctx, conf)

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
