package utils

import (
	"io/ioutil"
	"os"
)

func GetFileBytes(path string) ([]byte, error) {
	rc, err := os.Open(path)
	defer rc.Close()
	if err != nil {
		return nil, err
	}
	bin, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return bin, nil
}
