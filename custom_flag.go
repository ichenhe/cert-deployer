package main

import (
	"fmt"
	"github.com/ichenhe/cert-deployer/utils"
	"io/ioutil"
)

// fileType is used to receive and load file flags.
type fileType struct {
	filename string
	data     []byte
}

func (f *fileType) Set(value string) error {
	if !utils.IsFile(value) {
		return fmt.Errorf("file does not exist: %s", value)
	}
	if data, err := ioutil.ReadFile(value); err != nil {
		return fmt.Errorf("failed to read file '%s': %w", value, err)
	} else {
		f.filename = value
		f.data = data
	}
	return nil
}

func (f *fileType) String() string {
	return f.filename
}
