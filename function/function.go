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

var runBin = `using System.Net; using System.Diagnostics; using System; using System.IO; using Newtonsoft.Json; using System.Text; public static async Task<HttpResponseMessage> Run(HttpRequestMessage req, TraceWriter log) { Process process = new Process(); process.StartInfo.FileName = "D:/home/site/wwwroot/HttpTriggerCSharp1/main.exe"; var data = await req.Content.ReadAsStringAsync(); await WriteToFileAsync(data, log); //process.StartInfo.Arguments = "-r " + data; process.StartInfo.RedirectStandardOutput = true; process.StartInfo.UseShellExecute = false; process.Start(); string q = ""; while ( ! process.HasExited ) { q += process.StandardOutput.ReadToEnd(); } log.Info(q); return req.CreateResponse(HttpStatusCode.OK, "Hello "); } static async Task WriteToFileAsync(string text, TraceWriter log) { byte[] buffer = Encoding.UTF8.GetBytes(text); log.Info("BUFFER: "+ buffer[0] + " " + buffer[1]); Int32 offset = 0; Int32 sizeOfBuffer = 4096; FileStream fileStream = null; fileStream = new FileStream("D:/home/site/wwwroot/HttpTriggerCSharp1/tmp", FileMode.Create, FileAccess.Write, FileShare.None, bufferSize: sizeOfBuffer, useAsync: true); await fileStream.WriteAsync(buffer, offset, buffer.Length); fileStream.Dispose(); }`
var defaultFunctionJSONBin = `{ "bindings": [ { "authLevel": "function", "name": "req", "type": "httpTrigger", "direction": "in" }, { "name": "res", "type": "http", "direction": "out" } ], "disabled": false }`
var projectJSONBin = `{ "frameworks": { "net46":{ "dependencies": { "Newtonsoft.Json": "9.0.1" } } } }`

func Deploy(functionName string, functionDir string) error {
	bins, err := getBinaries(functionDir)
	//err = deleteFunction(functionName)
	if err != nil {
		return err
	}
	return createFunction(functionName, bins)
}

func getBinaries(functionDir string) (map[string]string, error) {
	binMap := make(map[string]string)
	b, err := build(functionDir + "main.go")
	if err != nil {
		return nil, err
	}
	binMap["main.exe"] = string(b)
	binMap["project.json"] = projectJSONBin
	binMap["run.csx"] = runBin

	//check for function.json
	if _, err := os.Stat(functionDir + "function.json"); os.IsNotExist(err) {
		binMap["function.json"] = defaultFunctionJSONBin
		return binMap, nil
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
	binMap["function.json"] = string(fB)
	return binMap, nil
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

func deleteFunction(functionName string) error {
	return function.Delete(functionName)
}

func createFunction(functionName string, bins map[string]string) error {
	fnBin := json.RawMessage(bins["function.json"])
	delete(bins, "function.json")

	dto := function.CreateFunctionDTO{Config: &fnBin, Files: bins}
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
