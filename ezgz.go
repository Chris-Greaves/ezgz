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
	// ErrNotFile is returned when the path given is not a file
	ErrNotFile = errors.New("path given is not a file")

	// ErrNotDir is returned when the path given is not a directory
	ErrNotDir = errors.New("path given is not a directory")

	// ErrInvalidSourcePath is returned when the path isn't recognised as a file or directory
	ErrInvalidSourcePath = errors.New("path given was not recognised as a file or directory")
)

// ZipToFile will take a source path to either a file or directory, and zip it into a file.
func ZipToFile(pathToSource, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	return ZipToWriter(pathToSource, file)
}

// ZipToWriter will take a source path to either a file or directory, and write it to the writer provided.
//
// Will return errInvalidSourcePath if source path isn't recognised as either file or directory.
func ZipToWriter(pathToSource string, writer io.Writer) error {
	if err := ZipFileToWriter(pathToSource, writer); err != nil {
		if err != ErrNotFile {
			return err
		}
	} else {
		return nil
	}

	if err := ZipFolderToWriter(pathToSource, writer); err != nil {
		if err != ErrNotDir {
			return err
		}
	} else {
		return nil
	}

	return ErrInvalidSourcePath
}

// ZipFileToWriter will take the file located by the source path, and write it, to the writer provided.
//
// Will return errNotFile is source path doesn't point to a file.
func ZipFileToWriter(pathToSource string, writer io.Writer) error {
	if !isFile(pathToSource) {
		return ErrNotFile
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

// ZipFolderToWriter will take the folder located by the source path, and write it, to the writer provided.
//
// Will return errNotDir is source path doesn't point to a directory.
func ZipFolderToWriter(pathToSource string, writer io.Writer) error {
	if !isDir(pathToSource) {
		return ErrNotDir
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
