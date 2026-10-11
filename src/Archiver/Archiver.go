package Archiver

import (
	"Helpers"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const DIRECTORYSEPERATOR = "/"

func CompressFile(target string, source string, excludePath string, comment string) (bool, error) {
	info, err := os.Stat(source)
	if err != nil {
		return false, err
	}

	var baseDir string
	if info.IsDir() {
		newSource := Helpers.RemoveAll(source, excludePath)
		baseDir = filepath.Base(newSource)
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = zipfile.Close()
	}()
	archive := zip.NewWriter(zipfile)
	defer func() {
		_ = archive.Close()
	}()
	if len(comment) > 0 {
		err := archive.SetComment(comment)
		if err != nil {
			return false, err
		}
	}

	_ = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			Helpers.HandleError(err, 1010, fmt.Sprintf("FileInfoHeader error is %v", err))
			return err
		}

		baseDir = Helpers.Trim(baseDir)
		if len(baseDir) > 0 {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += DIRECTORYSEPERATOR
		} else {
			header.Method = zip.Store
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			Helpers.HandleError(err, 1011, fmt.Sprintf("CreateHeader error is %v", err))
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			Helpers.HandleError(err, 1012, fmt.Sprintf("file open error is %v", err))
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			Helpers.HandleError(err, 1013, fmt.Sprintf("file copy error is %v", err))
			return err
		}

		err = file.Close()
		if err != nil {
			Helpers.HandleError(err, 1014, fmt.Sprintf("file close error is %v", err))
		}
		return err
	})

	return true, err
}
func Decompress(target string, source string) (bool, error) {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return false, err
	}

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			Helpers.CreateDirectory(target+DIRECTORYSEPERATOR+f.Name, true)
		} else {
			fe := target + DIRECTORYSEPERATOR + f.Name

			fz, err := f.Open()
			if err != nil {
				return false, err
			}
			fw, err := os.OpenFile(fe, os.O_CREATE|os.O_RDWR, f.Mode())
			if err == nil {
				if _, err := io.Copy(fw, fz); err != nil {
					return false, err
				}
			}
			err = fz.Close()
		}
	}
	err = archive.Close()
	return true, nil
}
