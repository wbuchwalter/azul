package golang

import (
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/wbuchwalter/azul/function"
	"github.com/wbuchwalter/azul/utils"
)

//Build main.go in a temp folder, read the bytes, delete the file
func Build(f *function.Function) (function.FilesMap, function.Config, error) {
	var conf function.Config
	fMap := getBootstrapFiles(f)

	dir, err := homedir.Dir()
	if err != nil {
		return nil, conf, err
	}
	dst := dir + "/.azul/main.exe"

	//-ldflags should be removed once C# -> Go communication is not done via a file anymore
	buildCmd := `GOOS=windows GOARCH=386 go build -ldflags="-X github.com/wbuchwalter/azul/azul-go.FunctionName=` + f.Name + `" ` + "-o " + dst + " " + f.Path + "main.go"
	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, conf, err
	}
	bin, err := utils.GetFileBytes(dst)
	if err != nil {
		return nil, conf, err
	}

	fMap["main.exe"] = function.FuncFile{IsHeavy: true, BContent: bin}

	return fMap, getDefaultConfig(), os.Remove(dst)
}

func getDefaultConfig() function.Config {
	var c function.Config
	c.Disabled = false
	c.Bindings = []function.Binding{
		function.Binding{AuthLevel: "anonymous", Name: "req", Type: "httpTrigger", Direction: "in"},
		function.Binding{Name: "res", Type: "http", Direction: "out"},
	}
	return c
}

func getBootstrapFiles(f *function.Function) function.FilesMap {
	fMap := make(function.FilesMap)
	fMap["project.json"] = function.FuncFile{IsHeavy: false, Content: `{ "frameworks": { "net46":{ "dependencies": { "Newtonsoft.Json": "9.0.1" } } } }`}
	fMap["run.csx"] = function.FuncFile{IsHeavy: false, Content: `using System.Net;
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
      string output = "";
      string err = "";
      while ( ! process.HasExited ) {
          output += process.StandardOutput.ReadToEnd();
          err += process.StandardError.ReadToEnd();
      }

      Regex regex = new Regex(@"^\[Error\]");
      Match match = regex.Match(err);
      if (match.Success) {
        log.Error(err);
        return req.CreateResponse((HttpStatusCode)500, err);
      } else {
        log.Info(err);
        return req.CreateResponse((HttpStatusCode)200, output);	
      }
  }

  static async Task WriteToFileAsync(string text, TraceWriter log)
  {
      byte[] buffer = Encoding.UTF8.GetBytes(text);
      Int32 offset = 0;
      Int32 sizeOfBuffer = 4096;
      FileStream fileStream = null;

      fileStream = new FileStream("D:/home/site/wwwroot/` + f.Name + `/tmp", FileMode.Create, FileAccess.Write, FileShare.None, bufferSize: sizeOfBuffer, useAsync: true);
      await fileStream.WriteAsync(buffer, offset, buffer.Length);
      fileStream.Dispose();
  }`}

	return fMap
}
