package main

import (
	"encoding/json"

	"github.com/wbuchwalter/lox/function"
)

type input struct {
	Name string `json:"name"`
}

type Output struct {
	Message string `json:"message"`
	Target  string `json:"target"`
}

func main() {
	function.Handle(func(event json.RawMessage) (interface{}, error) {
		var i input
		var output Output

		err := json.Unmarshal(event, &i)
		if err != nil || i.Name == "" {
			return nil, err
		}

		output.Message = "Hey!"
		output.Target = i.Name
		return output, nil
	})
}
