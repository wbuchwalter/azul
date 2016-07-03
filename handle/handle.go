package lox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type input struct {
	Body json.RawMessage `json:"body"`
}

type handlerFn func(req json.RawMessage) ([]byte, int)

func Handle(fn handlerFn) {
	var i input
	fmt.Println("OPENING FILE")
	data, err := ioutil.ReadFile("D:/home/site/wwwroot/HttpTriggerCSharp1/tmp")
	if err != nil {
		fmt.Println("Cannot read event file: ", err)
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		fmt.Println("ERR UNMARSHALING: ", err)
		fmt.Println("DATA WAS: ", data)
		return
	}

	fmt.Println("Body: ", len(string(i.Body)))
	ret, status := fn(i.Body)

	fmt.Println("Returns: ", string(ret), " Status: ", status)
}
