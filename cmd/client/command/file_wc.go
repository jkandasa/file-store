package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(wcCmd)
}

var wcCmd = &cobra.Command{
	Use:   "wc",
	Short: "Prints the word count from the available text files on the remote server",
	Example: `  # word count
  store wc`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)
		count, err := client.WordCount()
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error: %s\n", err.Error())
			return
		}

		fmt.Fprintln(ioStreams.Out, count)
	},
}
