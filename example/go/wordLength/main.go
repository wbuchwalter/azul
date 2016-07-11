package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul"
	"github.com/wbuchwalter/azul/logs"
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

		output.Length = len(i.Word)

		return output, nil
	})
}
