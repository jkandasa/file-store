package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jkandasa/file-store/pkg/types"
	"go.uber.org/zap"
)

// IsFileExists checks the file availability
func IsFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirExists checks the directory availability
func IsDirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CreateDir func
func CreateDir(dir string) error {
	if dir == "" {
		return nil
	}
	if !IsDirExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			zap.L().Error("failed to create a directory", zap.String("dir", dir))
			return err
		}
	}
	return nil
}

// removes the file
func RemoveFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		zap.L().Error("failed to remove a file", zap.String("file", file))
		return err
	}
	return nil
}

// CopyFile from a location to another location
func CopyFile(src, dst string, force bool) error {
	bufferSize := 1024

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if !force && IsFileExists(dst) {
		return fmt.Errorf("destination file exists: %s", dst)
	}

	// create target dir location
	dir, _ := filepath.Split(dst)
	err = CreateDir(dir)
	if err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, bufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// WriteFile func
func WriteFile(dir, filename string, data []byte) error {
	err := CreateDir(dir)
	if err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("%s/%s", dir, filename), data, os.ModePerm)
}

// AppendFile func
func AppendFile(dir, filename string, data []byte, offset int64) error {
	err := CreateDir(dir)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", dir, filename), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteAt(data, offset)
	return err
}

func GetFileInfo(filename string) (*types.File, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		zap.L().Error("error on getting file details", zap.String("filename", filename), zap.Error(err))
		return nil, err
	}
	rawFile, err := os.Open(filename)
	if err != nil {
		zap.L().Error("error on getting file content", zap.String("filename", filename), zap.Error(err))
		return nil, err
	}

	defer rawFile.Close()

	md5hash := md5.New()
	_, err = io.Copy(md5hash, rawFile)

	if err != nil {
		zap.L().Error("error on getting md5hash", zap.String("filename", filename), zap.Error(err))
		return nil, err
	}

	// get md5hash
	f := &types.File{
		Name:         fileInfo.Name(),
		Size:         fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime(),
		MD5Hash:      fmt.Sprintf("%x", md5hash.Sum(nil)),
	}

	fmt.Println("file info:", f)
	return f, nil
}
