Flow
====

#### usage demo

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogap/context"
	"github.com/gogap/flow"
)

type ctxKey struct{ Key string }

func main() {

	h1 := func(ctx context.Context) (err error) {

		v := flow.ValueConfig(ctx, "h1")

		fmt.Println("H1", v)

		return
	}

	h2 := func(ctx context.Context) (err error) {

		v := flow.ValueConfig(ctx, "h2")

		fmt.Println("H2", v)

		return
	}

	flow.RegisterHandler("h1", h1)
	flow.RegisterHandler("h2", h2)

	ctx := context.NewContext()

	flow.Begin(ctx).
		WithConfig("h1", flow.ConfigString(`{config = h2}`)).
		WithConfig("h2", flow.ConfigString(`{config = h2}`)).
		Then("h1").
		Then("h2").
		Subscribe(
			func(ctx context.Context) {
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

type ctxKey struct{ Key string }

func main() {

	myFlow := flow.New()

	h1 := func(ctx context.Context) (err error) {

		v := flow.ValueConfig(ctx, "h1")

		fmt.Println("H1", v)

		return
	}

	h2 := func(ctx context.Context) (err error) {

		v := flow.ValueConfig(ctx, "h2")

		fmt.Println("H2", v)

		return
	}

	myFlow.RegisterHandler("h1", h1)
	myFlow.RegisterHandler("h2", h2)

	ctx := context.NewContext()

	myFlow.Begin(ctx).
		WithConfig("h1", flow.ConfigString(`{config = h2}`)).
		WithConfig("h2", flow.ConfigString(`{config = h2}`)).
		Then("h1").
		Then("h2").
		Subscribe(
			func(ctx context.Context) {
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

#### handler options demo

```go
package main

import (
	"fmt"

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"
)

type ctxKey struct{ Key string }

func main() {

	var h1 flow.HandlerFunc = func(ctx context.Context) (err error) {

		v := ctx.Value(ctxKey{"H1"}).(config.Configuration)

		fmt.Println("H1", v)

		return
	}

	var h2 flow.HandlerFunc = func(ctx context.Context) (err error) {

		v := ctx.Value(ctxKey{"H2"}).(config.Configuration)

		fmt.Println("H2", v)

		return
	}

	ctx := context.NewContext()
	ctx.WithValue(ctxKey{"H1"}, config.NewConfig(flow.ConfigString(`{config = h1}`)))
	ctx.WithValue(ctxKey{"H2"}, config.NewConfig(flow.ConfigString(`{config = h2}`)))

	h := h1.Then(h2)

	h(ctx)
}

```


**output**

```bash
H1 {
  config : h1
}
H2 {
  config : h2
}
```


#### create aliyun vpc

```go
package main

import (
	"fmt"

	"github.com/gogap/context"
	"github.com/gogap/flow"

	_ "github.com/gogap/flow-contrib/handler/devops/aliyun"
	_ "github.com/gogap/flow/cache/redis"
)

var confStr = `
aliyun {
	region = cn-beijing
	access-key-id = ""
	access-key-secret = ""

	vpc  {
		test {
			cidr-block  = "172.16.0.0/16"
			description = "172.16.0.0/16"
		}
	}
}
`

func main() {

	var err error

	defer func() { fmt.Println(err) }()

	ctx := context.NewContext()

	ctx.WithValue("CODE", "test")

	err = flow.Begin(ctx).
		WithCache("go-redis").
		WithConfig("devops.aliyun", flow.ConfigString(confStr)).
		Then("devops.aliyun.vpc.create").
		Commit()

	if err != nil {
		return
	}
}

```

#### execute js

```go
package main

import (
	"fmt"

	"github.com/gogap/context"
	"github.com/gogap/flow"

	_ "github.com/gogap/flow-contrib/handler/lang/javascript/goja"
)

var confStr = `
src = test.js
`

func main() {

	var err error

	defer func() { fmt.Println(err) }()

	err = flow.Begin(context.NewContext()).
		WithConfig("lang.javascript.goja", flow.ConfigString(confStr)).
		Then("lang.javascript.goja").
		Commit()

	if err != nil {
		return
	}
}
```

`test.js`

```javascript
console.log("I am from goja")
```