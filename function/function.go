package function

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
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
	AuthLevel   string `json:"authLevel"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Direction   string `json:"direction"`
	WebHookType string `json:"webHookType"`
}

type filesMap map[string]string

//Build main.go in a temp folder, read the bytes, delete the file
func (f *Function) Build() ([]byte, error) {
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

// LoadConfig loads the config of the function stored in function.json.
// If no function.json is present, this will fallback to default values
func (f *Function) LoadConfig() error {
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

	return json.Unmarshal(b, &f.Config)
}

func (f *Function) defaults() {
	f.Config.Disabled = false
	f.Config.Bindings = []Binding{
		Binding{AuthLevel: "function", Name: "req", Type: "httpTrigger", Direction: "in"},
		Binding{Name: "res", Type: "http", Direction: "out"},
	}
}

// GetPredefinedFiles returns predefined files
func (f *Function) GetPredefinedFiles() filesMap {
	fMap := make(filesMap)

	fMap["project.json"] = `{ "frameworks": { "net46":{ "dependencies": { "Newtonsoft.Json": "9.0.1" } } } }`
	fMap["run.csx"] = `using System.Net;
using System.Diagnostics;
using System;
using System.Text.RegularExpressions;
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

		Regex regex = new Regex(@"^\[Error\]");
		Match match = regex.Match(err);
		if (match.Success) {
			log.Error(err);
			return req.CreateResponse((HttpStatusCode)500, err);
		} else {
			log.Info(err);
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
