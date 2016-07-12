package function

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Function represents an Azure Function
type Function struct {
	Config
	Name string
	Path string
}

//Config represents a function's configuration (function.json)
type Config struct {
	Bindings []Binding `json:"bindings"`
	Disabled bool      `json:"disabled"`
}

//Binding represents a function's binding
type Binding struct {
	AuthLevel   string `json:"authLevel,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Direction   string `json:"direction"`
	WebHookType string `json:"webHookType,omitempty"`
}

type FuncFile struct {
	IsHeavy  bool
	BContent []byte
	Content  string
}

type FilesMap map[string]FuncFile

// LoadConfig loads the config of the function stored in function.json.
// If no function.json is present, this will fallback to default values
func (f *Function) LoadConfig(defaultConfig Config) error {

	if _, err := os.Stat(f.Path + "function.json"); os.IsNotExist(err) {
		f.Config = defaultConfig
		return nil
	}

	rc, err := os.Open(f.Path + "function.json")
	defer rc.Close()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &f.Config)
}
