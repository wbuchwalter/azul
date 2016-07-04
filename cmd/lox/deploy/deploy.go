package deploy

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wbuchwalter/lox/cmd/lox/root"
	"github.com/wbuchwalter/lox/function"
)

var deployCmd = &cobra.Command{
	Use:   "deploy main.go",
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
		}
	},
}

func init() {
	root.RootCmd.AddCommand(deployCmd)
}

func preRun(args []string, wd string) error {
	if _, err := os.Stat(wd + "/lox.json"); os.IsNotExist(err) {
		return errors.New("lox.json file not found")
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

	if _, err := os.Stat(wd + "/" + dirname + "/main.go"); os.IsNotExist(err) {
		return errors.New("Function " + dirname + " found, but no main.go was present.")
	}

	return nil
}

func run(args []string, wd string) error {
	for i := 0; i < len(args); i++ {
		err := function.Deploy(wd + "/" + args[i] + "/")
		if err != nil {
			return err
		}
	}
	return nil
}
