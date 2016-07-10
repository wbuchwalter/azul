package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/wbuchwalter/lox/function"
)

//App represents an Azure Function App's configuration
type App struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Deploy (re)deploys a function
func (app *App) Deploy(f *function.Function) error {
	//build before anything else.
	//If the custom code doesn't build, we can fail cleanly without impacting what is already deployed
	bin, err := f.Build()
	if err != nil {
		return err
	}

	err = f.LoadConfig()
	if err != nil {
		return err
	}

	fMap := f.GetPredefinedFiles()

	//create the function and push config + predefined files
	err = app.create(f, fMap)
	if err != nil {
		return err
	}

	//upload the heavier executable
	err = app.pushFile(bin, "main.exe", f.Name)
	if err != nil {
		return err
	}

	//TEMPORARY - BUG: project.json needs to be remodified to trigger a nuget restore...
	return app.forceNugetRestore(f.Name, fMap["project.json"])
}

type CreateFunctionInput struct {
	FunctionName string            `json:"-"`
	Files        map[string]string `json:"files"`
	Config       *json.RawMessage  `json:"config"`
}

func (app *App) create(f *function.Function, fMap map[string]string) error {
	config, err := json.Marshal(f.Config)
	if err != nil {
		return err
	}
	rc := json.RawMessage(config)
	input := CreateFunctionInput{FunctionName: f.Name, Config: &rc, Files: fMap}
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
	return nil
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
	return "https://" + app.Username + ":" + app.Password + "@" + app.Name + ".scm.azurewebsites.net/api/"
}
