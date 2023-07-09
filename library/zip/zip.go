package zip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(srcFile string, destZip string, includePath *string) error {
	srcFilePath := ""
	srcIncludePath := ""

	if includePath != nil {
		srcIncludePath = *includePath
		srcFilePath = filepath.Join(srcFile, srcIncludePath)
	} else {
		srcFilePath = srcFile
	}

	if _, err := os.Stat(srcFilePath); os.IsNotExist(err) {
		return errors.New("资源不存在")
	}

	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	_ = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if includePath != nil && srcIncludePath != "" {
			inc := filepath.Join("/", srcIncludePath, "/")

			if strings.Index(path, inc+"/") == -1 {
				return nil
			}
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

//解压
func UnZip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer zipReader.Close()

	for _, f := range zipReader.File {
		fPath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}

			inFile.Close()
			outFile.Close()
		}
	}

	return nil
}
