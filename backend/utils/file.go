package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
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
