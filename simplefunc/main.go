package main

import (
	"encoding/json"
	"fmt"

	"github.com/wbuchwalter/bluefunc"
)

type message struct {
	Name string `json:"name"`
}

func main() {
	bluefunc.Handle(func(event json.RawMessage) ([]byte, int) {
		var m message
//t
		err := json.Unmarshal(event, &m)
		if err != nil || m.Name == "" {
			fmt.Println(err)
			return nil, 503
		}

		return []byte("Hello " + m.Name), 200
	})
}
