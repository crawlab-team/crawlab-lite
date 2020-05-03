package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

func SaveFile(src io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func GetFileMD5(file io.Reader) string {
	h := md5.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func ContainsOnlyOneDir(path string) string {
	fileList, err := ioutil.ReadDir(path)
	if err == nil && len(fileList) == 1 && fileList[0].IsDir() {
		return fileList[0].Name()
	}
	return ""
}
