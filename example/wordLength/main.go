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
		if err != nil {
			return nil, err
		}

		output.Length = len(i.Word)

		return output, nil
	}, "wordLength")
}
