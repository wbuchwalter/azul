package vfs

import (
	"io"
	"net/http"
)

func PushFile(rc io.Reader, fileName string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, "https://$funcwill:2xXQ3heWo7dD3mSmlvLhZnwzqJXMmrwHxugRCrnAnCb0idmo2vXCbiLMqqtY@funcwill.scm.azurewebsites.net/api/vfs/site/wwwroot/EmptyNodeJS1/"+fileName, rc)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}
