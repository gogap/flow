Flow
====


```go
package main

import (
	"fmt"
	"time"

	"github.com/gogap/context"
	"github.com/gogap/flow"
)

func main() {

	h1 := func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H1", flowOpts.Config)

		return
	}

	h2 := func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H2", flowOpts.Config)

		return
	}

	flow.RegisterHandler("h1", h1)
	flow.RegisterHandler("h2", h2)

	flow.Begin().
		Then("h1", flow.ConfigString(`{a = 1}`)).
		Then("h2", flow.ConfigString(`{a = 2}`)).
		Subscribe(
			func(ctx context.Context, opts ...flow.Option) {
				fmt.Println("subscribed")
			}).Commit()

	// delay exist console
	time.Sleep(time.Second)
}
```

or

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogap/context"
	"github.com/gogap/flow"
)

func main() {

	myFlow := flow.New()

	h1 := func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H1", flowOpts.Config)

		return
	}

	h2 := func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H2", flowOpts.Config)

		return
	}

	myFlow.RegisterHandler("h1", h1)
	myFlow.RegisterHandler("h2", h2)

	myFlow.Begin().
		Then("h1", flow.ConfigString(`{a = 1}`)).
		Then("h2", flow.ConfigString(`{a = 2}`)).
		Subscribe(
			func(ctx context.Context, opts ...flow.Option) {
				fmt.Println("subscribed")
			}).Commit()

	// delay exist console
	time.Sleep(time.Second)
}

```

**output**

```bash
H1 {
  a : 1
}
H2 {
  a : 2
}
subscribed
```



```go
package main

import (
	"fmt"

	"github.com/gogap/context"
	"github.com/gogap/flow"
)

func main() {

	var h1 flow.HandlerFunc = func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H1", flowOpts.Config)

		return
	}

	var h2 flow.HandlerFunc = func(ctx context.Context, opts ...flow.Option) (err error) {

		flowOpts := flow.NewOptions(opts...)

		fmt.Println("H2", flowOpts.Config)

		return
	}

	h := h1.Then(h2, flow.ConfigString(`{config = h2}`))

	h(context.NewContext(), flow.ConfigString(`{config = 1}`))
}
```


**output**

```bash
H1 {
  config : 1
}
H2 {
  config : h2
}
```