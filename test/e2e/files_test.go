package e2e

import (
	"testing"

	"github.com/jkandasa/file-store/pkg/types"
	"github.com/jkandasa/file-store/pkg/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FileStoreTestSuite struct {
	suite.Suite
}

func TestFileStoreTestSuite(t *testing.T) {
	suite.Run(t, new(FileStoreTestSuite))
}

func (suite *FileStoreTestSuite) SetupTest() {
	t = suite.T()
	err = prepare(t)
	if err != nil {
		require.FailNow(t, "failed in prepare")
	}
}

func (suite *FileStoreTestSuite) TestAddFiles() {
	// add single file test
	files := []string{"./resources/file_add1.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	availableFiles, err := client.ListFiles()
	require.NoError(t, err)
	fileInfo, err := utils.GetFileInfo(files[0])
	require.NoError(t, err)

	require.Equal(t, 1, len(availableFiles))
	require.Equal(t, availableFiles[0].Name, fileInfo.Name)
	require.Equal(t, availableFiles[0].MD5Hash, fileInfo.MD5Hash)
	require.Equal(t, availableFiles[0].Size, fileInfo.Size)

	// add multiple file test
	files = []string{"./resources/file_add2.txt", "./resources/file_add3.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	availableFiles, err = client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 3, len(availableFiles))

	// add existing file test
	files = []string{"./resources/file_add1.txt"}
	err = client.AddFiles(files)
	require.Error(t, err)

	// add invalid file test
	files = []string{"./resources/file_add_invalid.txt"}
	err = client.AddFiles(files)
	require.Error(t, err)

	// add duplicate content file test
	files = []string{"./resources/file_add_duplicate.txt"}
	fileInfo, err = utils.GetFileInfo(files[0])
	require.NoError(t, err)
	err = client.AddFiles(files)
	require.NoError(t, err)
	availableFiles, err = client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 4, len(availableFiles))

	var actualFile *types.File
	for _, file := range availableFiles {
		if fileInfo.Name == file.Name {
			actualFile = &file
			break
		}
	}
	require.NotNil(t, actualFile)

	require.Equal(t, actualFile.Name, fileInfo.Name)
	require.Equal(t, actualFile.MD5Hash, fileInfo.MD5Hash)
	require.Equal(t, actualFile.Size, fileInfo.Size)
}

func (suite *FileStoreTestSuite) TestRemoveFiles() {
	// add new files
	files := []string{"./resources/file_add1.txt", "./resources/file_add2.txt", "./resources/file_add3.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	availableFiles, err := client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 3, len(availableFiles))

	// delete a file
	targetFile := availableFiles[0].Name
	err = client.RemoveFiles([]string{targetFile})
	require.NoError(t, err)
	availableFiles, err = client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 2, len(availableFiles))
	for _, file := range availableFiles {
		require.NotEqual(t, targetFile, file.Name)
	}

	// delete all the files
	targetFiles := []string{}
	for _, file := range availableFiles {
		targetFiles = append(targetFiles, file.Name)
	}
	err = client.RemoveFiles(targetFiles)
	require.NoError(t, err)
	availableFiles, err = client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 0, len(availableFiles))
}

func (suite *FileStoreTestSuite) TestWordCountFiles() {
	// add new files
	files := []string{"./resources/file_wc1.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	count, err := client.WordCount()
	require.NoError(t, err)
	require.Equal(t, int64(4), count)

	// add additional file and verify wc
	files = []string{"./resources/file_wc2.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	count, err = client.WordCount()
	require.NoError(t, err)
	require.Equal(t, int64(12), count)

	// remove a file and verify wc
	files = []string{"file_wc1.txt"}
	err = client.RemoveFiles(files)
	require.NoError(t, err)
	count, err = client.WordCount()
	require.NoError(t, err)
	require.Equal(t, int64(8), count)

	// remove all files and verify wc
	err = client.RemoveAllFiles()
	require.NoError(t, err)
	count, err = client.WordCount()
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func (suite *FileStoreTestSuite) TestFreqWords() {
	// add files for the test
	files := []string{"./resources/freq_words_cobra.txt", "./resources/freq_words_testify.txt"}
	err = client.AddFiles(files)
	require.NoError(t, err)
	availableFiles, err := client.ListFiles()
	require.NoError(t, err)
	require.Equal(t, 2, len(availableFiles))

	// order by dsc, limit: 5
	result, err := client.FreqWords(types.FreqWordsRequest{Limit: 5, OrderBy: types.OrderByDsc})
	require.NoError(t, err)
	// cat test/e2e/resources/freq_words_*.txt | tr -s ' ' '\n' | sort | uniq -c | sort -n | tail -n 5
	expectedResult := map[string]uint64{"the": uint64(69), "//": uint64(63), "to": uint64(43), "a": uint64(37), "that": uint64(32)}
	require.EqualValues(t, expectedResult, result)

	// order by dsc, limit: 10
	result, err = client.FreqWords(types.FreqWordsRequest{Limit: 10, OrderBy: types.OrderByDsc})
	require.NoError(t, err)
	// cat test/e2e/resources/freq_words_*.txt | tr -s ' ' '\n' | sort | uniq -c | sort -n | tail -n 10
	expectedResult = map[string]uint64{
		"the": uint64(69), "//": uint64(63), "to": uint64(43), "a": uint64(37), "that": uint64(32),
		"is": uint64(32), "for": uint64(24), "and": uint64(23), "of": uint64(21), "assert": uint64(19),
	}
	require.EqualValues(t, expectedResult, result)
}
