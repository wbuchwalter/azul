package cmd

import "github.com/spf13/cobra"

var deployCmd = &cobra.Command{
	Use:   "deploy main.go",
	Short: "deploy the specified file to Azure Function",
	Long:  `deploy - Compile and deploy the specified go file to azure function`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
