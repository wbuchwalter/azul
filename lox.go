package lox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type handlerFn func(req json.RawMessage) (interface{}, error)

func Handle(fn handlerFn, fnName string) {
	input, err := ioutil.ReadFile("D:/home/site/wwwroot/" + fnName + "/tmp")
	if err != nil {
		os.Stderr.WriteString("Error reading function's input: " + err.Error())
		return
	}

	output, err := fn(input)
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
