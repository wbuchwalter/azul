package deploy

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/wbuchwalter/azul/cmd/azul/root"
	"github.com/wbuchwalter/azul/cmd/azul/utils"
)

var cmd = &cobra.Command{
	Use:     "delete <funcName>",
	Short:   "deprovision the specified function",
	PreRunE: preRun,
	RunE:    run,
}

func init() {
	root.RootCmd.AddCommand(cmd)
}

func run(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	app, err := utils.GetApp(wd)
	if err != nil {
		return err
	}

	return app.Delete(args[0])
}

func preRun(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(wd + "/azul.json"); os.IsNotExist(err) {
		return errors.New("azul.json file not found")
	}

	if len(args) != 1 {
		return errors.New("delete should be call with one argument")
	}

	return nil
}
