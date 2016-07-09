package functionClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

//GetFunctions returns the functions inside a function app
func GetFunctions(appName string) ([]string, error) {
	resp, err := http.Get("https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/functions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)

	var data []string
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

type CreateFunctionInput struct {
	FunctionName string            `json:"-"`
	Files        map[string]string `json:"files"`
	Config       *json.RawMessage  `json:"config"`
}

//Create create a new function within a function app
func Create(input CreateFunctionInput) error {
	m, err := json.Marshal(&input)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, "https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/functions/"+input.FunctionName, bytes.NewBuffer(m))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = (&http.Client{}).Do(req)
	return err
}
