package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the files to remote server",
	Example: `  # add files
  store add file1.txt file2.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)

		err := client.AddFiles(args)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on adding files. %s\n", err.Error())
			return
		}
		fmt.Fprintln(ioStreams.Out, "Files are added successfully")
	},
}
