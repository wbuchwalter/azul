package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wbuchwalter/lox/lox"
)

var deployCmd = &cobra.Command{
	Use:   "deploy main.go",
	Short: "deploy the specified file to Azure Function",
	Long:  `deploy - Compile and deploy the specified go file to azure function`,
	Run: func(cmd *cobra.Command, args []string) {
		err := preRun(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = lox.Deploy(args[0])
		if err != nil {
			fmt.Println("ERR: ", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

func preRun(args []string) error {
	if len(args) > 1 {
		return errors.New("Only one file should be passed to deploy.")
	}
	return nil
}
