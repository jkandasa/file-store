package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jkandasa/file-store/pkg/types"
	"github.com/jkandasa/file-store/pkg/utils"
	"github.com/jkandasa/file-store/pkg/version"
)

type Client struct {
	ServerAddress string
	Insecure      bool
}

func NewClient(serverAddress string, insecure bool) *Client {
	return &Client{ServerAddress: serverAddress, Insecure: insecure}
}

func (c *Client) GetServerVersion() (*version.Version, error) {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	res, err := client.executeJson(fmt.Sprintf("%s/api/version", c.ServerAddress), http.MethodGet, nil, nil, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	ver := &version.Version{}
	err = json.Unmarshal(res.Body, ver)
	if err != nil {
		return nil, err
	}
	return ver, nil
}

func (c *Client) ListFiles() ([]types.File, error) {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	res, err := client.executeJson(fmt.Sprintf("%s/api/file/list", c.ServerAddress), http.MethodGet, nil, nil, nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	files := make([]types.File, 0)
	err = json.Unmarshal(res.Body, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (c *Client) RemoveFiles(files []string) error {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	_, err := client.executeJson(fmt.Sprintf("%s/api/file/remove", c.ServerAddress), http.MethodPost, nil, nil, files, 0)
	return err
}

func (c *Client) AddFiles(files []string) error {
	return c.addFiles(files, "/api/file/add")
}

func (c *Client) UpdateFiles(files []string) error {
	return c.addFiles(files, "/api/file/update")
}

func (c *Client) addFiles(files []string, apiPath string) error {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	for _, filename := range files {
		file, err := utils.GetFileInfo(filename)
		if err != nil {
			return err
		}
		res, err := client.executeJson(fmt.Sprintf("%s%s", c.ServerAddress, apiPath), http.MethodPost, nil, nil, file, 0)
		if err != nil {
			return err
		}

		var responseText string
		err = json.Unmarshal(res.Body, &responseText)
		if err != nil {
			return err
		}

		switch types.FileResponse(responseText) {
		case types.FileResponseCloned, types.FileResponseUpToDate:
			// no action needed

		case types.FileResponseNameExists:
			return fmt.Errorf("file exists:%s", file.Name)

		case types.FileResponseNotAvailable:
			err = c.addFileBucket(filename, file, client)
			if err != nil {
				return err
			}

		default:
			return errors.New(responseText)

		}
	}
	return nil
}

func (c *Client) addFileBucket(filename string, fileInfo *types.File, client *HttpClient) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	uuid := uuid.NewString()

	buf := make([]byte, types.BUCKET_SIZE)
	index := int64(0)
	reachedEOF := false
	for {
		length, err := file.Read(buf)

		if err != nil {
			if err == io.EOF {
				reachedEOF = true
			} else {
				return err
			}
		}

		if length > 0 {
			bucket := &types.Bucket{
				UUID:         uuid,
				File:         *fileInfo,
				Index:        index,
				IsLastBucket: reachedEOF || length < types.BUCKET_SIZE,
				BucketSize:   types.BUCKET_SIZE,
				Data:         buf[:length],
			}
			_, err = client.executeJson(fmt.Sprintf("%s/api/file/bucket", c.ServerAddress), http.MethodPost, nil, nil, bucket, 0)
			if err != nil {
				return err
			}
		}

		// terminate the infinite for loop
		if reachedEOF {
			return nil
		}
		index++
	}
}

func (c *Client) WordCount() (int64, error) {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	res, err := client.executeJson(fmt.Sprintf("%s/api/wc/count", c.ServerAddress), http.MethodGet, nil, nil, nil, http.StatusOK)
	if err != nil {
		return 0, err
	}

	var count int64
	err = json.Unmarshal(res.Body, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Client) FreqWords(freqWordsReq types.FreqWordsRequest) (map[string]uint64, error) {
	client := newHttpClient(c.Insecure, DefaultTimeout.String())
	res, err := client.executeJson(fmt.Sprintf("%s/api/wc/freq-words", c.ServerAddress), http.MethodPost, nil, nil, freqWordsReq, http.StatusOK)
	if err != nil {
		return nil, err
	}

	freqWords := make(map[string]uint64)
	err = json.Unmarshal(res.Body, &freqWords)
	if err != nil {
		return nil, err
	}

	return freqWords, nil
}
