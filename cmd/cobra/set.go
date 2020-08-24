package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wedojava/enit"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := enit.File(encodingKey, secretPath())
		fmt.Println(args)
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			fmt.Println("value set error: ", err)
			return
		}
		fmt.Println("Value set successfully.")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
