package main

import (
	"github.com/jkandasa/file-store/cmd/client/command"
	types "github.com/jkandasa/file-store/pkg/types"
)

func main() {
	streams := types.NewStdStreams()
	command.Execute(streams)
}
