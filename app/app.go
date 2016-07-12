package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/wbuchwalter/azul/build"
	"github.com/wbuchwalter/azul/function"
)

//App represents an Azure Function App's configuration
type App struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Deploy (re)deploys a function
func (app *App) Deploy(f *function.Function) error {
	fMap, defaultConf, err := build.Build(f)
	if err != nil {
		return err
	}

	err = f.LoadConfig(defaultConf)
	if err != nil {
		return err
	}

	//contains all 'light' files as string
	lightMap := make(map[string]string)

	//contains all `heavy` fils as bytes (such as executables)
	heavyMap := make(map[string][]byte)

	for n, fi := range fMap {
		if fi.IsHeavy {
			heavyMap[n] = fi.BContent
		} else {
			lightMap[n] = fi.Content
		}
	}

	//create the function with config + 'light' files
	err = app.create(f, lightMap)
	if err != nil {
		return err
	}

	//upload the heavier files throught vfs api
	for n, b := range heavyMap {
		err = app.pushFile(b, n, f.Name)
		if err != nil {
			return err
		}
	}

	//TEMPORARY - BUG: project.json needs to be remodified to trigger a nuget restore...
	if val, ok := fMap["project.json"]; ok {
		return app.forceNugetRestore(f.Name, val.Content)
	}

	fmt.Println("https://" + app.Name + ".azurewebsites.net/api/" + f.Name)
	return nil
}

type logInfo struct {
	Href string `json:"href"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

//Logs returns the logs of the specified function
func (app *App) Logs(f *function.Function) error {
	serviceURI := app.getAPIBaseURI() + "vfs/logfiles/application/functions/function/hello-node/"
	res, err := http.Get(serviceURI)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var logInfos []logInfo
	err = json.Unmarshal(b, &logInfos)
	if err != nil {
		return err
	}

	if len(logInfos) < 1 {
		return errors.New("No logfiles available")
	}

	//always grab the first logfile for now...
	lres, err := http.Get(serviceURI + logInfos[0].Name)
	defer lres.Body.Close()
	b, err = ioutil.ReadAll(lres.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (app *App) LogStream() error {

	res, err := http.Get(app.getBaseURI() + "logstream")
	if err != nil {
		return err
	}

	//reader := bufio.NewReader(res.Body)
	defer res.Body.Close()
	p := make([]byte, 1)
	_, err = res.Body.Read(p)
	if err != nil {
		return err
	}

	return nil
}

//CreateFunctionInput : input for app.create
type CreateFunctionInput struct {
	FunctionName string            `json:"-"`
	RawFiles     map[string]string `json:"files"`
	Config       *json.RawMessage  `json:"config"`
}

func (app *App) create(f *function.Function, rawMap map[string]string) error {
	config, err := json.Marshal(f.Config)
	if err != nil {
		return err
	}
	rc := json.RawMessage(config)

	input := CreateFunctionInput{FunctionName: f.Name, Config: &rc, RawFiles: rawMap}
	m, err := json.Marshal(&input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, app.getAPIBaseURI()+"functions/"+input.FunctionName, bytes.NewBuffer(m))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = (&http.Client{}).Do(req)
	return err
}

//GetFunctions returns the functions inside a function app
func (app *App) GetFunctions() ([]string, error) {
	resp, err := http.Get(app.getAPIBaseURI() + "functions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)

	var data []string
	return data, json.Unmarshal(buffer, &data)
}

//Delete deletes the specified function within a function app
func (app *App) Delete(functionName string) error {
	req, err := http.NewRequest(http.MethodDelete, app.getAPIBaseURI()+"functions/"+functionName, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 204, 404: //expected, 404 means the function didn't existed in the first place
		return nil
	default:
		return errors.New("Unhandled error while trying to delete " + functionName + ". Received status " + strconv.Itoa(resp.StatusCode))
	}
}

//PushFile uploads a binary to the specified function
func (app *App) pushFile(buf []byte, fileName string, fnName string) error {
	rc := bytes.NewBuffer(buf)
	req, err := http.NewRequest(http.MethodPut, app.getAPIBaseURI()+"/vfs/site/wwwroot/"+fnName+"/"+fileName, rc)
	if err != nil {
		return err
	}
	_, err = (&http.Client{}).Do(req)
	return err
}

func (app *App) forceNugetRestore(functionName string, projectJSONContent string) error {
	req, err := http.NewRequest(http.MethodPut, app.getAPIBaseURI()+"/vfs/site/wwwroot/"+functionName+"/project.json", bytes.NewBuffer([]byte(" "+projectJSONContent)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("If-Match", "*")

	_, err = (&http.Client{}).Do(req)
	return err
}

func (app *App) getAPIBaseURI() string {
	return app.getBaseURI() + "api/"
}

func (app *App) getBaseURI() string {
	return "https://" + app.Username + ":" + app.Password + "@" + app.Name + ".scm.azurewebsites.net/"
}
