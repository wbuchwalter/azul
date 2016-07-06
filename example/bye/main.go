package main

import (
	"encoding/json"
	"fmt"

	"github.com/wbuchwalter/lox/function"
)

type message struct {
	Name string `json:"name"`
}

func main() {
	function.Handle(func(event json.RawMessage) ([]byte, int) {
		var m message

		err := json.Unmarshal(event, &m)
		if err != nil || m.Name == "" {
			fmt.Println(err)
			return nil, 503
		}

		return []byte(`{"val": "` + m.Name + `"}`), 200
	})
}
