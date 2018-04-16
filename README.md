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

type ctxKey struct{ Key string }

func main() {

	h1 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H1", conf)

		return
	}

	h2 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H2", conf)

		return
	}

	h3 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H3", conf)

		return
	}

	flow.RegisterHandler("h1", h1)
	flow.RegisterHandler("h2", h2)
	flow.RegisterHandler("h3", h3)

	ctx := context.NewContext()

	flow.Begin(ctx, config.ConfigString(`{config = default}`)).
		Then("h1", config.ConfigString(`{config = h1}`)).
		Then("h2", config.ConfigString(`{config = h2}`)).
		Then("h3").
		Subscribe(
			func(ctx context.Context) {
				fmt.Println("subscribed")
			}).Commit()

	// delay exist console
	time.Sleep(time.Second)
}


```

or you could define a Flow instance

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

	h1 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H1", conf)

		return
	}

	h2 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H2", conf)

		return
	}

	h3 := func(ctx context.Context, conf config.Configuration) (err error) {

		fmt.Println("H3", conf)

		return
	}

	myFlow.RegisterHandler("h1", h1)
	myFlow.RegisterHandler("h2", h2)
	myFlow.RegisterHandler("h3", h3)

	ctx := context.NewContext()

	myFlow.Begin(ctx, config.ConfigString(`{config = default}`)).
		Then("h1", config.ConfigString(`{config = h1}`)).
		Then("h2", config.ConfigString(`{config = h2}`)).
		Then("h3").
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
H3 {
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

	h1 := func(ctx context.Context, conf config.Configuration) (err error) {

		t := ctx.Value(ctxKey{"key1"}).(time.Time)

		fmt.Println("key1", t)

		return
	}

	flow.RegisterHandler("h1", h1)

	ctx := context.NewContext()

	ctx.WithValue(ctxKey{"key1"}, time.Now())

	flow.Begin(ctx).
		Then("h1", config.ConfigString(`{config = h1}`)).
		Commit()

	// delay exist console
	time.Sleep(time.Second)
}
```


**output**

```bash
key1 2018-04-16 20:51:03.220143 +0800 CST m=+0.000465875
```


#### create aliyun vpc

```go
package main

import (
	"fmt"

	"github.com/gogap/config"
	"github.com/gogap/context"
	"github.com/gogap/flow"

	_ "github.com/flow-contrib/aliyun"
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

	err = flow.Begin(ctx).
		Then("devops.aliyun.ecs.vpc.create", config.ConfigString(confStr)).
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

	_ "github.com/flow-contrib/goja"
)

var confStr = `
src = test.js
`

func main() {

	var err error

	defer func() { fmt.Println(err) }()

	err = flow.Begin(context.NewContext()).
		Then("lang.javascript.goja", config.ConfigString(confStr)).
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