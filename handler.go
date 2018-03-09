package flow

import (
	"github.com/gogap/context"
)

type SubscriberFunc func(context.Context)

type HandlerFunc func(context.Context, Params) error

func (p HandlerFunc) Run(ctx context.Context, params Params) error {
	return p(ctx, params)
}

func (p HandlerFunc) Then(next HandlerFunc, params Params) HandlerFunc {

	nextParams := params

	var h HandlerFunc = func(ctx context.Context, params Params) (err error) {

		err = p.Run(ctx, params)

		if err != nil {
			return
		}

		err = next.Run(ctx, nextParams)

		if err != nil {
			return
		}

		return
	}

	return h
}

func (p HandlerFunc) Subscribe(subscribers ...SubscriberFunc) HandlerFunc {
	var h HandlerFunc = func(ctx context.Context, params Params) (err error) {
		err = p.Run(ctx, params)

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
