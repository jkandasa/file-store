package command

import (
	"fmt"

	api "github.com/jkandasa/file-store/cmd/client/api"
	"github.com/jkandasa/file-store/pkg/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the client version information",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)
		fmt.Fprintln(ioStreams.Out, "Client Version:", version.Get().String())
		serverVersion, err := client.GetServerVersion()
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on getting server version. %s\n", err.Error())
			return
		}
		fmt.Fprintln(ioStreams.Out, "Server Version:", serverVersion.String())
	},
}
