package basic

import (
	"io/ioutil"
	"os"

	"github.com/wbuchwalter/azul/function"
)

//Build main.go in a temp folder, read the bytes, delete the file
func Build(f *function.Function) (function.FilesMap, function.Config, error) {
	fMap, err := getFiles(f)
	if err != nil {
		return nil, function.Config{}, err
	}
	return fMap, getDefaultConfig(), nil
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

func getFiles(f *function.Function) (function.FilesMap, error) {
	fMap := make(function.FilesMap)
	fi, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}

	for _, fi := range fi {
		fileHandle, err := os.Open(f.Path + fi.Name())
		defer fileHandle.Close()
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(fileHandle)
		if err != nil {
			return nil, err
		}

		fMap[fi.Name()] = function.FuncFile{IsHeavy: false, Content: string(b)}
	}

	return fMap, nil
}
