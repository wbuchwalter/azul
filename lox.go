package lox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wbuchwalter/lox/logs"
)

type handlerFn func(req json.RawMessage, logger logs.Logger) (interface{}, error)

func Handle(fn handlerFn, fnName string) {
	input, err := ioutil.ReadFile("D:/home/site/wwwroot/" + fnName + "/tmp")
	if err != nil {
		os.Stderr.WriteString("[Error] reading function's input: " + err.Error())
		return
	}

	logger := logs.Logger{Logs: make(chan string, 200)}

	output, err := fn(input, logger)

	//currently we only logs everything once the function is finished, not real time, this isnt great
	logger.Kill()
	logger.WriteToFile(os.Stderr)

	if err != nil {
		os.Stderr.WriteString("[Error] executing the function: " + err.Error())
		return
	}

	json, err := json.Marshal(output)
	if err != nil {
		os.Stderr.WriteString("[Error] marshaling output: " + err.Error())
		return
	}

	fmt.Print(string(json))
	return
}
