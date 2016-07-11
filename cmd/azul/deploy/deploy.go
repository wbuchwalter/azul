package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/wbuchwalter/azul/app"
	"github.com/wbuchwalter/azul/cmd/azul/root"
	"github.com/wbuchwalter/azul/function"
)

var deployCmd = &cobra.Command{
	Use:   "deploy funcname",
	Short: "deploy the specified function to Azure",
	Long:  `deploy - Compile and deploy the specified go function to Azure`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = preRun(args, wd)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = run(args, wd)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Succesfully deployed.")
	},
}

func init() {
	root.RootCmd.AddCommand(deployCmd)
}

func preRun(args []string, wd string) error {
	if _, err := os.Stat(wd + "/azul.json"); os.IsNotExist(err) {
		return errors.New("azul.json file not found")
	}

	for i := 0; i < len(args); i++ {
		err := checkFunctionSanity(wd, args[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func checkFunctionSanity(wd string, dirname string) error {
	if _, err := os.Stat(wd + "/" + dirname); os.IsNotExist(err) {
		return errors.New("Function " + dirname + " not found.")
	}
	return nil
}

func run(args []string, wd string) error {
	app, err := getApp(wd)
	if err != nil {
		return err
	}

	for i := 0; i < len(args); i++ {
		f := function.Function{Name: args[i], Path: wd + "/" + args[i] + "/"}
		err := app.Deploy(&f)
		if err != nil {
			return err
		}
	}
	return nil
}

func getApp(wd string) (*app.App, error) {
	file, err := os.Open(wd + "/azul.json")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var app app.App
	err = json.Unmarshal(b, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
