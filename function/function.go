package function

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/wbuchwalter/lox/service/function"
	"github.com/wbuchwalter/lox/service/vfs"
	"github.com/wbuchwalter/lox/utils"
)

//Function represents an Azure Function
type Function struct {
	Config
	Name string
	Path string
}

//Config represents a function's configuration (function.json)
type Config struct {
	Bindings []Binding `json:"bindings"`
	Disabled bool      `json:"disabled"`
}

//Binding represents a function's binding
type Binding struct {
	AuthLevel string `json:"authLevel"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
}

type filesMap map[string]string

//Deploy (re)deploys a function inside a functionApp
func (f *Function) Deploy() error {
	err := f.loadConfig()
	if err != nil {
		return err
	}

	fMap := f.getPredefinedFiles()

	//build before anything else.
	//If the custom code doesn't build, we can fail cleanly without impacting what is already deployed
	bin, err := f.build()
	if err != nil {
		return err
	}

	//delete the function if it already exists
	err = f.delete()
	if err != nil {
		return err
	}

	//create the function and push config + predefined files
	err = f.create(fMap)
	if err != nil {
		return err
	}

	//upload the heavier executable
	return f.uploadBinary(bin, "main.exe")
}

//build main.go in a temp folder, read the bytes, delete the file
func (f *Function) build() ([]byte, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
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

func (f *Function) loadConfig() error {
	f.defaults()
	if _, err := os.Stat(f.Path + "function.json"); os.IsNotExist(err) {
		return nil
	}

	rc, err := os.Open(f.Path + "function.json")
	defer rc.Close()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, f.Config)
}

func (f *Function) defaults() {
	f.Config.Disabled = false
	f.Config.Bindings = []Binding{
		Binding{AuthLevel: "function", Name: "req", Type: "httpTrigger", Direction: "in"},
		Binding{Name: "res", Type: "http", Direction: "out"},
	}
}

func (f *Function) delete() error {
	return functionClient.Delete(f.Name)
}

func (f *Function) uploadBinary(bin []byte, fileName string) error {
	return vfsClient.PushFile(bin, fileName, f.Name)
}

func (f *Function) create(fMap filesMap) error {
	config, err := json.Marshal(f.Config)
	if err != nil {
		return err
	}
	rm := json.RawMessage(config)
	return functionClient.Create(functionClient.CreateFunctionInput{FunctionName: f.Name, Config: &rm, Files: fMap})
}

func (f *Function) getPredefinedFiles() filesMap {
	fMap := make(filesMap)

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
