package command

import (
	"fmt"
	"os"
	"strings"

	printer "github.com/jkandasa/file-store/cmd/client/printer"
	types "github.com/jkandasa/file-store/pkg/types"

	"github.com/spf13/cobra"
)

const (
	SERVER_ADDRESS_ENV = "STORE_SERVER"
)

var (
	ioStreams types.IOStreams // read and write to this stream

	hideHeader    bool
	pretty        bool
	insecure      bool
	outputFormat  string
	serverAddress string

	rootCliLong = `Storage Client
  
This client helps you to control your storage server from the command line.
`
)

var rootCmd = &cobra.Command{
	Use:   "storage",
	Short: "storage",
	Long:  rootCliLong,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(ioStreams.Out)
		cmd.SetErr(ioStreams.ErrOut)
	},
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", printer.OutputConsole, "output format. options: yaml, json, console")
	rootCmd.PersistentFlags().BoolVar(&hideHeader, "hide-header", false, "hides the header on the console output")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "JSON pretty print")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", true, "connect to server in insecure mode")

	// update server address
	if os.Getenv(SERVER_ADDRESS_ENV) == "" {
		serverAddress = "http://127.0.0.1:8080"
	} else {
		serverAddress = strings.TrimSuffix(os.Getenv(SERVER_ADDRESS_ENV), "/")
	}
}

func Execute(streams types.IOStreams) {
	ioStreams = streams
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(ioStreams.ErrOut, err)
		os.Exit(1)
	}
}
