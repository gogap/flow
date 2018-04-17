package flow

import (
	"github.com/gogap/context"
)

type outputKey struct{}

type Output struct {
	Name  string
	Value interface{}
	Next  *Output
}

func AppendOutput(ctx context.Context, name string, value interface{}) {

	if ctx == nil {
		return
	}

	output, ok := ctx.Value(outputKey{}).(*Output)

	if !ok {
		ctx.WithValue(outputKey{}, &Output{Name: name, Value: value})
		return
	}

	if output == nil {
		output = &Output{Name: name, Value: value}
		ctx.WithValue(outputKey{}, output)
		return
	}

	for output != nil {
		if output.Next != nil {
			output = output.Next
			continue
		}

		output.Next = &Output{Name: name, Value: value}
		return
	}
}

func ListOutput(ctx context.Context) *Output {
	if ctx == nil {
		return nil
	}

	output, ok := ctx.Value(outputKey{}).(*Output)

	if !ok {
		return nil
	}

	return output
}
