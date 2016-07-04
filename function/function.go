package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/wbuchwalter/lox/webjobs-sdk/function"
)

func Deploy(functionName string) error {

	err := deleteFunction(functionName)
	// if err != nil {
	// 	return err
	// }

	// err = build(functionName + "main.go")
	// if err != nil {
	// 	return err
	// }

	// reader, err := os.Open(functionName + "main.exe")
	// if err != nil {
	// 	return err
	// }

	// resp, err := vfs.PushFile(reader, "main.exe")
	// fmt.Println(resp)
	return err
}

func build(filePath string) error {
	cmd := exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build "+filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

type input struct {
	Body json.RawMessage `json:"body"`
}

type handlerFn func(req json.RawMessage) ([]byte, int)

func Handle(fn handlerFn) error {
	var i input
	data, err := ioutil.ReadFile("D:/home/site/wwwroot/HttpTriggerCSharp1/tmp")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	ret, status := fn(i.Body)
	fmt.Println(status, " ", ret)
	return nil
}

func deleteFunction(functionName string) error {
	return function.Delete("funcwill")
}
