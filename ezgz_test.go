package ezgz

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fileContents = "Here is some example text"
)

type FakeWriter struct{}

func (w *FakeWriter) Write(p []byte) (n int, err error) { return 0, nil }

func TestFileSuccessfullyZipped(t *testing.T) {
	fileToZip := getFileToZip(t)
	defer os.Remove(fileToZip)

	outputFile := getOutputFilePath(t)
	defer os.RemoveAll(outputFile)

	err := ZipToFile(fileToZip, outputFile)
	if assert.NoError(t, err, "unexpected error occurred when converting file to zip") {
		_, err := os.Stat(outputFile)
		assert.NoError(t, err, "error getting stats of output file")
	}
}

func TestFolderSuccessfullyZipped(t *testing.T) {
	folderToZip := getFolderToZip(t)
	defer os.RemoveAll(folderToZip)

	outputFile := getOutputFilePath(t)
	defer os.RemoveAll(outputFile)

	err := ZipToFile(folderToZip, outputFile)
	if assert.NoError(t, err, "unexpected error occurred when converting folder to zip") {
		_, err := os.Stat(outputFile)
		assert.NoError(t, err, "error getting stats of output file")
	}
}

func TestNotFileErrorIsReturnedWhenPassingAFolder(t *testing.T) {
	folder := getFolderToZip(t)
	defer os.RemoveAll(folder)
	var writer = new(FakeWriter)

	err := ZipFileToWriter(folder, writer)

	if assert.Error(t, err) {
		assert.Equal(t, ErrNotFile, err, "error returned was not ErrNotFile as expected")
	}
}

func TestNotDirErrorIsReturnedWhenPassingAFile(t *testing.T) {
	file := getFileToZip(t)
	defer os.Remove(file)
	var writer = new(FakeWriter)

	err := ZipFolderToWriter(file, writer)

	if assert.Error(t, err) {
		assert.Equal(t, ErrNotDir, err, "error returned was not ErrNotDir as expected")
	}
}

func getFolderToZip(t *testing.T) string {
	testDir, err := ioutil.TempDir("", "ezgz-test-output-dir-")
	require.NoError(t, err, "unable to create temp file for test")

	file, err := ioutil.TempFile(testDir, "ezgz-test-file-*.txt")
	require.NoError(t, err, "unable to create temp file for test")
	defer file.Close()
	_, err = file.WriteString(fileContents)
	require.NoError(t, err, "unable to add content to temp file")
	return testDir
}

func getFileToZip(t *testing.T) string {
	file, err := ioutil.TempFile("", "ezgz-test-file-*.txt")
	require.NoError(t, err, "unable to create temp file for test")
	defer file.Close()
	_, err = file.WriteString(fileContents)
	require.NoError(t, err, "unable to add content to temp file")
	return file.Name()
}

func getOutputFilePath(t *testing.T) string {
	outputFile, err := ioutil.TempDir("", "ezgz-test-output-dir-")
	require.NoError(t, err, "unable to create temp file for test")
	return path.Join(outputFile, "output.gz")
}
