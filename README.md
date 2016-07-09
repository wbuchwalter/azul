> :warning: **Do not use in production**. `lox` is just a personal experiment at this time, probably full of bugs, and no support will be offered. :warning:

## `lox`: Deploy and run Azure Functions written with Golang.

## Installation
Get or update `lox`:
`go get -u github.com/wbuchwalter/lox/cmd/lox`

## Example

A `lox` project looks like this:

```
MyApp
├── lox.json
├── foo
│   └── main.go
└── bar
    └── main.go
```

`lox.json` will define on which Function App the `foo` and `bar` functions will be deployed.
This configuration file looks like this:

```json
{
  "name": "myfuncapp",
  "username": "myfuncapp",
  "password": "1xXA2heWo7dD3mSmlvLhZnwzqJXMmrwHxogFCrnAnCn0idmo2vXCbiLKqqtY"
}
```

`main.go` is where you define your actual function.
Here is an example of a function that returns the length of a word:
```go
package main

import (
	"encoding/json"

	"github.com/wbuchwalter/lox"
)

type input struct {
	Word string `json:"word"`
}

type Output struct {
	Length int `json:"length"`
}

func main() {
	lox.Handle(func(event json.RawMessage) (interface{}, error) {
		var i input
		var output Output

		err := json.Unmarshal(event, &i)
		if err != nil || i.Name == "" {
			return nil, err
		}

		output.Length = len(i.Word)

		return output, nil
	}, "hello")
}
```

## Limitations

Currently, `lox` only supports one kind of function:

```
{
  "bindings": [
    {
      "authLevel": "function",
      "name": "req",
      "type": "httpTrigger",
      "direction": "in"
    },
    {
      "name": "res",
      "type": "http",
      "direction": "out"
    }
  ],
  "disabled": false
}
```

Custom configs are probably coming... eventually :smiley: .

## FAQ

**Finding your Function App's credentials**

**Performance**
Zero efforts have been put in optimization, and no benchmark has been made. So I have no idea about the performances of `lox`.

**Why name it `lox`?**
Why not?




## Credits
Inspired by [TJ Holowaychuk](https://twitter.com/tjholowaychuk)'s [Apex](https://github.com/apex/apex).
