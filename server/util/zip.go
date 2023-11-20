package util

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Unzip(src string, dest string, skip int) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if skip == 1 {
			filePath = path.Join(dest, file.Name[strings.Index(file.Name, "/")+1:])
		}
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
