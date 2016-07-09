package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/wbuchwalter/lox/utils"
	"github.com/wbuchwalter/lox/webjobs-sdk/function"
	"github.com/wbuchwalter/lox/webjobs-sdk/vfs"
)

type Function struct {
	Name     string
	Path     string
	FilesMap map[string]string
}

func (f *Function) Deploy() error {
	err := f.delete()
	if err != nil {
		return err
	}

	fMap := f.getPredefinedFiles()
	if conf, ok, err := f.getCustomConfig(); ok {
		fMap["function.json"] = conf
	} else if err != nil {
		return err
	}

	f.FilesMap = fMap

	err = f.create()
	if err != nil {
		return err
	}

	bin, err := f.build()
	if err != nil {
		return err
	}

	return f.uploadBinary(bin, "main.exe")
}

//build main.go in a temp folder, read the bytes, delete the file
func (f *Function) build() ([]byte, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	fmt.Println("Build")
	dst := dir + "/.lox/main.exe"
	buildCmd := "GOOS=windows GOARCH=386 go build " + "-o " + dst + " " + f.Path + "main.go"
	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	bin, err := utils.GetFileBytes(dst)
	if err != nil {
		return nil, err
	}

	return bin, os.Remove(dst)
}

func (f *Function) getCustomConfig() (string, bool, error) {
	if _, err := os.Stat(f.Path + "function.json"); os.IsNotExist(err) {
		return "", false, nil
	}

	rc, err := os.Open(f.Path + "function.json")
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

func (f *Function) delete() error {
	return function.Delete(f.Name)
}

func (f *Function) uploadBinary(bin []byte, fileName string) error {
	return vfs.PushFile(bin, fileName, f.Name)
}

func (f *Function) create() error {
	config := json.RawMessage(f.FilesMap["function.json"])
	delete(f.FilesMap, "function.json")

	dto := function.CreateFunctionDTO{Config: &config, Files: f.FilesMap}
	return function.Create(f.Name, dto)
}

func (f *Function) getPredefinedFiles() map[string]string {
	fMap := make(map[string]string)

	fMap["function.json"] = `{ "bindings": [ { "authLevel": "function", "name": "req", "type": "httpTrigger", "direction": "in" }, { "name": "res", "type": "http", "direction": "out" } ], "disabled": false }`
	fMap["project.json"] = `{ "frameworks": { "net46":{ "dependencies": { "Newtonsoft.Json": "9.0.1" } } } }`
	fMap["run.csx"] = `using System.Net;
using System.Diagnostics;
using System;
using System.IO;
using Newtonsoft.Json;
using System.Text;
using Newtonsoft.Json.Linq;

public static async Task<HttpResponseMessage> Run(HttpRequestMessage req, TraceWriter log)
{
    Process process = new Process();
    process.StartInfo.FileName = "D:/home/site/wwwroot/` + f.Name + `/main.exe";

    var data = await req.Content.ReadAsStringAsync();
    await WriteToFileAsync(data, log);
    
    process.StartInfo.RedirectStandardOutput = true;
		process.StartInfo.RedirectStandardError = true;
    process.StartInfo.UseShellExecute = false;
    process.Start(); 
    string json = "";
    string err = "";
    while ( ! process.HasExited ) {
        json += process.StandardOutput.ReadToEnd();
        err += process.StandardError.ReadToEnd();
    }

		if(err != "") {
	    return req.CreateResponse((HttpStatusCode)500, err);
		} else {
    	return req.CreateResponse((HttpStatusCode)200, json);			
		}
}

static async Task WriteToFileAsync(string text, TraceWriter log)
{
    byte[] buffer = Encoding.UTF8.GetBytes(text);
    log.Info("in:" + text);
    Int32 offset = 0;
    Int32 sizeOfBuffer = 4096;
    FileStream fileStream = null;

    fileStream = new FileStream("D:/home/site/wwwroot/` + f.Name + `/tmp", FileMode.Create, FileAccess.Write, FileShare.None, bufferSize: sizeOfBuffer, useAsync: true);
    await fileStream.WriteAsync(buffer, offset, buffer.Length);
    fileStream.Dispose();
}`

	return fMap
}
