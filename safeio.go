package iohelper

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// SafeWrite writes the given data and hash to file. Uses backupFile as temporary storage to be resistant against unexpected crashes.
func SafeWrite(data []byte, file, backupFile string) error {
	// write to backup file first so previous file is retained after errors and unexpected crashes
	f, err := os.Create(backupFile)
	if err != nil {
		return err
	}

	// write hash encoded as text for better readability
	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	if _, err := f.WriteString("DataHash:" + base64.RawURLEncoding.EncodeToString(hash) + ";"); err != nil {
		f.Close()
		return err
	}

	// now write actual data
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// backup exists, now remove old file and replace by freshly written backup
	os.Remove(file)
	if err := os.Rename(backupFile, file); err != nil {
		return err
	}

	return nil
}

// SafeRead reads alls data from file. Alternatively reads from backupFile if the default file does not exist or is broken.
func SafeRead(file, backupFile string) ([]byte, error) {
	// define generic functionality for reading and validating a file with hash-prefix
	read := func(file string) ([]byte, error) {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		if len(data) < 53 {
			return nil, fmt.Errorf("Input file is too short")
		}

		// hash prefix is like "DataHash:n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg;" -> always 53 chars
		prefix := data[:53]
		if !strings.HasPrefix(string(prefix), "DataHash:") {
			return nil, fmt.Errorf("Input file does not contain DataHash: prefix")
		}
		if !strings.HasSuffix(string(prefix), ";") {
			return nil, fmt.Errorf("Input file hash is not terminated with semicolon")
		}

		// strip hash prefix to obtain data suffix
		content := data[53:]

		// validate content
		hasher := sha256.New()
		hasher.Write(content)
		hash := hasher.Sum(nil)
		if string(prefix[9:52]) != base64.RawURLEncoding.EncodeToString(hash) {
			return nil, fmt.Errorf("Input file hash does not match")
		}

		return content, nil
	}

	// try to read default file first
	data, err := read(file)
	if err != nil {
		// some error occured -> try to read backup file now
		data, backupErr := read(backupFile)
		if backupErr != nil {
			// still failed -> return error of default file here
			return nil, err
		}
		// return backup data
		return data, nil
	}
	// return data
	return data, nil
}
