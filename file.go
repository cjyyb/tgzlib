package tgzlib

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

//WalkFuc ...
type WalkFuc func(path string, data []byte) error

//ReadDir read directory children
func ReadDir(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("couldn't find the directory %s", name)
		}
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("the %s isn't a directory", name)
	}

	dc, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return dc, nil
}

//ReadFileContent read file content
func ReadFileContent(fileName string) ([]byte, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("couldn't find the file %s", fileName)
		}
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("the %s isn't a file", fileName)
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil

}

//FileOrDirectory ...
func FileOrDirectory(name string) (bool, error) {
	f, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("couldn't find the file or directory %s", name)
		}
		return false, err
	}
	return f.IsDir(), nil
}

//WalkWriteFile ...
func WalkWriteFile(name, prefix string, wf WalkFuc) error {
	if _, err := os.Lstat(name); err != nil {
		name, err = prelinimaryWalkWriteFile(name)
		if err != nil {
			return err
		}
	}
	ele := filepath.Base(name)
	if strings.HasPrefix(ele, ".") {
		return nil
	}
	base := filepath.Join(prefix, ele)
	b, err := FileOrDirectory(name)
	if err != nil {
		return err
	}
	if !b {
		content, err := ReadFileContent(name)
		if err != nil {
			return fmt.Errorf("read file %s failed", err)
		}
		fmt.Println(base)
		if err := wf(base, content); err != nil {
			return err
		}
		return nil
	}
	drs, err := ReadDir(name)
	if err != nil {
		return err
	}
	for _, v := range drs {
		dir := filepath.Join(name, v)
		if err := WalkWriteFile(dir, base, wf); err != nil {
			return err
		}
	}
	return nil
}

func prelinimaryWalkWriteFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if IsSymLink(info) {
		ap, err := filepath.EvalSymlinks(path)
		if err != nil {
			return "", err
		}
		return ap, nil
	}
	return path, nil
}

//IsSymLink ...
func IsSymLink(info os.FileInfo) bool {
	return info.Mode() & os.ModeSymlink == 0
}
