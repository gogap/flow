Flow
====

#### usage demo

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"
)

func main() {

	h1 := func(ctx context.Context, params flow.Params) (err error) {

		fmt.Println("H1", params.Val("config").(config.Configuration))

		return
	}

	h2 := func(ctx context.Context, params flow.Params) (err error) {

		fmt.Println("H2", params.Val("config").(config.Configuration))

		return
	}

	h1Config := config.NewConfig(config.ConfigString(`{config = h1}`))
	h2Config := config.NewConfig(config.ConfigString(`{config = h2}`))
	defaultConfig := config.NewConfig(config.ConfigString(`{config = default}`))

	flow.RegisterHandler("h1", h1)
	flow.RegisterHandler("h2", h2)

	ctx := context.NewContext()

	flow.Begin(ctx, flow.Params{"config": defaultConfig}).
		Then("h1", flow.Params{"config": h1Config}).
		Then("h2", flow.Params{"config": h2Config}).
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

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"
)

func main() {

	myFlow := flow.New()

	h1 := func(ctx context.Context, params flow.Params) (err error) {

		fmt.Println("H1", params.Val("config").(config.Configuration))

		return
	}

	h2 := func(ctx context.Context, params flow.Params) (err error) {

		fmt.Println("H2", params.Val("config").(config.Configuration))

		return
	}

	h1Config := config.NewConfig(config.ConfigString(`{config = h1}`))
	h2Config := config.NewConfig(config.ConfigString(`{config = h2}`))
	defaultConfig := config.NewConfig(config.ConfigString(`{config = default}`))

	myFlow.RegisterHandler("h1", h1)
	myFlow.RegisterHandler("h2", h2)

	ctx := context.NewContext()

	myFlow.Begin(ctx, flow.Params{"config": defaultConfig}).
		Then("h1", flow.Params{"config": h1Config}).
		Then("h2", flow.Params{"config": h2Config}).
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
  config : h1
}
H2 {
  config : h2
}
H2 {
  config : default
}
subscribed
```

#### handler context demo

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"
)

type ctxKey struct{ Key string }

func main() {

	myFlow := flow.New()

	h1 := func(ctx context.Context, params flow.Params) (err error) {

		v := ctx.Value(ctxKey{"H1"}).(config.Configuration)

		fmt.Println("H1", v)

		return
	}

	h2 := func(ctx context.Context, params flow.Params) (err error) {

		v := ctx.Value(ctxKey{"H2"}).(config.Configuration)

		fmt.Println("H2", v)

		return
	}

	h1Config := config.NewConfig(config.ConfigString(`{config = h1}`))
	h2Config := config.NewConfig(config.ConfigString(`{config = h2}`))

	myFlow.RegisterHandler("h1", h1)
	myFlow.RegisterHandler("h2", h2)

	ctx := context.NewContext()

	ctx.WithValue(ctxKey{"H1"}, h1Config)
	ctx.WithValue(ctxKey{"H2"}, h2Config)

	myFlow.Begin(ctx).
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
  config : h1
}
H2 {
  config : h2
}
subscribed
```


#### create aliyun vpc

```go
package main

import (
	"fmt"

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"

	_ "github.com/gogap/flow-contrib/handler/devops/aliyun"
)

var confStr = `
aliyun {
	region = cn-beijing
	access-key-id = ${DEVOPS_ALIYUN_ACCESS_KEY_ID}
	access-key-secret = ${DEVOPS_ALIYUN_ACCESS_KEY_SECRET}

	ecs {
		vpc  {
			test {
				cidr-block  = "172.16.0.0/16"
				description = "172.16.0.0/16"
			}
		}
	}
}
`

func main() {

	var err error

	defer func() { fmt.Println(err) }()

	ctx := context.NewContext()

	ctx.WithValue("CODE", "test")

	conf := config.NewConfig(config.ConfigString(confStr))

	err = flow.Begin(ctx).
		Then("devops.aliyun.ecs.vpc.create", flow.Params{"aliyun.config": conf}).
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

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"

	_ "github.com/gogap/flow-contrib/handler/lang/javascript/goja"
)

var confStr = `
src = test.js
`

func main() {

	var err error

	conf := config.NewConfig(config.ConfigString(confStr))

	defer func() { fmt.Println(err) }()

	err = flow.Begin(context.NewContext()).
		Then("lang.javascript.goja", flow.Params{"goja.config": conf}).
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