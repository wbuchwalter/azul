package main

import (
	"encoding/json"

	"github.com/wbuchwalter/lox"
	"github.com/wbuchwalter/lox/logs"
)

type input struct {
	Comment comment `json:"comment"`
}

type comment struct {
	Body string `json:"body"`
}

type Output struct {
	Body string `json:"body"`
}

func main() {
	lox.Handle(func(req json.RawMessage, logger logs.Logger) (interface{}, error) {
		var in input
		var out Output

		err := json.Unmarshal(req, &in)
		if err != nil {
			return nil, err
		}

		logger.Log("GitHub Webhook triggered! " + in.Comment.Body)

		out.Body = "New GitHub comment: " + in.Comment.Body
		return out, nil
	})
}
