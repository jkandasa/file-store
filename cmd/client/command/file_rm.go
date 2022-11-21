package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes the files from the remote server",
	Example: `  # remove files
  store rm file1.txt file2.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)
		err := client.RemoveFiles(args)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on removing files. %s\n", err.Error())
			return
		}
		fmt.Fprintln(ioStreams.Out, "files removed successfully")

	},
}
