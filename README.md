> :warning: **Do not use in production**. `lox` is just an experiment at this time, probably full of bugs, and no support will be offered. :warning:

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

## Finding your Function App's credentials

## Limitations

Currently, `lox` only supports one kind of functions:

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

Custom configs are coming eventually.

## Credits
Inspired by [TJ Holowaychuk](https://twitter.com/tjholowaychuk)'s [Apex](https://github.com/apex/apex).


### Next steps:
* Force nuget restore
* Return url of the created function
* clean run.csx
* clean code
* unit tests
* Documentation
* Publish
* expose commands (delete etc.)
* add commands (monitoring etc)
* allow multiple functions to be deployed at once
* Any way to reduce upload time of main.exe?
