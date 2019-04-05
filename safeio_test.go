package iohelper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeWrite(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	assert.NoError(t, SafeWrite([]byte("test"), dir+"/file.txt", dir+"/backup.txt"))
}

func TestSafeWriteEmpty(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	assert.NoError(t, SafeWrite([]byte{}, dir+"/file.txt", dir+"/backup.txt"))
}

func TestSafeRead(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("test"), dir+"/file.txt", dir+"/backup.txt")
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestSafeReadEmpty(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	ioutil.WriteFile(dir+"/file.txt", []byte("DataHash:47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU;not empty"), os.ModePerm)
	_, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.Error(t, err)
}

func TestSafeReadEmptyMalformed(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte{}, dir+"/file.txt", dir+"/backup.txt")
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, data)
}

func TestSafeReadBackup(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("test"), dir+"/file.txt", dir+"/backup.txt")
	os.Rename(dir+"/file.txt", dir+"/backup.txt")
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestSafeReadTooShort(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("test"), dir+"/file.txt", dir+"/backup.txt")
	os.Rename(dir+"/file.txt", dir+"/backup.txt")
	ioutil.WriteFile(dir+"/file.txt", []byte("Data"), os.ModePerm)
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestSafeReadNoHashPrefix(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("foobar"), dir+"/file.txt", dir+"/backup.txt")
	os.Rename(dir+"/file.txt", dir+"/backup.txt")
	ioutil.WriteFile(dir+"/file.txt", []byte("another file prefix that is absolutely not a valid sha256 hash of the real content"), os.ModePerm)
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "foobar", string(data))
}

func TestSafeReadMalformedPrefix(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("foobar"), dir+"/file.txt", dir+"/backup.txt")
	os.Rename(dir+"/file.txt", dir+"/backup.txt")
	ioutil.WriteFile(dir+"/file.txt", []byte("DataHash:n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg test"), os.ModePerm)
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "foobar", string(data))
}

func TestSafeReadWrongHash(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	SafeWrite([]byte("test"), dir+"/file.txt", dir+"/backup.txt")
	os.Rename(dir+"/file.txt", dir+"/backup.txt")
	ioutil.WriteFile(dir+"/file.txt", []byte("DataHash:n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg;foobar"), os.ModePerm)
	data, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(data))
}

func TestSafeReadError(t *testing.T) {
	dir := getTmpDir()
	defer os.RemoveAll(dir)
	_, err := SafeRead(dir+"/file.txt", dir+"/backup.txt")
	assert.Error(t, err)
}

func getTmpDir() string {
	dir, err := ioutil.TempDir("", "iohelper-test-")
	if err != nil {
		panic(err)
	}
	return dir
}
