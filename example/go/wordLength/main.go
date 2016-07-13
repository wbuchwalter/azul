package main

import (
	"encoding/json"

	"github.com/wbuchwalter/azul/azul-go"
	"github.com/wbuchwalter/azul/azul-go/logs"
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
