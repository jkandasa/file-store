package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates/sync the files on the remote server",
	Example: `  # update files
  store update file1.txt file2.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)

		err := client.UpdateFiles(args)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on updating files. %s\n", err.Error())
			return
		}
		fmt.Fprintln(ioStreams.Out, "Files are updated successfully")
	},
}
