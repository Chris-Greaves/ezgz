package ezgz

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	errNotFile          = errors.New("path given is not a file")
	errNotDir           = errors.New("path given is not a directory")
	errNeitherFileOrDir = errors.New("path given is neither a file or directory")
)

func ZipToFile(pathToSource, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	return ZipToWriter(pathToSource, file)
}

func ZipToWriter(pathToSource string, writer io.Writer) error {
	if err := ZipFileToWriter(pathToSource, writer); err != nil {
		if err != errNotFile {
			return err
		}
	} else {
		return nil
	}

	if err := ZipFolderToWriter(pathToSource, writer); err != nil {
		if err != errNotDir {
			return err
		}
	} else {
		return nil
	}

	return errNeitherFileOrDir
}

func ZipFileToWriter(pathToSource string, writer io.Writer) error {
	if !isFile(pathToSource) {
		return errNotFile
	}

	gzipwriter := gzip.NewWriter(writer)
	defer gzipwriter.Close()

	file, err := os.Open(pathToSource)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(gzipwriter, file)

	return err
}

func ZipFolderToWriter(pathToSource string, writer io.Writer) error {
	if !isDir(pathToSource) {
		return errNotDir
	}

	gzipwriter := gzip.NewWriter(writer)
	defer gzipwriter.Close()
	tarball := tar.NewWriter(gzipwriter)
	defer tarball.Close()

	info, err := os.Stat(pathToSource)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(pathToSource)
	}

	return filepath.Walk(pathToSource,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, pathToSource))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}
