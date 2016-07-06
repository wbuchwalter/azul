package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/wbuchwalter/lox/webjobs-sdk/function"
	"github.com/wbuchwalter/lox/webjobs-sdk/vfs"
)

func Deploy(fnName string, functionDir string) error {
	fMap := getPredefinedFiles(fnName)
	err := deleteFunction(fnName)
	if err != nil {
		return err
	}
	if conf, ok, err := getCustomConfig(functionDir); !ok {
		if err != nil {
			return err
		}
	} else {
		fMap["function.json"] = conf
	}

	err = createFunction(fnName, fMap)
	if err != nil {
		return err
	}

	bin, err := getBinary(functionDir)
	if err != nil {
		return err
	}

	return uploadBinary(bin, "main.exe", fnName)
}

func getBinary(functionDir string) ([]byte, error) {
	return build(functionDir + "main.go")
}

func getCustomConfig(functionDir string) (string, bool, error) {
	if _, err := os.Stat(functionDir + "function.json"); os.IsNotExist(err) {
		return "", false, nil
	}

	rc, err := os.Open(functionDir + "function.json")
	defer rc.Close()
	if err != nil {
		return "", false, err
	}
	fB, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", false, err
	}
	return string(fB), true, nil
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

func deleteFunction(fnName string) error {
	return function.Delete(fnName)
}

func createFunction(fnName string, fMap map[string]string) error {
	config := json.RawMessage(fMap["function.json"])
	delete(fMap, "function.json")

	dto := function.CreateFunctionDTO{Config: &config, Files: fMap}
	return function.Create(fnName, dto)
}

func uploadBinary(bin []byte, fileName string, fnName string) error {
	return vfs.PushFile(bin, fileName, fnName)
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

	ret, _ := fn(i.Body)
	fmt.Println(string(ret))
	return nil
}

//---------------------------------
func getPredefinedFiles(fnName string) map[string]string {
	fMap := make(map[string]string)

	fMap["function.json"] = `{ "bindings": [ { "authLevel": "function", "name": "req", "type": "httpTrigger", "direction": "in" }, { "name": "res", "type": "http", "direction": "out" } ], "disabled": false }`
	fMap["project.json"] = `{ "frameworks": { "net46":{ "dependencies": { "Newtonsoft.Json": "9.0.1" } } } }`
	fMap["run.csx"] = `using System.Net;
	using System.Diagnostics;
	using System;
	using System.IO;
	using Newtonsoft.Json;
	using System.Text;

	public static async Task<HttpResponseMessage> Run(HttpRequestMessage req, TraceWriter log)
	{
			Process process = new Process();
			process.StartInfo.FileName = "D:/home/site/wwwroot/` + fnName + `/main.exe";

			var data = await req.Content.ReadAsStringAsync();
			await WriteToFileAsync(data, log);
			
			process.StartInfo.RedirectStandardOutput = true;
			process.StartInfo.UseShellExecute = false;
			process.Start();
			string output = "";
			while ( ! process.HasExited ) {
					output += process.StandardOutput.ReadToEnd();
			}

			//JObject o = JObject.Parse(output);

			// string something = (string)o["something"];
			
			// log.Info(something);
		
			return req.CreateResponse(HttpStatusCode.OK, output);
	}

	static async Task WriteToFileAsync(string text, TraceWriter log)
	{
			byte[] buffer = Encoding.UTF8.GetBytes(text);
			log.Info("BUFFER: "+ buffer[0] + " " + buffer[1]);
			Int32 offset = 0;
			Int32 sizeOfBuffer = 4096;
			FileStream fileStream = null;

			fileStream = new FileStream("D:/home/site/wwwroot/` + fnName + `/tmp", FileMode.Create, FileAccess.Write, FileShare.None, bufferSize: sizeOfBuffer, useAsync: true);
			await fileStream.WriteAsync(buffer, offset, buffer.Length);
			fileStream.Dispose();
	}`

	return fMap
}
