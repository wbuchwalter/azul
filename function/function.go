package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/wbuchwalter/lox/auth"
)

func Deploy(filePath string) error {
	err := build(filePath)
	if err != nil {
		return err
	}

	authInfo := auth.GetAuthInfo("test", "test")
	fmt.Println(authInfo)
	return err
}

func build(filePath string) error {
	cmd := exec.Command("sh", "-c", "GOOS=windows GOARCH=386 go build "+filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// func

//push to https://${name}.scm.azurewebsites.net/api/vfs/site/wwwroot/${funcName}/${filepath}

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
