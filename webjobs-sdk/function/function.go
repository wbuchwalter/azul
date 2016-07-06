package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Function struct describe an Azure Functions
type Function struct {
	Name string
}

type CreateFunctionDTO struct {
	Files  map[string]string `json:"files"`
	Config *json.RawMessage  `json:"config"`
}

//GetFunctions returns the functions inside a function app
func GetFunctions(appName string) ([]Function, error) {
	resp, err := http.Get("https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/functions")
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
	req, err := http.NewRequest(http.MethodDelete, "https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/functions/"+functionName, nil)
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

//Create create a new function within a function app
func Create(functionName string, dto CreateFunctionDTO) error {
	m, err := json.Marshal(&dto)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, "https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/functions/"+functionName, bytes.NewBuffer(m))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = (&http.Client{}).Do(req)
	return err
}
