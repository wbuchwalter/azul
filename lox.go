package lox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type input struct {
	Body json.RawMessage `json:"body"`
}

type handlerFn func(req json.RawMessage) (interface{}, error)

func Handle(fn handlerFn, fnName string) {
	var i input
	data, err := ioutil.ReadFile("D:/home/site/wwwroot/" + fnName + "/tmp")
	if err != nil {
		os.Stderr.WriteString("Error reading function's input: " + err.Error())
		return
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		os.Stderr.WriteString("Error unmarshalling function's input: " + err.Error())
		return
	}

	output, err := fn(i.Body)
	if err != nil {
		os.Stderr.WriteString("Error executing the function: " + err.Error())
		return
	}

	json, err := json.Marshal(output)
	if err != nil {
		os.Stderr.WriteString("Error marshaling output: " + err.Error())
		return
	}

	fmt.Print(string(json))
	return
}
