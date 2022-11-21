package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	printer "github.com/jkandasa/file-store/cmd/client/printer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Prints the available file details from the server",
	Example: `  # list files
  store ls`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)
		files, err := client.ListFiles()
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on getting files. %s\n", err.Error())
			return
		}

		if len(files) == 0 && outputFormat == printer.OutputConsole {
			fmt.Fprintln(ioStreams.Out, "No files found")
			return
		}

		rows := make([]interface{}, len(files))
		for index := range files {
			rows[index] = files[index]
		}

		headers := []string{"name", "size", "md5_hash", "modified_time"}
		printer.Print(ioStreams.Out, headers, rows, hideHeader, outputFormat, pretty)
	},
}
