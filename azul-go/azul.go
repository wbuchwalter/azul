package azul

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/wbuchwalter/azul/azul-go/logs"
)

type handlerFn func(req json.RawMessage, logger logs.Logger) (interface{}, error)

var FunctionName string

func Handle(fn handlerFn) {
	input, err := ioutil.ReadFile("D:/home/site/wwwroot/" + FunctionName + "/tmp")
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

	//if returned type is a struct, we marshal it to json
	//otherwise we return it raw
	if reflect.TypeOf(output).Kind() == reflect.Struct {
		raw, err := json.Marshal(&output)
		if err != nil {
			os.Stderr.WriteString("[Error] marshaling output: " + err.Error())
		}
		fmt.Print(string(raw))
	} else {
		fmt.Print(output)
	}

	return
}
