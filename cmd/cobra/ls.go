package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wedojava/enit"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all secrets in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := enit.File(encodingKey, secretPath())
		list, err := v.List()
		if err != nil {
			fmt.Println("no value set")
			return
		}
		for k, v := range list {
			fmt.Printf("%s = %s\n", k, v)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
