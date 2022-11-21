package e2e

import (
	"testing"

	api "github.com/jkandasa/file-store/cmd/client/api"
)

var (
	t             *testing.T
	err           error
	serverAddress = "http://127.0.0.1:8080"
	insecure      = true
	client        *api.Client
)

func prepare(t *testing.T) error {
	client = api.NewClient(serverAddress, insecure)

	// delete all the files
	err = client.RemoveAllFiles()
	return err
}
