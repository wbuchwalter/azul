package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul/azul-go"
	"github.com/wbuchwalter/azul/azul-go/logs"
)

type event struct {
	Name string `json:"name"`
}

type response struct {
	Body string `json:"body"`
}

func main() {
	//	azul.FunctionName = "lol"
	azul.Handle(func(raw json.RawMessage, logger logs.Logger) (interface{}, error) {
		var input event
		err := json.Unmarshal(raw, &input)
		if err != nil {
			return nil, err
		}

		return response{Body: "hello " + input.Name}, nil
	})
}
