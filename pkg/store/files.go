package store

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/jkandasa/file-store/pkg/types"
	"github.com/jkandasa/file-store/pkg/utils"
	"go.uber.org/zap"
)

var (
	filesStore []types.File
	storeMutex sync.RWMutex
)

func UpdateFilesStore() {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	filesStore = make([]types.File, 0)

	// create sore root directory if not exists
	utils.CreateDir(types.STORE_DATA_PATH)

	// Load file details to the store
	files, err := os.ReadDir(types.STORE_DATA_PATH)
	if err != nil {
		zap.L().Error("error on reading home directory", zap.String("homeDir", types.STORE_DATA_PATH), zap.Error(err))
		return
	}
	for _, entry := range files {
		if !entry.IsDir() {
			file, err := entry.Info()
			if err != nil {
				zap.L().Error("error on getting file detail", zap.String("name", entry.Name()), zap.Error(err))
				continue
			}

			// if !strings.HasSuffix(strings.ToLower(file.Name()), types.TEXT_FILE_EXTENSION) {
			// 	continue
			// }

			rawFile, err := os.Open(filepath.Join(types.STORE_DATA_PATH, file.Name()))
			if err != nil {
				zap.L().Error("error on getting file content", zap.String("filename", file.Name()), zap.Error(err))
				continue
			}

			defer rawFile.Close()

			md5hash := md5.New()
			_, err = io.Copy(md5hash, rawFile)

			if err != nil {
				zap.L().Error("error on getting md5hash", zap.String("filename", file.Name()), zap.Error(err))
				continue
			}

			// get md5hash
			f := types.File{
				Name:         file.Name(),
				Size:         file.Size(),
				ModifiedTime: file.ModTime(),
				MD5Hash:      fmt.Sprintf("%x", md5hash.Sum(nil)),
			}
			filesStore = append(filesStore, f)
		}
	}
}

func Update(file types.File) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	for index := range filesStore {
		if file.Name == filesStore[index].Name {
			filesStore[index] = file
			return
		}
	}

	filesStore = append(filesStore, file)
}

func Delete(filename string) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	for index := range filesStore {
		if filename == filesStore[index].Name {
			filesStore = append(filesStore[:index], filesStore[index+1:]...)
		}
	}
}

func GetByName(name string) *types.File {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	for index := range filesStore {
		if name == filesStore[index].Name {
			return filesStore[index].Clone()
		}
	}

	return nil
}

func GetByHash(md5Hash string) *types.File {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	for index := range filesStore {
		if md5Hash == filesStore[index].MD5Hash {
			return filesStore[index].Clone()
		}
	}

	return nil
}

func ListFiles() []types.File {
	storeMutex.RLock()
	defer storeMutex.RUnlock()

	clonedFiles := make([]types.File, len(filesStore))

	for index := range filesStore {
		clonedFiles[index] = *filesStore[index].Clone()
	}

	return clonedFiles
}
