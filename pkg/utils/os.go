package utils

import "os"

func CreateDir(d string) error {
	return os.MkdirAll(d, os.ModePerm)
}
