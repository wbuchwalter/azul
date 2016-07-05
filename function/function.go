package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/wbuchwalter/lox/webjobs-sdk/function"
)

type functionBinary struct {
	fileName string
	data     []byte
}

func Deploy(functionName string, functionDir string) error {
	bins, err := getBinaries(functionDir)
	err = delete(functionName)
	if err != nil {
		return err
	}
	return create(functionName, bins)
}

func getBinaries(functionDir string) ([]functionBinary, error) {
	b, err := build(functionDir + "main.go")
	mBin := functionBinary{"main.exe", b}
	if err != nil {
		return nil, err
	}

	//check for function.json
	if _, err := os.Stat(functionDir + "function.json"); os.IsNotExist(err) {
		return []functionBinary{mBin}, nil
	}

	rc, err := os.Open(functionDir + "function.json")
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	fB, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	fBin := functionBinary{"function.json", fB}
	return []functionBinary{mBin, fBin}, nil
}

//build main.go in a temp folder, read the bytes, delete the file
func build(filePath string) ([]byte, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	dst := dir + "/.lox/main.exe"
	buildCmd := "GOOS=windows GOARCH=386 go build " + "-o " + dst + " " + filePath
	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	bin, err := getBytes(dst)
	if err != nil {
		return nil, err
	}

	return bin, os.Remove(dst)
}

func getBytes(path string) ([]byte, error) {
	rc, err := os.Open(path)
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	bin, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

func delete(functionName string) error {
	return function.Delete(functionName)
}

func create(functionName string, bins []functionBinary) error {
	fmap := make(map[string]string)
	fmap["run.csx"] = `"hey!"`
	dto := function.CreateFunctionDTO{Config: json.RawMessage(`{"bindings":[],"disabled":false}`), Files: fmap}
	return function.Create(functionName, dto)
}

//--------------------

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
