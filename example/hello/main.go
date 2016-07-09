package main

import (
	"encoding/json"

	"github.com/wbuchwalter/lox"
)

type input struct {
	Name string `json:"name"`
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

		output.Length = len(i.Name)

		return output, nil
	}, "hello")
}
