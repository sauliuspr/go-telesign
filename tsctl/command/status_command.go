package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	referenceID string
)

//NewStatusCommand return the cobra command for status
func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [reference_id]",
		Short: "Gets status the key or a range of keys",
		Run:   statusCommandFunc,
	}

	cmd.Flags().StringVar(&referenceID, "reference-id", "r", "Reference ID of Messaging Transaction")
	return cmd

}

func statusCommandFunc(cmd *cobra.Command, args []string) {
	ctx, cancel := commandCtx(cmd)
	fmt.Println("Running status command")

}
