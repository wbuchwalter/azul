package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Function struct describe an Azure Functions
type Function struct {
	Name string
}

//GetFunctions returns the functions inside a function app
func GetFunctions(appName string) ([]Function, error) {
	resp, err := http.Get("funcwill.scm.azurewebsites.net/api/functions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)

	var data []Function
	return data, json.Unmarshal(buffer, &data)
}

//Delete deletes the specified function within a function app
func Delete(functionName string) error {
	req, err := http.NewRequest(http.MethodDelete, "https://@funcwill.scm.azurewebsites.net/api/functions/"+functionName, nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 204: //expected
		return nil
	case 404:
		fmt.Println("Warning: Tried to delete unexisting function ", functionName)
	default:
		return errors.New("Unhandled error while trying to delete " + functionName + ". Received status " + strconv.Itoa(resp.StatusCode) + ", expected 204.")
	}

	return nil
}
