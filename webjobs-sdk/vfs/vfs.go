package vfs

import (
	"bytes"
	"net/http"
)

func PushFile(buf []byte, fileName string, fnName string) error {
	rc := bytes.NewBuffer(buf)
	req, err := http.NewRequest(http.MethodPut, "https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/vfs/site/wwwroot/"+fnName+"/"+fileName, rc)
	if err != nil {
		return err
	}
	_, err = (&http.Client{}).Do(req)
	return err
}
