


### Initail flow config

```go
flow1ConfigStr := `
	context {
		provider = LocalContextProvider
		options = {}
	}

	runner {
	
		type = PipeTaskRunner
		
		options = {
			singleton = false

			anko {
				singleton = true
				dir = "/anko_script"
			}

			goja {
				singleton = true
				dir = "/goja_script"
				timelimit = 180s
			}

			otto {
				dir = "/otto_script"
				timelimit = 5s
			}
		}
	}

	steps {
		order = [A,B,C]
		A.handler = anko
		B.handler = goja
		C.handler = otto
	}
`
```


#### Create flow

```go
var f *flow.Flow
	f, err = flow.NewFlow("flow1",
		flow.ConfigString(flow1ConfigStr),
	)
```

#### Example Script

`/anko_script/flow1/A.ank`

```go
println("print from ank")
val,exist=ctx.Get("hello")
println("context value",val,ctx.ID())
```

`/goja_script/flow1/B.js`

```go
console.log("B.js")
```



`/otto_script/flow1/C.js`

```go
var lyrics = [
  {line: 1, words: "I'm a lumberjack and I'm okay"},
  {line: 2, words: "I sleep all night and I work all day"},
  {line: 3, words: "He's a lumberjack and he's okay"},
  {line: 4, words: "He sleeps all night and he works all day"}
];

v=_.chain(lyrics)
  .map(function(line) { return line.words.split(' '); })
  .flatten()
  .reduce(function(counts, word) {
    counts[word] = (counts[word] || 0) + 1;
    return counts;
  }, {})
  .value();


console.log( JSON.stringify(v))
```


#### Run flow

```go
task := f.NewTask()

task.Context().Set("hello", "world")

task.Run()
```

Output:

```bash
print from ank
context value hello c13e3850-d1e7-4cc1-bd04-90f9db791dba

2017/09/09 17:00:31 B.js

{"He":1,"He's":1,"I":2,"I'm":2,"a":2,"all":4,"and":4,"day":2,"he":1,"he's":1,"lumberjack":2,"night":2,"okay":2,"sleep":1,"sleeps":1,"work":1,"works":1}
```
