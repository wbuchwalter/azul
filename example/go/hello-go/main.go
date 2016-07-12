package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul"
	"github.com/wbuchwalter/azul/logs"
)

type event struct {
	Name string `json:"name"`
}

type response struct {
	Body string `json:"body"`
}

func main() {
	azul.Handle(func(raw json.RawMessage, logger logs.Logger) (interface{}, error) {
		var input event
		err := json.Unmarshal(raw, &input)
		if err != nil {
			return nil, err
		}

		return response{Body: "hello " + input.Name}, nil
	})
}
