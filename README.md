> :warning: **Do not use in production**. `azul` is just a personal experiment at this time, probably full of bugs, and no support will be offered. :warning:

## `azul`: Deploy and run Azure Functions easily.

## Features
* Write Functions with Golang (not supported natively)
* Deploy/Delete Functions
* Logs

## Demo

![demo](demo2.gif)

## Installation
Get or update azul:
`go get -u github.com/wbuchwalter/azul/cmd/azul`

## Example

A `azul` project looks like this:

```
MyApp
├── azul.json
├── aFunc
│   └── main.go
└── anotherFunc
    └── index.js
```

`azul.json` will define on which Function App the `foo` and `bar` functions should be deployed.
This configuration file looks like this:

```json
{
  "name": "myfuncapp",
  "username": "$myfuncapp",
  "password": "1xXA2heWo7dD3mSmlvLhZnwzqJXMmrwHxogFCrnAnCn0idmo2vXCbiLKqqtY"
}
```


## Golang Support

`main.go` is where you define your actual function.
Here is an example of a function that returns the length of a word:
```go
package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul"
)

type input struct {
	Word string `json:"word"`
}

type Output struct {
	Length int `json:"length"`
}

func main() {
	azul.Handle(func(event json.RawMessage, logger logs.Logger) (interface{}, error) {
		var i input
		var output Output

		err := json.Unmarshal(event, &i)
		if err != nil {
			return nil, err
		}
		
		logger.Log("New request received, length: " + len(i.Word))
		output.Length = len(i.Word)

		return output, nil
	})
}

```

Deploy the function:  
`azul deploy wordLength`

## Limitations

* With Golang, `azul` only supports functions taking a request message as input (such as `httpTrigger`, `webHook`)
* Only `json` is supported as output for Golang Functions
* When calling `azul deploy`, `azul` will only returns the correct URL if `authLevel` is set to `anonymous` in the `function.json` (this is the case by default)


**fmt**

Do not use `fmt.Println` or similar functions writing to `stdin`.   
`stdin` is reserver by `azul` to communicate between the go function and the .NET runtime.
Instead use the `logger` passed as parameter to your function.


## FAQ

**Finding your Function App's credentials**  
In the Azure portal, open your `Function App`, click on `Function App Settings` -> `Advanced Settings`-> `Properties`-> `DEPLOYMENT TRIGGER URL`.  
The URL will look like this: `https://XXX:YYY@something.scm.azurewebsites.net/deploy`  
`XXX` is your username (you need to include the `$` if any), `YYY` is your password.

## Credits
Inspired by [TJ Holowaychuk](https://twitter.com/tjholowaychuk)'s [Apex](https://github.com/apex/apex).
