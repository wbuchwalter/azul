package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul/azul-go"
	"github.com/wbuchwalter/azul/azul-go/logs"
)

func main() {
	azul.Handle(func(raw json.RawMessage, logger logs.Logger) (interface{}, error) {
		return "Hello World!", nil
	})
}
