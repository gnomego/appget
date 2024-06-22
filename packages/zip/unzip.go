package zip

import (
	z "archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipArchive string, dest string) error {

	archive, err := z.OpenReader(zipArchive)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dest, f.Name)
		fmt.Println("extracting " + filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			fmt.Println("invalid file path " + filePath)
			return fmt.Errorf("%s: illegal file path", filePath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func UnzipSingleFile(zipArchive string, file string, dest string) (string, error) {
	r, err := z.OpenReader(zipArchive)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != file {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		// Create the file
		destPath := filepath.Join(dest, f.Name)
		dstFile, err := os.Create(destPath)
		if err != nil {
			return "", err
		}

		// Copy the file
		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return "", err
		}

		// Close the file
		err = dstFile.Close()
		if err != nil {
			return "", err
		}

		return destPath, nil

	}

	return "", fmt.Errorf("file not found in archive")
}
